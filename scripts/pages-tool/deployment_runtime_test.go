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

// TestDeploymentRuntimeExchange joins guarded configuration, real local issuance, admission, creation and reporting.
func TestDeploymentRuntimeExchange(t *testing.T) {
	for _, kind := range []string{"valid", "issuance", "candidate", "upload", "cancelled", "authorize", "admission", "expired", "foreign-plan", "status", "redirect", "lost", "mime", "encoding", "large", "ack-id", "ack-url", "ack-duplicate", "fresh-token", "report", "report-plan"} {
		t.Run(kind, func(t *testing.T) {
			r, candidate, uploaded := receiverFixture(t)
			plan := planFor(r, "999")
			var issued atomic.Int32
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

			issuerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if issued.Add(1) == 1 {
					if kind == "issuance" {
						w.WriteHeader(503)
						return
					}
					fmt.Fprint(w, `{"value":"original.fixture.token"}`)
				} else {
					fresh.Add(1)
					if kind == "fresh-token" {
						w.WriteHeader(503)
						return
					}
					fmt.Fprint(w, `{"value":"fresh.fixture.token"}`)
				}
			}))
			defer issuerServer.Close()
			issuer, e := workflowOIDCAt(issuerServer.URL, issuerServer.URL+"/issue", "fixture-credential", "https://github.com/example", true)
			if e != nil {
				t.Fatal(e)
			}
			policy, values := runtimeFixture(r)
			runtime := deploymentRuntime{policy: policy, issuerOrigin: "https://runner.actions.githubusercontent.com"}
			builds := 0
			factory := func(p deploymentPolicy, o string, settings deploymentSettings, actual receiverPlan, get func(string) string) (deploymentClients, error) {
				builds++
				if p != policy || actual != plan || settings.Repository != policy.Repository {
					t.Error("reviewed settings changed")
				}
				return deploymentClients{issuer, receiver, creator}, nil
			}
			read := func(key string) string {
				if strings.Contains(key, "TOKEN") {
					t.Fatal("credential access outside trusted factory")
				}
				return values[key]
			}
			raw, _ := json.Marshal(r)

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
			out, e := runtime.runWith(ctx, raw, candidate, uploaded, read, factory)
			if (e == nil) != (kind == "valid") {
				t.Fatalf("unexpected outcome %v", e)
			}
			if (builds == 0) != (kind == "candidate" || kind == "upload" || kind == "cancelled") {
				t.Fatal("factory before verified content")
			}
			before := creates.Load()
			if _, e = runtime.runWith(ctx, raw, candidate, uploaded, read, factory); e != errReceiver || creates.Load() != before {
				t.Fatal("sender replayed")
			}
			wantCreate := int32(1)
			if kind == "issuance" || kind == "candidate" || kind == "upload" || kind == "cancelled" || kind == "authorize" || kind == "admission" || kind == "expired" || kind == "foreign-plan" {
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

// runtimeFixture supplies reviewed nonsecret policy and exact workflow observations for local tests.
func runtimeFixture(r DeployRequest) (deploymentPolicy, map[string]string) {
	p := deploymentPolicy{"example/homepage", "123", "example/homepage/.github/workflows/cms-pages-deploy.yml@refs/heads/main", r.Stage.HeadSHA, "https://cms.example.test", "https://github.com/example"}
	values := map[string]string{"CMS_DEPLOY_ENABLED": "true", "GITHUB_ACTIONS": "true", "GITHUB_EVENT_NAME": "workflow_dispatch", "GITHUB_SERVER_URL": "https://github.com", "GITHUB_API_URL": "https://api.github.com", "GITHUB_REPOSITORY": p.Repository, "GITHUB_REPOSITORY_ID": p.RepositoryID, "GITHUB_WORKFLOW_REF": p.WorkflowRef, "GITHUB_WORKFLOW_SHA": p.HeadSHA, "GITHUB_SHA": p.HeadSHA, "GITHUB_REF": "refs/heads/main", "GITHUB_RUN_ATTEMPT": "1", "RUNNER_ENVIRONMENT": "github-hosted", "CMS_DEPLOY_ENVIRONMENT": "github-pages", "CMS_RECEIVER_ORIGIN": p.ReceiverOrigin, "CMS_OIDC_AUDIENCE": p.Audience, "GITHUB_RUN_ID": "999"}
	return p, values
}

// TestDeploymentRuntimePreflight checks bad policy/request/code/run are rejected without factory or credential access.
func TestDeploymentRuntimePreflight(t *testing.T) {
	for _, kind := range []string{"policy", "request", "code", "source-run", "disabled"} {
		t.Run(kind, func(t *testing.T) {
			r, candidate, uploaded := receiverFixture(t)
			p, values := runtimeFixture(r)
			switch kind {
			case "policy":
				p.RepositoryID = "bad"
			case "code":
				r.Stage.HeadSHA = strings.Repeat("b", 40)
			case "source-run":
				values["GITHUB_RUN_ID"] = r.RunID
			case "disabled":
				values["CMS_DEPLOY_ENABLED"] = "false"
			}
			raw, _ := json.Marshal(r)
			if kind == "request" {
				raw = []byte(`{"unknown":true}`)
			}
			runtime := deploymentRuntime{policy: p}
			get := func(key string) string {
				if strings.Contains(key, "TOKEN") {
					t.Fatal("credential read")
				}
				return values[key]
			}
			build := func(deploymentPolicy, string, deploymentSettings, receiverPlan, func(string) string) (deploymentClients, error) {
				t.Fatal("factory reached")
				return deploymentClients{}, errReceiver
			}
			if _, e := runtime.runWith(context.Background(), raw, candidate, uploaded, get, build); e != errReceiver {
				t.Fatal("invalid preflight")
			}
		})
	}
}

// TestRuntimeProductionConstructors exercises production construction/cleanup without performing any network operation.
func TestRuntimeProductionConstructors(t *testing.T) {
	for _, kind := range []string{"valid", "policy", "issuer", "receiver", "provider-token"} {
		t.Run(kind, func(t *testing.T) {
			r, _, _ := receiverFixture(t)
			p, values := runtimeFixture(r)
			settings, e := validateDeploymentSettings(p, func(k string) string { return values[k] })
			if e != nil {
				t.Fatal(e)
			}
			origin := "https://runner.actions.githubusercontent.com"
			values["ACTIONS_ID_TOKEN_REQUEST_URL"] = origin + "/issue?api-version=2.0"
			values["ACTIONS_ID_TOKEN_REQUEST_TOKEN"] = "fixture-runner-credential"
			values["GITHUB_TOKEN"] = strings.Repeat("p", 40)
			switch kind {
			case "policy":
				p.RepositoryID = "bad"
			case "issuer":
				origin = "https://invalid.example"
			case "receiver":
				settings.Receiver.Origin = "http://invalid.example"
			case "provider-token":
				values["GITHUB_TOKEN"] = "bad"
			}
			clients, e := runtimeClients(p, origin, settings, planFor(r, "999"), func(k string) string { return values[k] })
			defer clients.close()
			if (e == nil) != (kind == "valid") {
				t.Fatal("constructor outcome", e)
			}
		})
	}
	runtime := deploymentRuntime{}
	if _, e := runtime.run(context.Background(), nil, nil, nil, nil); e != errReceiver {
		t.Fatal("unguarded production entry")
	}
}
