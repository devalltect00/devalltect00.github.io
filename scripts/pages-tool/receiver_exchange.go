// Receiver exchange clients validate admission and durable reports without exposing runtime or Pages writes.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"io"
	"mime"
	"net/http"
	"regexp"
	"sync/atomic"
	"time"
)

// receiverPlan binds lossless provider IDs to the exact typed request and reviewed deployment run.
type receiverPlan struct {
	AttemptID    string `json:"attempt_id"`
	RunID        string `json:"run_id"`
	ArtifactID   string `json:"artifact_id"`
	RequestHash  string `json:"request_sha256"`
	BuildVersion string `json:"build_version"`
}

// receiverPermit is a short-lived server admission; local validation never replaces durable server consumption.
type receiverPermit struct {
	receiverPlan
	Scope     string    `json:"scope"`
	ExpiresAt time.Time `json:"expires_at"`
}

// receiverReport is historical acknowledgement evidence, never permission to write or finish a job.
type receiverReport struct {
	receiverPlan
	Scope       string    `json:"scope"`
	OperationID string    `json:"operation_id"`
	RecordedAt  time.Time `json:"recorded_at"`
	Replayed    bool      `json:"replayed"`
}

// receiverExchange shares the reviewed fixed transport and spends admission before its first network attempt.
// Report is explicitly callable again with the same acknowledgement; it never invokes admission or Pages.
type receiverExchange struct {
	client   *receiverClient
	admitted atomic.Bool
}

// planFor derives the provider operation version independently using the private/public protocol formula.
func planFor(r DeployRequest, run string) receiverPlan {
	raw, _ := json.Marshal(r)
	digest := hash(raw)
	return receiverPlan{r.ID, run, r.ArtifactID, digest, hash([]byte("devalltect.astro.pages-operation.v1\n" + run + "\n" + digest))[:40]}
}

// validFor checks every plan field and the live ten-second deadline immediately before future protected use.
func (p receiverPermit) validFor(r DeployRequest, run string, now time.Time) bool {
	return p.receiverPlan == planFor(r, run) && p.Scope == "astro-pages-write-admission-v1" && p.ExpiresAt.After(now) && !p.ExpiresAt.After(now.Add(10*time.Second))
}

// beginVerified verifies complete archives and the consumed authorization before one irreversible admission POST.
// The original bearer must be supplied; the server, not a local token claim, enforces its consumed identity.
func (x *receiverExchange) beginVerified(ctx context.Context, r DeployRequest, candidate, uploaded []byte, token string, decision receiverDecision) (receiverPermit, error) {
	var result receiverPermit
	if x.client == nil || !x.admitted.CompareAndSwap(false, true) {
		return result, errReceiver
	}
	if !decision.validFor(r, x.client.config.RunID, time.Now()) || !x.localValid(ctx, r, candidate, uploaded, token) {
		return result, errReceiver
	}
	if x.post(ctx, r.ID, "begin-write", token, []byte("{}"), &result) != nil || !result.validFor(r, x.client.config.RunID, time.Now()) || result.ExpiresAt.After(decision.ExpiresAt) {
		return receiverPermit{}, errReceiver
	}
	return result, nil
}

// reportVerified forwards only a bounded acknowledgement object and verifies the exact historical receipt.
// Retrying this report is explicit and cannot repeat a provider write; server authentication/binding remain authoritative.
func (x *receiverExchange) reportVerified(ctx context.Context, r DeployRequest, candidate, uploaded []byte, token string, acknowledgement []byte) (receiverReport, error) {
	var result receiverReport
	var object map[string]json.RawMessage
	if len(acknowledgement) > 8192 || jsonv2.Unmarshal(acknowledgement, &object) != nil || object == nil || !x.localValid(ctx, r, candidate, uploaded, token) {
		return result, errReceiver
	}
	body, e := json.Marshal(struct {
		Acknowledgement json.RawMessage `json:"acknowledgement"`
	}{acknowledgement})
	if e != nil || len(body) > 8256 {
		return result, errReceiver
	}
	if x.post(ctx, r.ID, "report", token, body, &result) != nil || result.receiverPlan != planFor(r, x.client.config.RunID) || result.Scope != "astro-pages-operation-report-v1" || result.OperationID != result.BuildVersion || result.RecordedAt.IsZero() || result.RecordedAt.After(time.Now().Add(2*time.Minute)) {
		return receiverReport{}, errReceiver
	}
	return result, nil
}

// localValid rejects invalid content/context before egress; public content is never authentication authority.
func (x *receiverExchange) localValid(ctx context.Context, r DeployRequest, candidate, uploaded []byte, token string) bool {
	if x.client == nil || len(token) == 0 || len(token) > 16384 || !regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`).MatchString(token) || r.Stage.HeadSHA != x.client.config.HeadSHA || x.client.config.RunID == r.RunID || x.client.config.RunID == r.Stage.RunID {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	return verifyDeploymentInput(ctx, r, candidate, uploaded) == nil && ctx.Err() == nil
}

// post uses the existing fixed origin and bounded transport; it never retries or returns remote errors/bearers.
func (x *receiverExchange) post(ctx context.Context, id, operation, token string, body []byte, out any) error {
	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	req, e := http.NewRequestWithContext(ctx, "POST", x.client.config.Origin+"/api/v1/machine/astro/deployments/"+id+"/"+operation, bytes.NewReader(body))
	if e != nil {
		return errReceiver
	}
	req.GetBody = nil // Disable transport body replay, including a stale reused connection.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	response, e := x.client.http.Do(req)
	if e != nil {
		return errReceiver
	}
	defer response.Body.Close()
	typ, _, e := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if e != nil || typ != "application/json" || response.StatusCode != 200 || response.ContentLength > 4096 || response.Header.Get("Content-Encoding") != "" || response.Header.Get("Cache-Control") != "no-store" {
		return errReceiver
	}
	raw, e := io.ReadAll(io.LimitReader(response.Body, 4097))
	if e != nil || len(raw) == 0 || len(raw) > 4096 {
		return errReceiver
	}
	envelope := struct {
		Data json.RawMessage `json:"data"`
	}{}
	if jsonv2.Unmarshal(raw, &envelope, jsonv2.RejectUnknownMembers(true)) != nil || len(envelope.Data) == 0 || bytes.TrimSpace(envelope.Data)[0] != '{' || jsonv2.Unmarshal(envelope.Data, out, jsonv2.RejectUnknownMembers(true)) != nil || ctx.Err() != nil {
		return errReceiver
	}
	var fields map[string]json.RawMessage
	if jsonv2.Unmarshal(envelope.Data, &fields) != nil {
		return errReceiver
	}
	required := []string{"scope", "attempt_id", "run_id", "artifact_id", "request_sha256", "build_version"}
	if operation == "begin-write" {
		required = append(required, "expires_at")
	} else {
		required = append(required, "operation_id", "recorded_at", "replayed")
	}
	for _, key := range required {
		if len(fields[key]) == 0 || bytes.Equal(bytes.TrimSpace(fields[key]), []byte("null")) {
			return errReceiver
		}
	}
	return nil
}
