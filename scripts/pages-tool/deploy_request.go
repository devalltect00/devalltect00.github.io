// Deployment intents pin exact already-verified bytes; unsigned intent is never deployment authority.
package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"strconv"
)

// MaxDeployTar bounds complete canonical TAR bytes with candidate header padding.
const MaxDeployTar = MaxExpanded + MaxEntries*1024 + 1024

// MaxDeployWrapper bounds a one-member Actions wrapper before allocation.
const MaxDeployWrapper = MaxDeployTar + (1 << 20)

// DeployScope versions the separate deployment handoff without changing staging semantics.
const DeployScope = "astro-pages-deploy-v1"

// DeployRequest binds one new deployment intent to exact source, staging identity and full uploaded-byte evidence.
// It deliberately accepts no token, callback, destination URL, environment or workflow supplied by a client.
type DeployRequest struct {
	SchemaVersion int            `json:"schema_version"`
	Scope         string         `json:"scope"`
	ID            string         `json:"request_id"`
	Stage         StageRequest   `json:"stage"`
	RunID         string         `json:"stage_run_id"`
	ArtifactID    string         `json:"stage_artifact_id"`
	Evidence      StagedEvidence `json:"uploaded"`
}

// ParseDeployRequest strictly rejects oversized, unknown, duplicate, noncanonical or inconsistent public input.
func ParseDeployRequest(raw []byte) (DeployRequest, error) {
	var r DeployRequest
	if len(raw) == 0 || len(raw) > 12288 || jsonv2.Unmarshal(raw, &r, jsonv2.RejectUnknownMembers(true)) != nil {
		return DeployRequest{}, errInvalid
	}
	b, _ := json.Marshal(r.Stage)
	if _, e := ParseStageRequest(b); e != nil {
		return DeployRequest{}, errInvalid
	}
	positive := func(s string) bool {
		n, e := strconv.ParseInt(s, 10, 64)
		return e == nil && n > 0 && strconv.FormatInt(n, 10) == s
	}
	if r.SchemaVersion != 1 || r.Scope != DeployScope || !uuid.MatchString(r.ID) || r.ID == r.Stage.ID || r.ID == r.Stage.CandidateID || !positive(r.RunID) || !positive(r.ArtifactID) || r.RunID == r.Stage.RunID || !digest.MatchString(r.Evidence.ZIPDigest) || !digest.MatchString(r.Evidence.TARDigest) || r.Evidence.PayloadDigest != r.Stage.PayloadDigest || r.Evidence.TARBytes <= 0 || r.Evidence.TARBytes > MaxDeployTar {
		return DeployRequest{}, errInvalid
	}
	return r, nil
}

// StagedEvidence binds wrapper, canonical payload and TAR bytes; no field grants deployment authority.
type StagedEvidence struct {
	ZIPDigest     string `json:"zip_sha256"`
	PayloadDigest string `json:"payload_sha256"`
	TARDigest     string `json:"tar_sha256"`
	TARBytes      int    `json:"tar_bytes"`
}
