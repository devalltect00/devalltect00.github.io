// Staging requests identify exact public bytes; they cannot authorize deployment or finish a CMS job.
package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"regexp"
	"strconv"
)

// StageScope versions the optional read-and-upload workflow handoff.
const StageScope = "astro-pages-stage-v1"

// StageRequest binds a fresh staging correlation to one candidate, reviewed code and exact Pages payload.
type StageRequest struct {
	SchemaVersion int     `json:"schema_version"`
	Scope         string  `json:"scope"`
	ID            string  `json:"request_id"`
	CandidateID   string  `json:"candidate_attempt_id"`
	RunID         string  `json:"candidate_run_id"`
	ArtifactID    string  `json:"candidate_artifact_id"`
	HeadSHA       string  `json:"head_sha"`
	ZIPDigest     string  `json:"zip_sha256"`
	PayloadDigest string  `json:"payload_sha256"`
	Build         Request `json:"build"`
}

// ParseStageRequest rejects unknown/duplicate data, noncanonical identities and precision-losing IDs.
func ParseStageRequest(raw []byte) (StageRequest, error) {
	var r StageRequest
	if len(raw) < 1 || len(raw) > 8192 || jsonv2.Unmarshal(raw, &r, jsonv2.RejectUnknownMembers(true)) != nil {
		return r, errInvalid
	}
	b, _ := json.Marshal(r.Build)
	_, e := ParseRequest(b)
	positive := func(v string) bool {
		n, e := strconv.ParseInt(v, 10, 64)
		return e == nil && n > 0 && strconv.FormatInt(n, 10) == v
	}
	if e != nil || r.SchemaVersion != 1 || r.Scope != StageScope || !uuid.MatchString(r.ID) || !uuid.MatchString(r.CandidateID) || r.ID == r.CandidateID || !positive(r.RunID) || !positive(r.ArtifactID) || !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(r.HeadSHA) || !digest.MatchString(r.ZIPDigest) || !digest.MatchString(r.PayloadDigest) {
		return r, errInvalid
	}
	return r, nil
}

// Request is frozen public identity, never a credential or current CMS authority.
type Request struct {
	SchemaVersion    int    `json:"schema_version"`
	Scope            string `json:"scope"`
	ReleaseID        string `json:"release_id"`
	Sequence         string `json:"sequence"`
	SnapshotDigest   string `json:"snapshot_digest"`
	ProjectionDigest string `json:"projection_digest"`
}

var uuid = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
var digest = regexp.MustCompile(`^[a-f0-9]{64}$`)

// ParseRequest strictly checks the independently defined public identity.
func ParseRequest(raw []byte) (Request, error) {
	var r Request
	if jsonv2.Unmarshal(raw, &r, jsonv2.RejectUnknownMembers(true)) != nil {
		return r, errInvalid
	}
	n, e := strconv.ParseInt(r.Sequence, 10, 64)
	if e != nil || n < 1 || strconv.FormatInt(n, 10) != r.Sequence || r.SchemaVersion != 1 || r.Scope != "astro-project-profile-build-v1" || !uuid.MatchString(r.ReleaseID) || !digest.MatchString(r.SnapshotDigest) || !digest.MatchString(r.ProjectionDigest) {
		return r, errInvalid
	}
	return r, nil
}
