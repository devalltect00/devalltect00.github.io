// Protocol and archive fixtures contain only synthetic public content and credentials.
package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// fixture independently builds expected USTAR/gzip around complete public ZIP data.
func fixture(t *testing.T) (StageRequest, []byte, []byte) {
	t.Helper()
	r := StageRequest{SchemaVersion: 1, Scope: StageScope, ID: "33333333-3333-4333-8333-333333333333", CandidateID: "22222222-2222-4222-8222-222222222222", RunID: "123", ArtifactID: "456", HeadSHA: strings.Repeat("a", 40), Build: Request{1, "astro-project-profile-build-v1", "11111111-1111-4111-8111-111111111111", "1", strings.Repeat("b", 64), strings.Repeat("c", 64)}}
	marker, _ := json.Marshal(map[string]any{"schema_version": 1, "source": "cms-build-preview", "scope": "astro-project-profile-preview-v1", "release_id": r.Build.ReleaseID, "sequence": "1", "snapshot_digest": r.Build.SnapshotDigest, "projection_digest": r.Build.ProjectionDigest, "consumed_kinds": []string{"project", "profile"}})
	files := map[string][]byte{"index.html": []byte("public en"), "id/index.html": []byte("public id"), ".well-known/homepage-content.json": marker, "images/demo.bin": {0, 1, 255}, "zero.txt": {}}
	keys := []string{}
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var z, tb, gb bytes.Buffer
	zw := zip.NewWriter(&z)
	tw := tar.NewWriter(&tb)
	for _, name := range keys {
		f, e := zw.Create(name)
		if e != nil {
			t.Fatal(e)
		}
		f.Write(files[name])
		if tw.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(files[name])), ModTime: time.Unix(0, 0).UTC(), Typeflag: tar.TypeReg, Format: tar.FormatUSTAR}) != nil {
			t.Fatal("header")
		}
		tw.Write(files[name])
	}
	zw.Close()
	tw.Close()
	g := gzip.NewWriter(&gb)
	g.Header.OS = 255
	g.Write(tb.Bytes())
	g.Close()
	r.ZIPDigest = hash(z.Bytes())
	r.PayloadDigest = hash(gb.Bytes())
	return r, z.Bytes(), tb.Bytes()
}

// TestPackAndCLI protects binary/empty members, canonical output and new-file-only local writes.
func TestPackAndCLI(t *testing.T) {
	r, z, want := fixture(t)
	out, e := pack(context.Background(), r, z)
	if e != nil || !bytes.Equal(out, want) {
		t.Fatal("canonical pack", e)
	}
	root := t.TempDir()
	request := filepath.Join(root, "request.json")
	archive := filepath.Join(root, "candidate.zip")
	target := filepath.Join(root, "artifact.tar")
	body, _ := json.Marshal(r)
	os.WriteFile(request, body, 0600)
	os.WriteFile(archive, z, 0600)
	args := []string{"pack", "--request", request, "--zip", archive, "--out", target}
	var printed bytes.Buffer
	if execute(args, func(string) string { return "" }, &printed) != nil {
		t.Fatal("local command")
	}
	got, _ := os.ReadFile(target)
	if !bytes.Equal(got, want) || !strings.Contains(printed.String(), "No deployment") {
		t.Fatal("wrong result")
	}
	if execute(args, func(string) string { return "" }, &printed) != errInvalid {
		t.Fatal("overwrite")
	}
	for _, a := range [][]string{{}, {"pack"}, {"pack", "--out", target}, {"pack", "--bad"}, {"pack", "--out", target, "extra"}, {"unknown", "--out", target}, {"stage", "--request", request, "--out", target}, {"stage", "--out", target}} {
		if execute(a, func(string) string { return "" }, io.Discard) == nil {
			t.Fatal("bad command")
		}
	}
	if _, e = read(root, 100); e == nil {
		t.Fatal("directory input")
	}
	if _, e = read(request, 1); e == nil {
		t.Fatal("oversize input")
	}
	if _, e = read(filepath.Join(root, "missing"), 10); e == nil {
		t.Fatal("missing input")
	}
}

// TestRequestAndSettings enforces canonical lossless identity and every opt-in workflow context field.
func TestRequestAndSettings(t *testing.T) {
	r, _, _ := fixture(t)
	body, _ := json.Marshal(r)
	if _, e := ParseStageRequest(body); e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{"{}", string(body[:len(body)-1]) + ",\"unknown\":true}", string(body[:len(body)-1]) + ",\"scope\":\"duplicate\"}", strings.Replace(string(body), "\"candidate_run_id\":\"123\"", "\"candidate_run_id\":123", 1), strings.Replace(string(body), "\"candidate_run_id\":\"123\"", "\"candidate_run_id\":\"0123\"", 1), strings.Replace(string(body), "\"sequence\":\"1\"", "\"sequence\":\"9223372036854775808\"", 1), strings.Replace(string(body), r.ID, r.CandidateID, 1), strings.Repeat(" ", 8193)} {
		if _, e := ParseStageRequest([]byte(raw)); e == nil {
			t.Fatal("bad request")
		}
	}
	env := map[string]string{"CMS_STAGE_ENABLED": "true", "GITHUB_REF": "refs/heads/main", "GITHUB_RUN_ATTEMPT": "1", "CMS_STAGE_ID": r.ID, "GITHUB_SHA": r.HeadSHA, "GITHUB_REPOSITORY": "owner/site", "GITHUB_REPOSITORY_ID": "42", "CMS_CANDIDATE_WORKFLOW_ID": "3", "GITHUB_TOKEN": strings.Repeat("x", 40)}
	if _, e := settings(func(k string) string { return env[k] }, r); e != nil {
		t.Fatal(e)
	}
	for key := range env {
		old := env[key]
		env[key] = "invalid"
		if _, e := settings(func(k string) string { return env[k] }, r); e == nil {
			t.Fatal("unsafe context", key)
		}
		env[key] = old
	}
}

// TestUnsafeArchives rejects source/payload drift, cancelled work and malformed ZIP structure before output.
func TestUnsafeArchives(t *testing.T) {
	for _, kind := range []string{"zip-pin", "payload-pin", "cancel", "short", "corrupt", "marker", "unsafe", "symlink", "directory", "long", "large"} {
		t.Run(kind, func(t *testing.T) {
			r, z, _ := fixture(t)
			ctx := context.Background()
			switch kind {
			case "zip-pin":
				r.ZIPDigest = strings.Repeat("0", 64)
			case "payload-pin":
				r.PayloadDigest = strings.Repeat("0", 64)
			case "cancel":
				c, stop := context.WithCancel(ctx)
				stop()
				ctx = c
			case "short":
				z = []byte("invalid")
				r.ZIPDigest = hash(z)
			case "corrupt":
				z[0] ^= 1
				r.ZIPDigest = hash(z)
			case "marker":
				r.Build.Sequence = "2"
			default:
				var b bytes.Buffer
				w := zip.NewWriter(&b)
				prior, _ := zip.NewReader(bytes.NewReader(z), int64(len(z)))
				for _, old := range prior.File {
					f, _ := w.Create(old.Name)
					reader, _ := old.Open()
					io.Copy(f, reader)
					reader.Close()
				}
				h := &zip.FileHeader{Name: "../private", Method: zip.Store}
				if kind == "symlink" {
					h.Name = "link"
					h.SetMode(os.ModeSymlink | 0777)
				}
				if kind == "directory" {
					h.Name = "empty/"
					h.SetMode(os.ModeDir | 0755)
				}
				if kind == "long" {
					h.Name = strings.Repeat("x", 101)
				}
				if kind == "large" {
					h.Name = "large.bin"
					h.Method = zip.Deflate
				}
				f, _ := w.CreateHeader(h)
				if kind == "large" {
					f.Write(make([]byte, MaxExpanded+1))
				} else {
					f.Write([]byte("bad"))
				}
				w.Close()
				z = b.Bytes()
				r.ZIPDigest = hash(z)
			}
			if _, e := pack(ctx, r, z); (e == nil) != (kind == "directory") {
				t.Fatal("archive policy", kind, e)
			}
		})
	}
}

// transport substitutes synthetic API/blob responses while production origin policy stays fixed.
type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// reply supplies a bounded response without logging arbitrary provider text.
func reply(status int, body []byte, kind string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(bytes.NewReader(body)), ContentLength: int64(len(body)), Header: http.Header{"Content-Type": []string{kind}}}
}

// TestProviderProtocol covers before/after provenance, private header isolation and malformed remote state.
func TestProviderProtocol(t *testing.T) {
	for _, kind := range []string{"valid", "rerun", "fork", "code", "workflow", "title", "pulls", "failed", "expired", "size", "url", "digest", "changed-run", "changed-artifact", "expired-after", "redirect-host", "redirect-status", "blob-error", "blob-status", "blob-mime", "blob-encoding", "blob-size", "json", "api-error"} {
		t.Run(kind, func(t *testing.T) {
			r, z, want := fixture(t)
			now := time.Now().UTC().Add(-time.Hour)
			runs, arts, blobs := 0, 0, 0
			c := &client{config: runtime{"owner/site", r.HeadSHA, strings.Repeat("x", 40), 42, 3}}
			c.api = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: transport(func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "api.github.com" || req.Method != "GET" || req.Header.Get("Authorization") != "Bearer "+c.config.token || req.Header.Get("Cookie") != "" {
					t.Fatal("API policy")
				}
				if kind == "api-error" {
					return nil, errors.New("synthetic private body")
				}
				var body any
				if strings.HasSuffix(req.URL.Path, "/zip") {
					out := reply(302, nil, "")
					out.Header.Set("Location", "https://fixture.blob.core.windows.net/file?sig=synthetic")
					if kind == "redirect-host" {
						out.Header.Set("Location", "https://127.0.0.1/private?sig=synthetic")
					}
					if kind == "redirect-status" {
						out.StatusCode = 200
					}
					return out, nil
				}
				if strings.Contains(req.URL.Path, "/runs/") {
					runs++
					repo := map[string]any{"id": 42, "full_name": "owner/site"}
					head := map[string]any{"id": 42, "full_name": "owner/site"}
					m := map[string]any{"id": 123, "run_attempt": 1, "workflow_id": 3, "head_sha": r.HeadSHA, "head_branch": "main", "event": "workflow_dispatch", "display_title": "CMS Astro candidate " + r.CandidateID, "path": ".github/workflows/cms-build-preview.yml", "status": "completed", "conclusion": "success", "created_at": now, "updated_at": now.Add(10 * time.Minute), "repository": repo, "head_repository": head, "pull_requests": []any{}}
					switch kind {
					case "rerun":
						m["run_attempt"] = 2
					case "fork":
						head["id"] = 99
					case "code":
						m["head_sha"] = strings.Repeat("b", 40)
					case "workflow":
						m["workflow_id"] = 9
					case "title":
						m["display_title"] = "other"
					case "pulls":
						delete(m, "pull_requests")
					case "failed":
						m["conclusion"] = "failure"
					case "changed-run":
						if runs > 1 {
							m["updated_at"] = now.Add(11 * time.Minute)
						}
					}
					body = m
				} else {
					arts++
					root := "https://api.github.com/repos/owner/site/actions/artifacts/456"
					m := map[string]any{"id": 456, "name": "astro-cms-candidate-" + r.CandidateID, "size_in_bytes": len(z), "digest": "sha256:" + r.ZIPDigest, "expired": false, "url": root, "archive_download_url": root + "/zip", "created_at": now.Add(5 * time.Minute), "updated_at": now.Add(5 * time.Minute), "expires_at": now.Add(24 * time.Hour), "workflow_run": map[string]any{"id": 123, "repository_id": 42, "head_repository_id": 42, "head_branch": "main", "head_sha": r.HeadSHA}}
					switch kind {
					case "expired":
						m["expired"] = true
					case "size":
						m["size_in_bytes"] = MaxArchive + 1
					case "url":
						m["url"] = "https://evil.invalid"
					case "digest":
						m["digest"] = "sha256:" + strings.Repeat("0", 64)
					case "changed-artifact":
						if arts > 1 {
							m["expires_at"] = now.Add(25 * time.Hour)
						}
					case "expired-after":
						if arts > 1 {
							m["expired"] = true
						}
					}
					body = m
				}
				raw, _ := json.Marshal(body)
				if kind == "json" {
					raw = []byte("bad")
				}
				return reply(200, raw, "application/json"), nil
			})}
			c.blob = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: transport(func(req *http.Request) (*http.Response, error) {
				blobs++
				if req.Header.Get("Authorization") != "" || req.Header.Get("Cookie") != "" || req.Header.Get("Referer") != "" || req.Method != "GET" {
					t.Fatal("credential forwarding")
				}
				out := reply(200, z, "application/zip")
				switch kind {
				case "blob-status":
					out.StatusCode = 302
				case "blob-mime":
					out.Header.Set("Content-Type", "text/html")
				case "blob-encoding":
					out.Header.Set("Content-Encoding", "gzip")
				case "blob-size":
					out.ContentLength++
				case "blob-error":
					return nil, errors.New("synthetic signed capability")
				}
				return out, nil
			})}
			got, e := c.prepare(context.Background(), r)
			if kind == "valid" {
				if e != nil || !bytes.Equal(got, want) || runs != 2 || arts != 2 || blobs != 1 {
					t.Fatal("valid protocol", e, runs, arts, blobs)
				}
			} else if e == nil {
				t.Fatal("unsafe provider", kind)
			}
		})
	}
}

// TestEgressPolicy checks signed hosts, public address classification and transport settings without external DNS/network.
func TestEgressPolicy(t *testing.T) {
	for _, url := range []string{"https://fixture.blob.core.windows.net/x?sig=synthetic", "https://results-receiver.actions.githubusercontent.com/x?sig=synthetic"} {
		if _, e := artifactLocation(url); e != nil {
			t.Fatal(e)
		}
	}
	for _, url := range []string{"http://fixture.blob.core.windows.net/x?sig=x", "https://token@fixture.blob.core.windows.net/x?sig=x", "https://fixture.blob.core.windows.net:443/x?sig=x", "https://fixture.blob.core.windows.net.evil/x?sig=x", "https://fixture.blob.core.windows.net/x", "https://fixture.blob.core.windows.net/x?sig=x#bad"} {
		if _, e := artifactLocation(url); e == nil {
			t.Fatal("unsafe URL")
		}
	}
	for _, s := range []string{"127.0.0.1", "10.0.0.1", "100.64.0.1", "192.0.2.1", "198.18.0.1", "198.51.100.1", "203.0.113.1", "240.0.0.1", "::1", "fc00::1", "fe80::1", "2001:db8::1", "2002::1"} {
		if publicArtifactIP(netip.MustParseAddr(s)) {
			t.Fatal("unsafe IP")
		}
	}
	if !publicArtifactIP(netip.MustParseAddr("8.8.8.8")) || !publicArtifactIP(netip.MustParseAddr("2606:4700:4700::1111")) {
		t.Fatal("public IP")
	}
	for _, addr := range []string{"bad", "example.invalid:80"} {
		if _, e := dialArtifact(context.Background(), "tcp", addr); e == nil {
			t.Fatal("bad dial")
		}
	}
	c := newClient(runtime{})
	defer c.api.CloseIdleConnections()
	defer c.blob.CloseIdleConnections()
	for _, h := range []*http.Client{c.api, c.blob} {
		tr := h.Transport.(*http.Transport)
		if tr.Proxy != nil || !tr.DisableCompression || h.Jar != nil || h.CheckRedirect(nil, nil) != http.ErrUseLastResponse {
			t.Fatal("unsafe transport")
		}
	}
}

// failingOutput ensures a failed local status write is reported instead of pretending upload success.
type failingOutput struct{}

func (failingOutput) Write([]byte) (int, error) { return 0, errors.New("synthetic output") }

// TestLocalFailures covers invalid request/ZIP/output paths and disabled valid staging input without egress.
func TestLocalFailures(t *testing.T) {
	r, z, _ := fixture(t)
	root := t.TempDir()
	request := filepath.Join(root, "request.json")
	archive := filepath.Join(root, "candidate.zip")
	body, _ := json.Marshal(r)
	os.WriteFile(request, body, 0600)
	os.WriteFile(archive, z, 0600)
	args := []string{"pack", "--request", request, "--zip", archive, "--out", filepath.Join(root, "new.tar")}
	if execute(args, func(string) string { return "" }, failingOutput{}) == nil {
		t.Fatal("status writer error")
	}
	for _, kind := range []string{"request", "zip", "pin", "parent"} {
		a := append([]string(nil), args...)
		a[6] = filepath.Join(t.TempDir(), "out.tar")
		switch kind {
		case "request":
			a[2] = filepath.Join(root, "absent")
		case "zip":
			a[4] = filepath.Join(root, "absent")
		case "pin":
			bad := r
			bad.PayloadDigest = strings.Repeat("0", 64)
			raw, _ := json.Marshal(bad)
			name := filepath.Join(t.TempDir(), "bad.json")
			os.WriteFile(name, raw, 0600)
			a[2] = name
		case "parent":
			a[6] = filepath.Join(root, "absent", "out.tar")
		}
		if execute(a, func(string) string { return "" }, io.Discard) == nil {
			t.Fatal("local failure hidden", kind)
		}
	}
	if execute([]string{"stage", "--out", filepath.Join(root, "stage.tar")}, func(k string) string {
		if k == "CMS_STAGE_REQUEST" {
			return string(body)
		}
		return ""
	}, io.Discard) == nil {
		t.Fatal("disabled staging")
	}
}

// TestPortablePolicy explicitly exercises cross-platform aliases before archive allocations or writes.
func TestPortablePolicy(t *testing.T) {
	for _, name := range []string{"", "/absolute", "../escape", "a/../b", "a\\b", "a//b", ".env", "a/.private", "a/CON.txt", "a/LPT9", "a/name.", "a/name ", strings.Repeat("a", 1025), "a/é.png", "a/" + strings.Repeat("x", 241)} {
		if portableMember(name) {
			t.Fatal("unsafe member")
		}
	}
	for _, name := range []string{".well-known/homepage-content.json", "images/demo.png", "a/console.txt", "path with spaces/file.txt"} {
		if !portableMember(name) {
			t.Fatal("portable member rejected")
		}
	}
}
