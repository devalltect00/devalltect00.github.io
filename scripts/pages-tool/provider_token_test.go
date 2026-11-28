// Synthetic opaque credentials protect compatibility and header safety; no real credential or network is used.
package main

import (
	"context"
	"strings"
	"testing"
)

// TestOpaqueProviderTokenBoundaries covers modern installation shapes and finite resource/injection limits.
func TestOpaqueProviderTokenBoundaries(t *testing.T) {
	cases := []struct {
		name, value string
		valid       bool
	}{
		{"legacy", strings.Repeat("a", 40), true},
		{"stateless-installation", "ghs_123456_" + strings.Repeat("a", 240) + "." + strings.Repeat("b", 240) + ".signature-with_dash", true},
		{"opaque-bearer", strings.Repeat("a", 32) + "._~+/=-", true},
		{"minimum", strings.Repeat("a", 32), true},
		{"maximum", strings.Repeat("a", 4096), true},
		{"empty", "", false}, {"short", strings.Repeat("a", 31), false}, {"over-limit", strings.Repeat("a", 4097), false},
		{"crlf", strings.Repeat("a", 32) + "\r\nAuthorization: other", false},
		{"space", strings.Repeat("a", 32) + " ", false}, {"tab", strings.Repeat("a", 32) + "\t", false},
		{"nul", strings.Repeat("a", 32) + "\x00", false}, {"delete", strings.Repeat("a", 32) + "\x7f", false},
		{"unicode", strings.Repeat("a", 32) + "é", false}, {"colon", strings.Repeat("a", 32) + ":", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if validProviderToken(tc.value) != tc.valid {
				t.Fatal("opaque token boundary mismatch")
			}
		})
	}
}

// TestOpaqueProviderTokenAcrossWorkflowClients rejects a regression in any of the three credential entry points.
func TestOpaqueProviderTokenAcrossWorkflowClients(t *testing.T) {
	request, _, _ := receiverFixture(t)
	plan := planFor(request, "999")
	long := "ghs_123456_" + strings.Repeat("a", 260) + "." + strings.Repeat("b", 260) + ".signature-with_dash"
	for _, token := range []string{strings.Repeat("a", 40), long, strings.Repeat("a", 4096), strings.Repeat("a", 4097), strings.Repeat("a", 40) + "\r\n", strings.Repeat("a", 40) + " "} {
		want := validProviderToken(token)
		env := map[string]string{"CMS_STAGE_ENABLED": "true", "GITHUB_REF": "refs/heads/main", "GITHUB_RUN_ATTEMPT": "1", "CMS_STAGE_ID": request.Stage.ID, "GITHUB_SHA": request.Stage.HeadSHA, "GITHUB_REPOSITORY": "example/homepage", "GITHUB_REPOSITORY_ID": "123", "CMS_CANDIDATE_WORKFLOW_ID": "22", "CMS_STAGE_WORKFLOW_ID": "33", "GITHUB_TOKEN": token}
		get := func(k string) string { return env[k] }
		if _, err := settings(get, request.Stage); (err == nil) != want {
			t.Fatal("staging credential compatibility")
		}
		acquirer, err := newDeploymentAcquirer(deploymentSettings{Repository: "example/homepage", RepositoryID: "123", Receiver: receiverConfig{HeadSHA: request.Stage.HeadSHA}}, get)
		if (err == nil) != want {
			t.Fatal("artifact acquisition credential compatibility")
		}
		if acquirer != nil {
			acquirer.close()
		}
		creator, err := newPagesCreator("example", "homepage", token, plan)
		if (err == nil) != want {
			t.Fatal("Pages creator credential compatibility")
		}
		if creator != nil {
			creator.Close()
		}
	}
	// An opaque value must reach only the fixed GitHub API header unchanged, never a decoded claim or public log.
	c := newClient(runtime{"example/homepage", request.Stage.HeadSHA, long, 123, 22})
	defer c.api.CloseIdleConnections()
	defer c.blob.CloseIdleConnections()
	req, err := c.request(context.Background(), "/actions/runs/123")
	if err != nil || req.Header.Get("Authorization") != "Bearer "+long || req.URL.Host != "api.github.com" {
		t.Fatal("opaque credential forwarding")
	}
}
