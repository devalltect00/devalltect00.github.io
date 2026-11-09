// Deployment acquisition fixtures verify exact API/blob separation and full content without network or live credentials.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestDeploymentAcquisition rejects substituted provenance, stale metadata and altered archive bytes before any mutation.
func TestDeploymentAcquisition(t *testing.T) {
	for _, kind := range []string{"valid", "stage-run", "stage-workflow", "rerun", "fork", "code", "title", "live", "expired", "digest", "size", "redirect", "candidate-bytes", "uploaded-bytes", "changed-after", "cancelled"} {
		t.Run(kind, func(t *testing.T) {
			request, candidate, uploaded := receiverFixture(t)
			now := time.Now().UTC().Add(-time.Hour)
			source := &client{config: runtime{"example/homepage", request.Stage.HeadSHA, strings.Repeat("p", 40), 123, 22}}
			stageChecks := 0
			blobs := 0
			source.api = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: transport(func(req *http.Request) (*http.Response, error) {
				if req.Method != "GET" || req.URL.Host != "api.github.com" || req.Header.Get("Authorization") != "Bearer "+strings.Repeat("p", 40) {
					t.Error("API authority")
				}
				stage := strings.Contains(req.URL.Path, "/789") || strings.Contains(req.URL.Path, "/890")
				if strings.HasSuffix(req.URL.Path, "/zip") {
					out := reply(302, nil, "")
					name := "candidate"
					if stage {
						name = "uploaded"
					}
					out.Header.Set("Location", "https://fixture.blob.core.windows.net/"+name+"?sig=fixture")
					if kind == "redirect" {
						out.Header.Set("Location", "https://127.0.0.1/file?sig=fixture")
					}
					return out, nil
				}
				runID, artifactID, workflow, title, path, name, digest, size := int64(123), int64(456), int64(22), "CMS Astro candidate "+request.Stage.CandidateID, ".github/workflows/cms-build-preview.yml", "astro-cms-candidate-"+request.Stage.CandidateID, request.Stage.ZIPDigest, len(candidate)
				if stage {
					runID = 789
					artifactID = 890
					workflow = 33
					title = "CMS Astro staging " + request.Stage.ID
					path = ".github/workflows/cms-pages-stage.yml"
					name = "astro-cms-pages-" + request.Stage.ID
					digest = request.Evidence.ZIPDigest
					size = len(uploaded)
				}
				var body map[string]any
				if strings.Contains(req.URL.Path, "/runs/") {
					body = map[string]any{"id": runID, "run_attempt": 1, "workflow_id": workflow, "head_sha": request.Stage.HeadSHA, "head_branch": "main", "event": "workflow_dispatch", "display_title": title, "path": path, "status": "completed", "conclusion": "success", "created_at": now, "updated_at": now.Add(10 * time.Minute), "repository": map[string]any{"id": 123, "full_name": "example/homepage"}, "head_repository": map[string]any{"id": 123, "full_name": "example/homepage"}, "pull_requests": []any{}}
					if stage {
						stageChecks++
						switch kind {
						case "stage-run":
							body["id"] = 790
						case "stage-workflow":
							body["workflow_id"] = 34
						case "rerun":
							body["run_attempt"] = 2
						case "fork":
							body["head_repository"] = map[string]any{"id": 999, "full_name": "fork/homepage"}
						case "code":
							body["head_sha"] = strings.Repeat("0", 40)
						case "title":
							body["display_title"] = "other"
						case "live":
							body["status"] = "in_progress"
						case "changed-after":
							if stageChecks > 1 {
								body["updated_at"] = now.Add(11 * time.Minute)
							}
						}
					}
				} else {
					root := "https://api.github.com/repos/example/homepage/actions/artifacts/" + req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
					body = map[string]any{"id": artifactID, "name": name, "size_in_bytes": size, "digest": "sha256:" + digest, "expired": false, "url": root, "archive_download_url": root + "/zip", "created_at": now.Add(5 * time.Minute), "updated_at": now.Add(5 * time.Minute), "expires_at": now.Add(24 * time.Hour), "workflow_run": map[string]any{"id": runID, "repository_id": 123, "head_repository_id": 123, "head_branch": "main", "head_sha": request.Stage.HeadSHA}}
					if stage {
						switch kind {
						case "expired":
							body["expired"] = true
						case "digest":
							body["digest"] = "sha256:" + strings.Repeat("0", 64)
						case "size":
							body["size_in_bytes"] = MaxDeployWrapper + 1
						}
					}
				}
				raw, _ := json.Marshal(body)
				return reply(200, raw, "application/json"), nil
			})}
			source.blob = &http.Client{Transport: transport(func(req *http.Request) (*http.Response, error) {
				blobs++
				if req.Header.Get("Authorization") != "" || req.Header.Get("Cookie") != "" {
					t.Fatal("credential forwarded")
				}
				body := candidate
				if req.URL.Path == "/uploaded" {
					body = uploaded
				}
				body = append([]byte(nil), body...)
				if kind == "candidate-bytes" && req.URL.Path == "/candidate" || kind == "uploaded-bytes" && req.URL.Path == "/uploaded" {
					body[0] ^= 1
				}
				return reply(200, body, "application/zip"), nil
			})}
			acquire := deploymentAcquirer{source, 33}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "cancelled" {
				cancel()
			}
			c, u, e := acquire.acquire(ctx, request)
			if kind == "valid" {
				if e != nil || hash(c) != hash(candidate) || hash(u) != hash(uploaded) || blobs != 2 {
					t.Fatal("exact acquisition", e)
				}
			} else if e != errReceiver || c != nil || u != nil {
				t.Fatal("unsafe acquisition", e)
			}
		})
	}
}

// TestWorkflowCommandDefaults rejects execution before inputs/credentials and prevents partial source workflow configuration.
func TestWorkflowCommandDefaults(t *testing.T) {
	if execute([]string{"deploy-workflow"}, func(key string) string {
		if key != "CMS_DEPLOY_ENABLED" {
			t.Fatal("credential read")
		}
		return ""
	}, nil) != errReceiver {
		t.Fatal("not disabled")
	}
	for _, candidate := range []string{"", "022", "33"} {
		settings := deploymentSettings{RepositoryID: "123"}
		get := func(key string) string {
			switch key {
			case "CMS_CANDIDATE_WORKFLOW_ID":
				return candidate
			case "CMS_STAGE_WORKFLOW_ID":
				return "33"
			default:
				t.Fatal("token read before IDs")
				return ""
			}
		}
		if _, e := newDeploymentAcquirer(settings, get); e != errReceiver {
			t.Fatal("unsafe workflow IDs")
		}
	}
}
