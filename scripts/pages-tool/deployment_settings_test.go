// Deployment configuration tests exercise all reviewed pins without requesting credentials or accessing providers.
package main

import (
	"strings"
	"testing"
)

// TestDeploymentSettings rejects every substituted observation and malformed policy before credential access.
func TestDeploymentSettings(t *testing.T) {
	p := deploymentPolicy{"example/homepage", "123", "example/homepage/.github/workflows/cms-pages-deploy.yml@refs/heads/main", strings.Repeat("a", 40), "https://cms.example.test", "https://github.com/example"}
	values := map[string]string{"CMS_DEPLOY_ENABLED": "true", "GITHUB_ACTIONS": "true", "GITHUB_EVENT_NAME": "workflow_dispatch", "GITHUB_SERVER_URL": "https://github.com", "GITHUB_API_URL": "https://api.github.com", "GITHUB_REPOSITORY": p.Repository, "GITHUB_REPOSITORY_ID": p.RepositoryID, "GITHUB_WORKFLOW_REF": p.WorkflowRef, "GITHUB_WORKFLOW_SHA": p.HeadSHA, "GITHUB_SHA": p.HeadSHA, "GITHUB_REF": "refs/heads/main", "GITHUB_RUN_ATTEMPT": "1", "RUNNER_ENVIRONMENT": "github-hosted", "CMS_DEPLOY_ENVIRONMENT": "github-pages", "CMS_RECEIVER_ORIGIN": p.ReceiverOrigin, "CMS_OIDC_AUDIENCE": p.Audience, "GITHUB_RUN_ID": "999"}
	read := func(key string) string {
		if strings.Contains(key, "TOKEN") || strings.Contains(key, "PASSWORD") {
			t.Fatal("credential read")
		}
		return values[key]
	}
	out, e := validateDeploymentSettings(p, read)
	if e != nil || out.Receiver.RunID != "999" {
		t.Fatal("valid settings", e)
	}
	for key, original := range values {
		t.Run(key, func(t *testing.T) {
			values[key] = "substituted"
			defer func() { values[key] = original }()
			if _, e := validateDeploymentSettings(p, read); e != errReceiver {
				t.Fatal("substitution accepted")
			}
		})
	}
	for _, kind := range []string{"repository", "id", "workflow", "code", "origin", "audience", "nil"} {
		t.Run("policy/"+kind, func(t *testing.T) {
			copy := p
			get := read
			switch kind {
			case "repository":
				copy.Repository = "../other"
			case "id":
				copy.RepositoryID = "0123"
			case "workflow":
				copy.WorkflowRef = "other"
			case "code":
				copy.HeadSHA = "bad"
			case "origin":
				copy.ReceiverOrigin = "http://cms.example.test"
			case "audience":
				copy.Audience = "other"
			case "nil":
				get = nil
			}
			if _, e := validateDeploymentSettings(copy, get); e != errReceiver {
				t.Fatal("invalid policy")
			}
		})
	}
}
