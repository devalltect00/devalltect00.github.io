// Workflow execution acquires exact source artifacts before invoking the guarded one-shot deployment runtime.
package main

import (
	"context"
	"io"
	"strings"
	"time"
)

// runWorkflowDeployment accepts reviewed local policy/request files; unsigned input never selects repository, workflow or URL.
// Only GET acquisition precedes full-byte validation; OIDC acquisition and Pages mutations follow the guarded runtime checks.
func runWorkflowDeployment(args []string, get func(string) string, output io.Writer) error {
	if get == nil || get("CMS_DEPLOY_ENABLED") != "true" || output == nil || len(args) != 6 {
		return errReceiver
	}
	flags := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		key, value := args[i], args[i+1]
		if key != "--policy" && key != "--policy-sha256" && key != "--request" || value == "" || strings.HasPrefix(value, "--") || flags[key] != "" {
			return errReceiver
		}
		flags[key] = value
	}
	policy, file, raw, e := loadDeploymentPolicy(flags, get)
	if e != nil {
		return errReceiver
	}
	request, e := ParseDeployRequest(raw)
	if e != nil {
		return errReceiver
	}
	settings, e := validateDeploymentSettings(policy, get)
	if e != nil || request.Stage.HeadSHA != settings.Receiver.HeadSHA || settings.Receiver.RunID == request.RunID || settings.Receiver.RunID == request.Stage.RunID {
		return errReceiver
	}
	ctx, cancel := context.WithTimeout(context.Background(), 135*time.Second)
	defer cancel()
	source, e := newDeploymentAcquirer(settings, get)
	if e != nil {
		return errReceiver
	}
	defer source.close()
	candidate, uploaded, e := source.acquire(ctx, request)
	if e != nil {
		return errReceiver
	}
	runtime := deploymentRuntime{policy: policy, issuerOrigin: file.IssuerOrigin}
	result, e := runtime.run(ctx, raw, candidate, uploaded, get)
	if e != nil || result.Report.Scope != "astro-pages-operation-report-v1" {
		return errReceiver
	}
	_, e = io.WriteString(output, "Exact Pages operation reported. Completion and ordering hold remain unchanged.\n")
	if e != nil {
		return errReceiver
	}
	return nil
}
