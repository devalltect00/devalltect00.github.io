// Deployment settings validate reviewed workflow context before credentials or network clients are used.
package main

import (
	"regexp"
	"strings"
)

// deploymentPolicy comes from reviewed protected configuration, never unsigned dispatch JSON.
// Environment observations are cross-checked against it; they do not attest GitHub identity by themselves.
type deploymentPolicy struct {
	Repository, RepositoryID, WorkflowRef, HeadSHA, ReceiverOrigin, Audience string
}

// deploymentSettings contains validated nonsecret context; it carries no token or deployment permission.
type deploymentSettings struct {
	Receiver                           receiverConfig
	Repository, RepositoryID, Audience string
}

// validateDeploymentSettings is side-effect free and never reads any credential environment variable.
// The caller must supply reviewed policy separately from the untrusted workflow/dispatch observations.
func validateDeploymentSettings(p deploymentPolicy, get func(string) string) (deploymentSettings, error) {
	fail := func() (deploymentSettings, error) { return deploymentSettings{}, errReceiver }
	if get == nil || !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,99}/[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`).MatchString(p.Repository) {
		return fail()
	}
	if _, e := positiveProviderID(p.RepositoryID); e != nil {
		return fail()
	}
	if p.WorkflowRef != p.Repository+"/.github/workflows/cms-pages-deploy.yml@refs/heads/main" || p.Audience != "https://github.com/"+strings.Split(p.Repository, "/")[0] {
		return fail()
	}
	expected := map[string]string{
		"CMS_DEPLOY_ENABLED": "true", "GITHUB_ACTIONS": "true", "GITHUB_EVENT_NAME": "workflow_dispatch",
		"GITHUB_SERVER_URL": "https://github.com", "GITHUB_API_URL": "https://api.github.com",
		"GITHUB_REPOSITORY": p.Repository, "GITHUB_REPOSITORY_ID": p.RepositoryID,
		"GITHUB_WORKFLOW_REF": p.WorkflowRef, "GITHUB_WORKFLOW_SHA": p.HeadSHA, "GITHUB_SHA": p.HeadSHA,
		"GITHUB_REF": "refs/heads/main", "GITHUB_RUN_ATTEMPT": "1", "RUNNER_ENVIRONMENT": "github-hosted",
		"CMS_DEPLOY_ENVIRONMENT": "github-pages", "CMS_RECEIVER_ORIGIN": p.ReceiverOrigin, "CMS_OIDC_AUDIENCE": p.Audience,
	}
	for key, value := range expected {
		if get(key) != value {
			return fail()
		}
	}
	receiver := receiverConfig{p.ReceiverOrigin, get("GITHUB_RUN_ID"), p.HeadSHA, "refs/heads/main", "1"}
	client, e := newReceiverClient(receiver)
	if e != nil {
		return fail()
	}
	client.close()
	return deploymentSettings{receiver, p.Repository, p.RepositoryID, p.Audience}, nil
}
