// Receiver exchange fixtures use synthetic loopback authority, complete archives and no live bearer or Pages calls.
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// receiverFixture builds real independently reconstructible bytes; the HTTP decision is still synthetic authority.
func receiverFixture(t *testing.T) (DeployRequest, []byte, []byte) {
	t.Helper()
	stage, candidate, tar := fixture(t)
	uploaded := deployWrapper(t, []string{"artifact.tar"}, tar, 0644)
	return DeployRequest{SchemaVersion: 1, Scope: DeployScope, ID: "66666666-6666-4666-8666-666666666666", Stage: stage, RunID: "789", ArtifactID: "890", Evidence: StagedEvidence{ZIPDigest: hash(uploaded), PayloadDigest: stage.PayloadDigest, TARDigest: hash(tar), TARBytes: len(tar)}}, candidate, uploaded
}

// TestReceiverExchange checks exact pinned response, no retries/redirects and zero egress for invalid local content.
func TestReceiverExchange(t *testing.T) {
	for _, kind := range []string{"valid", "scope", "attempt", "run", "hash", "expired", "future", "unknown", "duplicate", "null", "status", "redirect", "mime", "cache", "encoding", "large", "candidate", "upload", "token", "code", "source-run", "cancelled"} {
		t.Run(kind, func(t *testing.T) {
			r, candidate, uploaded := receiverFixture(t)
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls.Add(1)
				body, _ := io.ReadAll(req.Body)
				if req.Method != "POST" || req.URL.Path != "/api/v1/machine/astro/deployments/"+r.ID+"/authorize" || req.URL.RawQuery != "" || req.Header.Get("Authorization") != "Bearer one.two.three" || req.Header.Get("Cookie") != "" || req.Header.Get("Origin") != "" || string(body) != "{}" {
					t.Error("fixed machine transport")
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				raw, _ := json.Marshal(r)
				d := map[string]any{"scope": "astro-receiver-authorization-v1", "attempt_id": r.ID, "run_id": "999", "request_sha256": hash(raw), "expires_at": time.Now().Add(20 * time.Second).UTC().Format(time.RFC3339Nano)}
				switch kind {
				case "scope":
					d["scope"] = "wrong"
				case "attempt":
					d["attempt_id"] = r.Stage.ID
				case "run":
					d["run_id"] = "1000"
				case "hash":
					d["request_sha256"] = strings.Repeat("0", 64)
				case "expired":
					d["expires_at"] = time.Now().Add(-time.Second).Format(time.RFC3339Nano)
				case "future":
					d["expires_at"] = time.Now().Add(time.Minute).Format(time.RFC3339Nano)
				case "unknown":
					d["callback"] = "https://invalid.example"
				case "duplicate":
					fmt.Fprint(w, `{"data":null,"data":null}`)
					return
				case "null":
					fmt.Fprint(w, `{"data":null}`)
					return
				case "status":
					w.WriteHeader(403)
				case "redirect":
					http.Redirect(w, req, "/other", 302)
					return
				case "mime":
					w.Header().Set("Content-Type", "text/plain")
				case "cache":
					w.Header().Del("Cache-Control")
				case "encoding":
					w.Header().Set("Content-Encoding", "gzip")
				case "large":
					fmt.Fprint(w, strings.Repeat("x", 4097))
					return
				}
				json.NewEncoder(w).Encode(map[string]any{"data": d})
			}))
			defer server.Close()
			cfg := receiverConfig{Origin: server.URL, RunID: "999", HeadSHA: r.Stage.HeadSHA, Ref: "refs/heads/main", Attempt: "1"}
			if kind == "code" {
				cfg.HeadSHA = strings.Repeat("b", 40)
			}
			if kind == "source-run" {
				cfg.RunID = r.RunID
			}
			client, e := receiverClientWithOrigin(cfg, true)
			if e != nil {
				t.Fatal(e)
			}
			defer client.close()
			token := "one.two.three"
			if kind == "token" {
				token = "bad\r\nAuthorization: other"
			}
			if kind == "candidate" {
				candidate[0] ^= 1
			}
			if kind == "upload" {
				uploaded[0] ^= 1
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "cancelled" {
				cancel()
			}
			out, e := client.authorizeVerified(ctx, r, candidate, uploaded, token)
			if kind == "valid" {
				if e != nil || out.RunID != "999" || !out.validFor(r, "999", time.Now()) {
					t.Fatal(out, e)
				}
				if out.validFor(r, "999", out.ExpiresAt) {
					t.Fatal("expired decision accepted")
				}
				if _, e = client.authorizeVerified(ctx, r, candidate, uploaded, token); e != errReceiver {
					t.Fatal("client resent")
				}
			} else if e != errReceiver {
				t.Fatal("unsafe receiver decision", out, e)
			}
			want := int32(1)
			if kind == "candidate" || kind == "upload" || kind == "token" || kind == "code" || kind == "source-run" || kind == "cancelled" {
				want = 0
			}
			if calls.Load() != want {
				t.Fatal("unexpected egress", calls.Load(), want)
			}
		})
	}
}

// TestReceiverConfiguration ensures no user-origin, weak workflow context or transport override can enable egress.
func TestReceiverConfiguration(t *testing.T) {
	cfg := receiverConfig{Origin: "https://cms.example.test", RunID: "999", HeadSHA: strings.Repeat("a", 40), Ref: "refs/heads/main", Attempt: "1"}
	c, e := newReceiverClient(cfg)
	if e != nil {
		t.Fatal(e)
	}
	defer c.close()
	transport := c.http.Transport.(*http.Transport)
	if transport.Proxy != nil || transport.TLSClientConfig.MinVersion < tls.VersionTLS12 || !transport.DisableCompression || c.http.Jar != nil {
		t.Fatal("transport policy")
	}
	for _, origin := range []string{"http://cms.example.test", "https://user@cms.example.test", "https://cms.example.test/path", "https://cms.example.test?token=x", "https://cms.example.test#fragment", "https://cms.example.test:8443", "https://127.0.0.1"} {
		copy := cfg
		copy.Origin = origin
		if _, e := newReceiverClient(copy); e != errReceiver {
			t.Fatal("unsafe origin")
		}
	}
	copy := cfg
	copy.Ref = "refs/heads/other"
	if _, e := newReceiverClient(copy); e != errReceiver {
		t.Fatal("foreign ref")
	}
	copy = cfg
	copy.Attempt = "2"
	if _, e := newReceiverClient(copy); e != errReceiver {
		t.Fatal("rerun")
	}
	copy = cfg
	copy.RunID = "01"
	if _, e := newReceiverClient(copy); e != errReceiver {
		t.Fatal("noncanonical run")
	}
}
