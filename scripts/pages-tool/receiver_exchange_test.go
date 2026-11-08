// Exchange transport tests use complete local archives and synthetic authority; no hosted endpoint is contacted.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestReceiverAdmissionReport checks exact plans, bounded wire behavior, local rejection and admission no-retry.
func TestReceiverAdmissionReport(t *testing.T) {
	for _, operation := range []string{"begin-write", "report"} {
		for _, kind := range []string{"valid", "scope", "attempt", "run", "artifact", "hash", "version", "expired", "future", "unknown", "duplicate", "null", "missing", "status", "redirect", "mime", "cache", "encoding", "large", "lost", "candidate", "upload", "token", "code", "source-run", "cancelled", "ack", "ack-duplicate", "ack-large", "decision"} {
			t.Run(operation+"/"+kind, func(t *testing.T) {
				r, candidate, uploaded := receiverFixture(t)
				var calls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					calls.Add(1)
					body, _ := io.ReadAll(req.Body)
					if req.Method != "POST" || req.URL.Path != "/api/v1/machine/astro/deployments/"+r.ID+"/"+operation || req.URL.RawQuery != "" || req.Header.Get("Authorization") != "Bearer one.two.three" || req.Header.Get("Cookie") != "" || req.Header.Get("Origin") != "" {
						t.Error("fixed machine transport")
					}
					if operation == "begin-write" && string(body) != "{}" {
						t.Error("admission input")
					}
					if operation == "report" && string(body) != `{"acknowledgement":{"id":"provider-fixture"}}` {
						t.Error("report input")
					}
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Cache-Control", "no-store")
					plan := planFor(r, "999")
					encoded, _ := json.Marshal(plan)
					d := map[string]any{}
					json.Unmarshal(encoded, &d)
					d["scope"] = "astro-pages-write-admission-v1"
					d["expires_at"] = time.Now().Add(5 * time.Second)
					if operation == "report" {
						delete(d, "expires_at")
						d["scope"] = "astro-pages-operation-report-v1"
						d["operation_id"] = plan.BuildVersion
						d["recorded_at"] = time.Now().Add(-time.Hour)
						d["replayed"] = calls.Load() > 1
					}
					switch kind {
					case "scope":
						d["scope"] = "bad"
					case "attempt":
						d["attempt_id"] = r.Stage.ID
					case "run":
						d["run_id"] = "1000"
					case "artifact":
						d["artifact_id"] = "891"
					case "hash":
						d["request_sha256"] = strings.Repeat("0", 64)
					case "version":
						d["build_version"] = strings.Repeat("0", 40)
					case "expired":
						if operation == "begin-write" {
							d["expires_at"] = time.Now().Add(-time.Second)
						} else {
							d["recorded_at"] = time.Time{}
						}
					case "future":
						if operation == "begin-write" {
							d["expires_at"] = time.Now().Add(time.Minute)
						} else {
							d["recorded_at"] = time.Now().Add(time.Hour)
						}
					case "unknown":
						d["callback"] = "https://invalid.example"
					case "duplicate":
						fmt.Fprint(w, `{"data":{},"data":{}}`)
						return
					case "null":
						if operation == "report" {
							d["replayed"] = nil
						} else {
							d["expires_at"] = nil
						}
					case "missing":
						if operation == "report" {
							delete(d, "replayed")
						} else {
							delete(d, "expires_at")
						}
					case "status":
						w.WriteHeader(429)
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
					case "lost":
						conn, _, _ := w.(http.Hijacker).Hijack()
						conn.Close()
						return
					}
					json.NewEncoder(w).Encode(map[string]any{"data": d})
				}))
				defer server.Close()
				client, e := receiverClientWithOrigin(receiverConfig{server.URL, "999", r.Stage.HeadSHA, "refs/heads/main", "1"}, true)
				if e != nil {
					t.Fatal(e)
				}
				defer client.close()
				x := receiverExchange{client: client}
				raw, _ := json.Marshal(r)
				decision := receiverDecision{r.ID, "999", hash(raw), time.Now().Add(20 * time.Second), "astro-receiver-authorization-v1"}
				token := "one.two.three"
				ack := []byte(`{"id":"provider-fixture"}`)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				local := false
				switch kind {
				case "candidate":
					candidate[0] ^= 1
					local = true
				case "upload":
					uploaded[0] ^= 1
					local = true
				case "token":
					token = "bad\r\nheader"
					local = true
				case "code":
					client.config.HeadSHA = strings.Repeat("b", 40)
					local = true
				case "source-run":
					client.config.RunID = r.RunID
					local = true
				case "cancelled":
					cancel()
					local = true
				case "ack":
					if operation == "report" {
						ack = []byte(`null`)
						local = true
					}
				case "ack-duplicate":
					if operation == "report" {
						ack = []byte(`{"id":"x","id":"y"}`)
						local = true
					}
				case "ack-large":
					if operation == "report" {
						ack = []byte(strings.Repeat("x", 8193))
						local = true
					}
				case "decision":
					if operation == "begin-write" {
						decision.ExpiresAt = time.Now().Add(-time.Second)
						local = true
					}
				}
				expectValid := kind == "valid" || operation == "begin-write" && strings.HasPrefix(kind, "ack") || operation == "report" && kind == "decision"
				if operation == "begin-write" {
					out, e := x.beginVerified(ctx, r, candidate, uploaded, token, decision)
					if (e == nil) != expectValid {
						t.Fatalf("admission outcome %v", e)
					}
					if e == nil && (!out.validFor(r, "999", time.Now()) || out.validFor(r, "999", out.ExpiresAt)) {
						t.Fatal("permit validity")
					}
					if _, e = x.beginVerified(ctx, r, candidate, uploaded, token, decision); e != errReceiver {
						t.Fatal("admission resent")
					}
				} else {
					out, e := x.reportVerified(ctx, r, candidate, uploaded, token, ack)
					if (e == nil) != expectValid {
						t.Fatalf("report outcome %v", e)
					}
					if kind == "valid" {
						if out.Replayed {
							t.Fatal("first report replay")
						}
						out, e = x.reportVerified(ctx, r, candidate, uploaded, token, ack)
						if e != nil || !out.Replayed {
							t.Fatal("historical exact report replay", e)
						}
					}
				}
				want := int32(1)
				if local {
					want = 0
				}
				if operation == "report" && kind == "valid" {
					want = 2
				}
				if calls.Load() != want {
					t.Fatal("unexpected egress", calls.Load(), want)
				}
			})
		}
	}
}

// TestReceiverAdmissionConcurrent proves competing local callers cannot send two admission requests.
func TestReceiverAdmissionConcurrent(t *testing.T) {
	r, candidate, uploaded := receiverFixture(t)
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		json.NewEncoder(w).Encode(struct {
			Data receiverPermit `json:"data"`
		}{receiverPermit{planFor(r, "999"), "astro-pages-write-admission-v1", time.Now().Add(5 * time.Second)}})
	}))
	defer server.Close()
	client, e := receiverClientWithOrigin(receiverConfig{server.URL, "999", r.Stage.HeadSHA, "refs/heads/main", "1"}, true)
	if e != nil {
		t.Fatal(e)
	}
	defer client.close()
	x := receiverExchange{client: client}
	raw, _ := json.Marshal(r)
	decision := receiverDecision{r.ID, "999", hash(raw), time.Now().Add(20 * time.Second), "astro-receiver-authorization-v1"}
	var wg sync.WaitGroup
	var success atomic.Int32
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, e := x.beginVerified(context.Background(), r, candidate, uploaded, "one.two.three", decision); e == nil {
				success.Add(1)
			}
		}()
	}
	wg.Wait()
	if calls.Load() != 1 || success.Load() != 1 {
		t.Fatal("concurrent admission duplicated", calls.Load(), success.Load())
	}
}
