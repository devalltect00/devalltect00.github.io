// Sender fixtures join real local HTTP clients with complete archives, never live provider credentials.
package main

import (
	"context"
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

// TestProtectedSenderExchange proves ordered local authority exchanges and credential separation across origins.
func TestProtectedSenderExchange(t *testing.T) {
	for _, kind := range []string{"valid", "candidate", "upload", "cancelled", "authorize", "admission", "expired", "foreign-plan", "status", "redirect", "lost", "mime", "encoding", "large", "ack-id", "ack-url", "ack-duplicate", "fresh-token", "report", "report-plan"} {
		t.Run(kind, func(t *testing.T) {
			r, candidate, uploaded := receiverFixture(t)
			plan := planFor(r, "999")
			var phase atomic.Int32
			var creates atomic.Int32
			var reports atomic.Int32
			var fresh atomic.Int32
			cms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				raw, _ := io.ReadAll(req.Body)
				if req.Method != "POST" || req.URL.RawQuery != "" || req.Header.Get("Origin") != "" || req.Header.Get("Cookie") != "" {
					t.Error("machine request")
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				d := map[string]any{}
				b, _ := json.Marshal(plan)
				json.Unmarshal(b, &d)
				switch req.URL.Path {
				case "/api/v1/machine/astro/deployments/" + r.ID + "/authorize":
					if !phase.CompareAndSwap(0, 1) || string(raw) != "{}" || req.Header.Get("Authorization") != "Bearer original.fixture.token" {
						t.Error("authorization sequence/credential")
					}
					if kind == "authorize" {
						w.WriteHeader(403)
						return
					}
					delete(d, "artifact_id")
					delete(d, "build_version")
					d["scope"] = "astro-receiver-authorization-v1"
					d["expires_at"] = time.Now().Add(25 * time.Second)
				case "/api/v1/machine/astro/deployments/" + r.ID + "/begin-write":
					if !phase.CompareAndSwap(1, 2) || string(raw) != "{}" || req.Header.Get("Authorization") != "Bearer original.fixture.token" {
						t.Error("admission sequence/credential")
					}
					if kind == "admission" {
						w.WriteHeader(403)
						return
					}
					d["scope"] = "astro-pages-write-admission-v1"
					d["expires_at"] = time.Now().Add(8 * time.Second)
					if kind == "expired" {
						d["expires_at"] = time.Now().Add(-time.Second)
					}
					if kind == "foreign-plan" {
						d["artifact_id"] = "891"
					}
				case "/api/v1/machine/astro/deployments/" + r.ID + "/report":
					reports.Add(1)
					if phase.Load() != 3 || req.Header.Get("Authorization") != "Bearer fresh.fixture.token" {
						t.Error("report sequence/credential")
					}
					var envelope struct {
						Acknowledgement json.RawMessage `json:"acknowledgement"`
					}
					if json.Unmarshal(raw, &envelope) != nil || !strings.Contains(string(envelope.Acknowledgement), plan.BuildVersion) {
						t.Error("actual acknowledgement missing")
					}
					if kind == "report" && reports.Load() == 1 {
						w.WriteHeader(503)
						return
					}
					d["scope"] = "astro-pages-operation-report-v1"
					d["operation_id"] = plan.BuildVersion
					d["recorded_at"] = time.Now().Add(-time.Second)
					d["replayed"] = reports.Load() > 1
					if kind == "report-plan" {
						d["run_id"] = "1000"
					}
				default:
					t.Error("unexpected receiver path")
					w.WriteHeader(404)
					return
				}
				json.NewEncoder(w).Encode(map[string]any{"data": d})
			}))
			defer cms.Close()
			provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				creates.Add(1)
				if !phase.CompareAndSwap(2, 3) || req.Method != "POST" || req.URL.Path != "/repos/example/homepage/pages/deployments" || req.URL.RawQuery != "" || req.Header.Get("Authorization") != "Bearer "+strings.Repeat("p", 40) || req.Header.Get("Origin") != "" || req.Header.Get("Cookie") != "" {
					t.Error("provider sequence/path/credential")
				}
				var body map[string]json.RawMessage
				json.NewDecoder(req.Body).Decode(&body)
				if len(body) != 4 || string(body["artifact_id"]) != r.ArtifactID || string(body["environment"]) != `"github-pages"` || string(body["pages_build_version"]) != `"`+plan.BuildVersion+`"` || string(body["oidc_token"]) != `"original.fixture.token"` {
					t.Error("exact provider body")
				}
				w.Header().Set("Content-Type", "application/json")
				ack := map[string]any{"id": plan.BuildVersion, "status_url": "https://api.github.com/repos/example/homepage/pages/deployments/" + plan.BuildVersion + "/status"}
				switch kind {
				case "status":
					w.WriteHeader(429)
				case "redirect":
					http.Redirect(w, req, "/other", 302)
					return
				case "lost":
					conn, _, _ := w.(http.Hijacker).Hijack()
					conn.Close()
					return
				case "mime":
					w.Header().Set("Content-Type", "text/plain")
				case "encoding":
					w.Header().Set("Content-Encoding", "gzip")
				case "large":
					fmt.Fprint(w, strings.Repeat("x", 8193))
					return
				case "ack-id":
					ack["id"] = r.Stage.HeadSHA
				case "ack-url":
					ack["status_url"] = "https://invalid.example/status"
				case "ack-duplicate":
					fmt.Fprint(w, `{"id":"x","id":"y"}`)
					return
				}
				json.NewEncoder(w).Encode(ack)
			}))
			defer provider.Close()
			receiver, e := receiverClientWithOrigin(receiverConfig{cms.URL, "999", r.Stage.HeadSHA, "refs/heads/main", "1"}, true)
			if e != nil {
				t.Fatal(e)
			}
			defer receiver.close()
			creator, e := pagesCreatorAt("example", "homepage", strings.Repeat("p", 40), plan, provider.URL)
			if e != nil {
				t.Fatal(e)
			}
			defer creator.Close()
			sender := protectedSender{receiver: receiver, creator: creator}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "candidate" {
				candidate[0] ^= 1
			}
			if kind == "upload" {
				uploaded[0] ^= 1
			}
			if kind == "cancelled" {
				cancel()
			}
			out, e := sender.send(ctx, r, candidate, uploaded, "original.fixture.token", func(context.Context) (string, error) {
				fresh.Add(1)
				if kind == "fresh-token" {
					return "", errReceiver
				}
				return "fresh.fixture.token", nil
			})
			if (e == nil) != (kind == "valid") {
				t.Fatalf("unexpected outcome %v", e)
			}
			before := creates.Load()
			if _, e = sender.send(ctx, r, candidate, uploaded, "original.fixture.token", nil); e != errReceiver || creates.Load() != before {
				t.Fatal("sender replayed")
			}
			wantCreate := int32(1)
			if kind == "candidate" || kind == "upload" || kind == "cancelled" || kind == "authorize" || kind == "admission" || kind == "expired" || kind == "foreign-plan" {
				wantCreate = 0
			}
			if creates.Load() != wantCreate {
				t.Fatal("create count", creates.Load(), wantCreate)
			}
			validAck := kind == "valid" || kind == "fresh-token" || kind == "report" || kind == "report-plan"
			if (len(out.Acknowledgement) > 0) != validAck || (fresh.Load() == 1) != validAck {
				t.Fatal("ack/token callback boundary")
			}
			if kind == "report" {
				exchange := receiverExchange{client: receiver}
				receipt, e := exchange.reportVerified(ctx, r, candidate, uploaded, "fresh.fixture.token", out.Acknowledgement)
				if e != nil || !receipt.Replayed || creates.Load() != 1 {
					t.Fatal("report-only recovery", e)
				}
			}
		})
	}
}

// TestPagesCreatorConfiguration rejects arbitrary origins, repository paths, credentials and noncanonical plans.
func TestPagesCreatorConfiguration(t *testing.T) {
	r, _, _ := receiverFixture(t)
	plan := planFor(r, "999")
	c, e := newPagesCreator("example", "homepage", strings.Repeat("p", 40), plan)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	for _, origin := range []string{"https://invalid.example", "http://api.github.com", "https://api.github.com/path", "https://user@api.github.com", "https://api.github.com?"} {
		if _, e := pagesCreatorAt("example", "homepage", strings.Repeat("p", 40), plan, origin); e != errReceiver {
			t.Fatal("unsafe origin")
		}
	}
	for _, field := range []string{"owner", "repo", "token", "run", "artifact", "hash", "version"} {
		p := plan
		owner, repo, token := "example", "homepage", strings.Repeat("p", 40)
		switch field {
		case "owner":
			owner = "../other"
		case "repo":
			repo = "../other"
		case "token":
			token = "bad\r\n"
		case "run":
			p.RunID = "0999"
		case "artifact":
			p.ArtifactID = "0"
		case "hash":
			p.RequestHash = "bad"
		case "version":
			p.BuildVersion = r.Stage.HeadSHA
		}
		if _, e := newPagesCreator(owner, repo, token, p); e != errReceiver {
			t.Fatal("unsafe configuration", field)
		}
	}
	ctx := context.Background()
	if _, e := c.create(ctx, plan, "original.fixture.token"); e != errReceiver {
		t.Fatal("deadline required")
	}
}
