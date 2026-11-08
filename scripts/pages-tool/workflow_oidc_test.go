// Issuance tests use synthetic loopback tokens and never request live credentials.
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// TestWorkflowOIDCAcquisition bounds issuance, enforces fixed origin/query/header and rejects unsafe replies.
func TestWorkflowOIDCAcquisition(t *testing.T) {
	for _, kind := range []string{"valid", "status", "redirect", "mime", "encoding", "large", "duplicate", "unknown", "null", "empty", "malformed", "token-large", "cancelled", "lost"} {
		t.Run(kind, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != "GET" || r.URL.Path != "/issue" || r.URL.Query().Get("audience") != "https://github.com/example" || r.URL.Query().Get("api-version") != "2.0" || r.Header.Get("Authorization") != "Bearer fixture-runner-credential" || r.Header.Get("Cookie") != "" || r.Header.Get("Origin") != "" {
					t.Error("issuance transport")
				}
				w.Header().Set("Content-Type", "application/json")
				switch kind {
				case "status":
					w.WriteHeader(429)
				case "redirect":
					http.Redirect(w, r, "/other", 302)
					return
				case "mime":
					w.Header().Set("Content-Type", "text/plain")
				case "encoding":
					w.Header().Set("Content-Encoding", "gzip")
				case "large":
					fmt.Fprint(w, strings.Repeat("x", 32769))
					return
				case "duplicate":
					fmt.Fprint(w, `{"value":"a.b.c","value":"a.b.c"}`)
					return
				case "unknown":
					fmt.Fprint(w, `{"value":"a.b.c","other":true}`)
					return
				case "null":
					fmt.Fprint(w, `{"value":null}`)
					return
				case "empty":
					fmt.Fprint(w, `{}`)
					return
				case "malformed":
					fmt.Fprint(w, `{"value":"bad"}`)
					return
				case "token-large":
					fmt.Fprintf(w, `{"value":"%s.b.c"}`, strings.Repeat("a", 16385))
					return
				case "lost":
					conn, _, _ := w.(http.Hijacker).Hijack()
					conn.Close()
					return
				}
				fmt.Fprint(w, `{"value":"fixture.payload.signature"}`)
			}))
			defer server.Close()
			client, e := workflowOIDCAt(server.URL, server.URL+"/issue?api-version=2.0", "fixture-runner-credential", "https://github.com/example", true)
			if e != nil {
				t.Fatal(e)
			}
			defer client.close()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "cancelled" {
				cancel()
			}
			token, e := client.acquire(ctx)
			if (e == nil) != (kind == "valid") {
				t.Fatal("unsafe issuance outcome", e)
			}
			if e != nil && token != "" {
				t.Fatal("token returned on failure")
			}
			_, _ = client.acquire(ctx)
			if _, e := client.acquire(ctx); e != errReceiver {
				t.Fatal("budget reset")
			}
			want := int32(2)
			if kind == "cancelled" {
				want = 0
			}
			if calls.Load() != want {
				t.Fatal("issuance retries", calls.Load())
			}
		})
	}
}

// TestWorkflowOIDCPolicy rejects endpoint substitution, ambiguous audiences and credential injection before egress.
func TestWorkflowOIDCPolicy(t *testing.T) {
	origin := "https://runner.actions.githubusercontent.com"
	endpoint := origin + "/issue?api-version=2.0"
	c, e := workflowOIDCAt(origin, endpoint, "fixture-token", "https://github.com/example", false)
	if e != nil {
		t.Fatal(e)
	}
	defer c.close()
	for _, bad := range []string{"https://evil.example/issue", origin + ".evil.example/issue", origin + "/issue?audience=other", origin + "/issue#fragment", "http://runner.actions.githubusercontent.com/issue", "https://user@runner.actions.githubusercontent.com/issue"} {
		if _, e := workflowOIDCAt(origin, bad, "fixture-token", "https://github.com/example", false); e != errReceiver {
			t.Fatal("unsafe endpoint")
		}
	}
	if _, e := workflowOIDCAt(origin, endpoint, "bad\r\nheader", "https://github.com/example", false); e != errReceiver {
		t.Fatal("credential injection")
	}
	read := func(key string) string {
		if strings.Contains(key, "TOKEN") {
			t.Fatal("credential read before guard")
		}
		return ""
	}
	if _, e := newWorkflowOIDC(deploymentPolicy{}, origin, read); e != errReceiver {
		t.Fatal("guard bypass")
	}
}
