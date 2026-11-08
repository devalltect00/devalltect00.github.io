// Receiver-side local validation neither extracts files nor authorizes Pages publication.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
)

// deploymentDirectory bounds the single ordinary wrapper inventory before archive/zip allocates it.
func deploymentDirectory(raw []byte) bool {
	if len(raw) < 22 || len(raw) > MaxDeployWrapper {
		return false
	}
	end := len(raw) - 22
	b := raw[end:]
	return binary.LittleEndian.Uint32(b) == 0x06054b50 && binary.LittleEndian.Uint16(b[4:]) == 0 && binary.LittleEndian.Uint16(b[6:]) == 0 && binary.LittleEndian.Uint16(b[8:]) == 1 && binary.LittleEndian.Uint16(b[10:]) == 1 && binary.LittleEndian.Uint16(b[20:]) == 0 && binary.LittleEndian.Uint32(b[12:]) <= 4096 && uint64(binary.LittleEndian.Uint32(b[16:]))+uint64(binary.LittleEndian.Uint32(b[12:])) == uint64(end)
}

// verifyDeploymentInput independently reconstructs source TAR and compares every pinned uploaded byte.
// It returns integrity only: no provider provenance/current CMS authority, filesystem write or deployment is performed.
func verifyDeploymentInput(ctx context.Context, r DeployRequest, candidate, uploaded []byte) error {
	raw, _ := json.Marshal(r)
	if _, e := ParseDeployRequest(raw); e != nil || ctx.Err() != nil || !deploymentDirectory(uploaded) || hash(uploaded) != r.Evidence.ZIPDigest {
		return errInvalid
	}
	canonical, e := pack(ctx, r.Stage, candidate)
	if e != nil || len(canonical) != r.Evidence.TARBytes || hash(canonical) != r.Evidence.TARDigest {
		return errInvalid
	}
	z, e := zip.NewReader(bytes.NewReader(uploaded), int64(len(uploaded)))
	if e != nil || len(z.File) != 1 {
		return errInvalid
	}
	f := z.File[0]
	if f.Name != "artifact.tar" || !f.Mode().IsRegular() || f.Flags & ^uint16(0x808) != 0 || (f.Method != zip.Store && f.Method != zip.Deflate) || f.UncompressedSize64 != uint64(len(canonical)) || f.CompressedSize64 > uint64(len(uploaded)) {
		return errInvalid
	}
	reader, e := f.Open()
	if e != nil {
		return errInvalid
	}
	body, e := io.ReadAll(io.LimitReader(contextArchiveReader{ctx, reader}, int64(len(canonical))+1))
	closed := reader.Close()
	if e != nil || closed != nil || ctx.Err() != nil || !bytes.Equal(body, canonical) {
		return errInvalid
	}
	return nil
}
