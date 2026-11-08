// Deployment command loads exact local files under explicit opt-in, never downloading caller-selected URLs.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	jsonv2 "encoding/json/v2"
	"io"
	"os"
	"strings"
	"time"
)

// deploymentPolicyFile is operator-reviewed configuration, not a public deployment request or authenticated grant.
type deploymentPolicyFile struct {
	SchemaVersion  int    `json:"schema_version"`
	Scope          string `json:"scope"`
	Repository     string `json:"repository"`
	RepositoryID   string `json:"repository_id"`
	WorkflowRef    string `json:"workflow_ref"`
	HeadSHA        string `json:"head_sha"`
	ReceiverOrigin string `json:"receiver_origin"`
	Audience       string `json:"audience"`
	IssuerOrigin   string `json:"issuer_origin"`
}

// readDeploymentFile refuses leaf links/non-files and verifies the opened file identity before bounded reading.
// Parent directories and concurrent writers remain trusted operator workspace inputs; digests/content checks follow.
func readDeploymentFile(path string, limit int64) ([]byte, error) {
	before, e := os.Lstat(path)
	if e != nil || !before.Mode().IsRegular() || before.Size() < 1 || before.Size() > limit {
		return nil, errReceiver
	}
	file, e := os.Open(path)
	if e != nil {
		return nil, errReceiver
	}
	defer file.Close()
	opened, e := file.Stat()
	if e != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) || opened.Size() != before.Size() {
		return nil, errReceiver
	}
	raw, e := io.ReadAll(io.LimitReader(file, limit+1))
	if e != nil || len(raw) == 0 || int64(len(raw)) > limit || int64(len(raw)) != opened.Size() {
		return nil, errReceiver
	}
	return raw, nil
}

// runDeploymentCommand exposes only production runtime construction; tests cannot select adapters through CLI arguments.
func runDeploymentCommand(args []string, get func(string) string, output io.Writer) error {
	return deploymentCommandWith(args, get, output, func(ctx context.Context, p deploymentPolicy, origin string, raw, candidate, uploaded []byte, get func(string) string) (senderOutcome, error) {
		runtime := deploymentRuntime{policy: p, issuerOrigin: origin}
		return runtime.run(ctx, raw, candidate, uploaded, get)
	})
}

// deploymentCommandWith permits a package-local runner fixture; no command flag enables this injection seam.
// Success writes only a constant report acknowledgement, never credentials, policy, paths or raw provider data.
func deploymentCommandWith(args []string, get func(string) string, output io.Writer, run func(context.Context, deploymentPolicy, string, []byte, []byte, []byte, func(string) string) (senderOutcome, error)) error {
	if get == nil || get("CMS_DEPLOY_ENABLED") != "true" || run == nil || output == nil || len(args) != 10 {
		return errReceiver
	}
	flags := map[string]string{}
	for i := 0; i < len(args); i += 2 {
		key, value := args[i], args[i+1]
		if key != "--policy" && key != "--policy-sha256" && key != "--request" && key != "--candidate" && key != "--uploaded" {
			return errReceiver
		}
		if flags[key] != "" || value == "" || strings.HasPrefix(value, "--") {
			return errReceiver
		}
		flags[key] = value
	}
	if !digest.MatchString(flags["--policy-sha256"]) {
		return errReceiver
	}
	raw, e := readDeploymentFile(flags["--policy"], 8192)
	if e != nil {
		return errReceiver
	}
	hash := sha256.Sum256(raw)
	if hex.EncodeToString(hash[:]) != flags["--policy-sha256"] {
		return errReceiver
	}
	var file deploymentPolicyFile
	if jsonv2.Unmarshal(raw, &file, jsonv2.RejectUnknownMembers(true)) != nil || file.SchemaVersion != 1 || file.Scope != "astro-protected-runtime-policy-v1" || file.IssuerOrigin == "" {
		return errReceiver
	}
	policy := deploymentPolicy{file.Repository, file.RepositoryID, file.WorkflowRef, file.HeadSHA, file.ReceiverOrigin, file.Audience}
	if _, e = validateDeploymentSettings(policy, get); e != nil {
		return errReceiver
	}
	request, e := readDeploymentFile(flags["--request"], 12288)
	if e != nil {
		return errReceiver
	}
	candidate, e := readDeploymentFile(flags["--candidate"], MaxArchive)
	if e != nil {
		return errReceiver
	}
	uploaded, e := readDeploymentFile(flags["--uploaded"], MaxDeployWrapper)
	if e != nil {
		return errReceiver
	}
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
	defer cancel()
	result, e := run(ctx, policy, file.IssuerOrigin, request, candidate, uploaded, get)
	if e != nil || result.Report.Scope != "astro-pages-operation-report-v1" {
		return errReceiver
	}
	_, e = io.WriteString(output, "Exact Pages operation reported. Completion and ordering hold remain unchanged.\n")
	if e != nil {
		return errReceiver
	}
	return nil
}
