// Local command tests use real bounded input files and a package-local runner fixture, never external services.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDeploymentCommand validates argument/file/policy safety and keeps output constant after successful reporting.
func TestDeploymentCommand(t *testing.T) {
	for _, kind := range []string{"valid", "disabled", "unknown-flag", "duplicate-flag", "missing-flag", "digest", "unknown-policy", "duplicate-policy", "scope", "settings", "missing-file", "directory", "empty", "oversize", "bad-candidate", "bad-upload", "runner-error", "bad-report", "output-error"} {
		t.Run(kind, func(t *testing.T) {
			r, candidate, uploaded := receiverFixture(t)
			policy, values := runtimeFixture(r)
			file := deploymentPolicyFile{1, "astro-protected-runtime-policy-v1", policy.Repository, policy.RepositoryID, policy.WorkflowRef, policy.HeadSHA, policy.ReceiverOrigin, policy.Audience, "https://runner.actions.githubusercontent.com"}
			raw, _ := json.Marshal(file)
			request, _ := json.Marshal(r)
			switch kind {
			case "unknown-policy":
				raw = append(raw[:len(raw)-1], []byte(`,"extra":true}`)...)
			case "duplicate-policy":
				raw = append(raw[:len(raw)-1], []byte(`,"scope":"other"}`)...)
			case "scope":
				file.Scope = "wrong"
				raw, _ = json.Marshal(file)
			case "settings":
				values["GITHUB_RUN_ATTEMPT"] = "2"
			case "disabled":
				values["CMS_DEPLOY_ENABLED"] = "false"
			}
			dir := t.TempDir()
			paths := map[string]string{}
			for name, data := range map[string][]byte{"policy": raw, "request": request, "candidate": candidate, "uploaded": uploaded} {
				path := filepath.Join(dir, name)
				if e := os.WriteFile(path, data, 0600); e != nil {
					t.Fatal(e)
				}
				paths[name] = path
			}
			args := []string{"--policy", paths["policy"], "--policy-sha256", hash(raw), "--request", paths["request"], "--candidate", paths["candidate"], "--uploaded", paths["uploaded"]}
			switch kind {
			case "unknown-flag":
				args[0] = "--url"
			case "duplicate-flag":
				args[6] = "--request"
			case "missing-flag":
				args = args[:8]
			case "digest":
				args[3] = strings.Repeat("0", 64)
			case "missing-file":
				os.Remove(paths["uploaded"])
			case "directory":
				os.Remove(paths["uploaded"])
				os.Mkdir(paths["uploaded"], 0700)
			case "empty":
				os.WriteFile(paths["uploaded"], nil, 0600)
			case "oversize":
				f, _ := os.OpenFile(paths["uploaded"], os.O_WRONLY, 0600)
				f.Truncate(MaxDeployWrapper + 1)
				f.Close()
			case "bad-candidate":
				candidate[0] ^= 1
				os.WriteFile(paths["candidate"], candidate, 0600)
			case "bad-upload":
				uploaded[0] ^= 1
				os.WriteFile(paths["uploaded"], uploaded, 0600)
			}
			var buffer bytes.Buffer
			var output io.Writer = &buffer
			if kind == "output-error" {
				output = brokenDeploymentOutput{}
			}
			calls := 0
			read := func(key string) string {
				if strings.Contains(key, "TOKEN") {
					t.Fatal("credential read before runtime")
				}
				return values[key]
			}
			run := func(ctx context.Context, p deploymentPolicy, origin string, req, c, u []byte, get func(string) string) (senderOutcome, error) {
				calls++
				if p != policy || origin != file.IssuerOrigin {
					t.Fatal("policy changed")
				}
				if _, ok := ctx.Deadline(); !ok {
					t.Fatal("missing deadline")
				}
				parsed, e := ParseDeployRequest(req)
				if e != nil || verifyDeploymentInput(ctx, parsed, c, u) != nil {
					return senderOutcome{}, errReceiver
				}
				if kind == "runner-error" {
					return senderOutcome{}, errors.New("private upstream body must never escape")
				}
				scope := "astro-pages-operation-report-v1"
				if kind == "bad-report" {
					scope = "wrong"
				}
				return senderOutcome{Report: receiverReport{Scope: scope}}, nil
			}
			e := deploymentCommandWith(args, read, output, run)
			if (e == nil) != (kind == "valid") {
				t.Fatal("unsafe command outcome", e)
			}
			if e != nil && e != errReceiver {
				t.Fatal("raw diagnostic escaped")
			}
			if kind == "valid" && !strings.HasPrefix(buffer.String(), "Exact Pages operation reported.") {
				t.Fatal("success output")
			}
			if kind != "valid" && kind != "output-error" && buffer.Len() != 0 {
				t.Fatal("failure output leak")
			}
			expected := kind == "valid" || kind == "bad-candidate" || kind == "bad-upload" || kind == "runner-error" || kind == "bad-report" || kind == "output-error"
			if (calls == 1) != expected {
				t.Fatal("runner reached before preflight", calls)
			}
		})
	}
}

// brokenDeploymentOutput exercises a caller output failure without echoing provider acknowledgement data.
type brokenDeploymentOutput struct{}

// Write always fails so the command must emit its constant error and never repeat deployment work.
func (brokenDeploymentOutput) Write([]byte) (int, error) {
	return 0, errors.New("fixture output failure")
}

// TestDeploymentCommandDisabled checks main routing rejects defaults without touching local inputs or credentials.
func TestDeploymentCommandDisabled(t *testing.T) {
	if execute([]string{"deploy-local"}, func(key string) string {
		if key != "CMS_DEPLOY_ENABLED" {
			t.Fatal("unexpected read")
		}
		return ""
	}, io.Discard) != errReceiver {
		t.Fatal("not disabled by default")
	}
}
