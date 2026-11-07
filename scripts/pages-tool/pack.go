// Packing preserves every candidate file; only the pinned private gzip identity can produce a TAR.
package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/binary"
	jsonv2 "encoding/json/v2"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"
)

const (
	MaxArchive  = 8 << 20
	MaxExpanded = 32 << 20
	MaxEntries  = 256
)

// hash covers binary assets without interpretation.
func hash(raw []byte) string { return fmt.Sprintf("%x", sha256.Sum256(raw)) }

// portableMember rejects traversal, ambiguous Windows names, private hidden roots and case collisions before use.
func portableMember(name string) bool {
	if len(name) < 1 || len(name) > 1024 || strings.HasPrefix(name, "/") {
		return false
	}
	name = strings.TrimSuffix(name, "/")
	if name == "" || path.Clean(name) != name {
		return false
	}
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || strings.ContainsRune("_-. /", c)) {
			return false
		}
	}
	for i, part := range strings.Split(name, "/") {
		if part == "" || part == "." || part == ".." || len(part) > 240 || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") || strings.HasPrefix(part, ".") && !(i == 0 && part == ".well-known") {
			return false
		}
		base := strings.ToUpper(strings.SplitN(part, ".", 2)[0])
		if base == "CON" || base == "PRN" || base == "AUX" || base == "NUL" || len(base) == 4 && (strings.HasPrefix(base, "COM") || strings.HasPrefix(base, "LPT")) && base[3] >= '0' && base[3] <= '9' {
			return false
		}
	}
	return true
}

// boundedDirectory checks ordinary EOCD metadata before archive/zip allocates the file inventory. ZIP64/multidisk/trailers are unsupported.
func boundedDirectory(raw []byte) bool {
	if len(raw) < 22 || len(raw) > MaxArchive {
		return false
	}
	// Actions candidates have no archive comment: the ordinary EOCD must be the final 22 bytes.
	end := len(raw) - 22
	b := raw[end:]
	if binary.LittleEndian.Uint32(b) != 0x06054b50 || binary.LittleEndian.Uint16(b[4:]) != 0 || binary.LittleEndian.Uint16(b[6:]) != 0 || binary.LittleEndian.Uint16(b[20:]) != 0 {
		return false
	}
	count := binary.LittleEndian.Uint16(b[10:])
	size := uint64(binary.LittleEndian.Uint32(b[12:]))
	offset := uint64(binary.LittleEndian.Uint32(b[16:]))
	return count >= 3 && count <= MaxEntries && binary.LittleEndian.Uint16(b[8:]) == count && size <= 128<<10 && offset+size == uint64(end)
}

// contextArchiveReader checks cancellation between bounded decompression reads.
type contextArchiveReader struct {
	context.Context
	r io.Reader
}

func (r contextArchiveReader) Read(p []byte) (int, error) {
	if e := r.Err(); e != nil {
		return 0, e
	}
	if len(p) > 32768 {
		p = p[:32768]
	}
	return r.r.Read(p)
}

// pack reconstructs the versioned canonical payload without a private module dependency.
// Exact hashes make writer drift fail closed; joint compiled tests must accompany changes to either writer.
func pack(ctx context.Context, r StageRequest, raw []byte) ([]byte, error) {
	if ctx.Err() != nil || !boundedDirectory(raw) || hash(raw) != r.ZIPDigest {
		return nil, errInvalid
	}
	z, e := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if e != nil || len(z.File) > MaxEntries {
		return nil, errInvalid
	}
	sort.Slice(z.File, func(i, j int) bool { return z.File[i].Name < z.File[j].Name })
	var out, compressed bytes.Buffer
	g := gzip.NewWriter(&compressed)
	g.Header.OS = 255
	t := tar.NewWriter(io.MultiWriter(&out, g))
	seen := map[string]bool{}
	var total int64
	var marker []byte
	regular := 0
	for _, f := range z.File {
		if ctx.Err() != nil || !portableMember(f.Name) || f.Flags&^uint16(0x080e) != 0 || f.Method != zip.Store && f.Method != zip.Deflate || f.UncompressedSize64 > MaxExpanded || f.CompressedSize64 > uint64(len(raw)) {
			return nil, errInvalid
		}
		dir := f.FileInfo().IsDir()
		if !dir && !f.Mode().IsRegular() || dir != strings.HasSuffix(f.Name, "/") {
			return nil, errInvalid
		}
		key := strings.ToLower(strings.TrimSuffix(f.Name, "/"))
		if _, ok := seen[key]; ok {
			return nil, errInvalid
		}
		for parent := path.Dir(key); parent != "."; parent = path.Dir(parent) {
			if isDir, ok := seen[parent]; ok && !isDir {
				return nil, errInvalid
			}
		}
		if !dir {
			for prior := range seen {
				if strings.HasPrefix(prior, key+"/") {
					return nil, errInvalid
				}
			}
		}
		seen[key] = dir
		if dir {
			if f.UncompressedSize64 != 0 || f.CompressedSize64 != 0 || f.CRC32 != 0 || f.Method != zip.Store {
				return nil, errInvalid
			}
			continue
		}
		limit := int64(MaxExpanded) - total
		if f.Name == ".well-known/homepage-content.json" && limit > 8192 {
			limit = 8192
		}
		if f.UncompressedSize64 > uint64(limit) {
			return nil, errInvalid
		}
		file, e := f.Open()
		if e != nil {
			return nil, errInvalid
		}
		body, e := io.ReadAll(io.LimitReader(contextArchiveReader{ctx, file}, limit+1))
		closed := file.Close()
		if e != nil || closed != nil || int64(len(body)) > limit || uint64(len(body)) != f.UncompressedSize64 {
			return nil, errInvalid
		}
		total += int64(len(body))
		regular++
		if f.Name == ".well-known/homepage-content.json" {
			marker = body
		}
		if t.WriteHeader(&tar.Header{Name: f.Name, Mode: 0644, Size: int64(len(body)), ModTime: time.Unix(0, 0).UTC(), Typeflag: tar.TypeReg, Format: tar.FormatUSTAR}) != nil {
			return nil, errInvalid
		}
		if _, e = t.Write(body); e != nil {
			return nil, errInvalid
		}
	}
	if regular < 3 || total < 1 || t.Close() != nil || g.Close() != nil || ctx.Err() != nil || hash(compressed.Bytes()) != r.PayloadDigest {
		return nil, errInvalid
	}
	// The public marker binds request labels; private Go remains responsible for frozen visitor semantics.
	var m struct {
		SchemaVersion    int      `json:"schema_version"`
		Source           string   `json:"source"`
		Scope            string   `json:"scope"`
		ReleaseID        string   `json:"release_id"`
		Sequence         string   `json:"sequence"`
		SnapshotDigest   string   `json:"snapshot_digest"`
		ProjectionDigest string   `json:"projection_digest"`
		Kinds            []string `json:"consumed_kinds"`
	}
	if jsonv2.Unmarshal(marker, &m, jsonv2.RejectUnknownMembers(true)) != nil || m.SchemaVersion != 1 || m.Source != "cms-build-preview" || m.Scope != "astro-project-profile-preview-v1" || m.ReleaseID != r.Build.ReleaseID || m.Sequence != r.Build.Sequence || m.SnapshotDigest != r.Build.SnapshotDigest || m.ProjectionDigest != r.Build.ProjectionDigest || strings.Join(m.Kinds, ",") != "project,profile" {
		return nil, errInvalid
	}
	for _, name := range []string{"index.html", "id/index.html"} {
		if isDir, ok := seen[name]; !ok || isDir {
			return nil, errInvalid
		}
	}
	return out.Bytes(), nil
}
