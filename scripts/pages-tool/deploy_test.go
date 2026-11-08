// Receiver integrity tests use complete canonical local fixtures without provider authority or deployment.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// deployWrapper is fixture-only ZIP writing; production validation never writes or extracts archives.
func deployWrapper(t *testing.T, names []string, body []byte, mode os.FileMode) []byte {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	for _, name := range names {
		h := &zip.FileHeader{Name: name, Method: zip.Deflate}
		h.SetMode(mode)
		w, e := z.CreateHeader(h)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = w.Write(body); e != nil {
			t.Fatal(e)
		}
	}
	if e := z.Close(); e != nil {
		t.Fatal(e)
	}
	return b.Bytes()
}

// TestReceiverDeploymentInput checks strict wire input and independent canonical TAR/wrapper comparison.
func TestReceiverDeploymentInput(t *testing.T) {
	for _, kind := range []string{"valid", "scope", "duplicate-json", "unknown", "numeric-id", "reused-id", "payload", "candidate", "changed-tar", "truncated", "extra", "traversal", "link", "trailer", "pin", "tar-pin", "size", "cancelled"} {
		t.Run(kind, func(t *testing.T) {
			stage, candidate, tar := fixture(t)
			uploaded := deployWrapper(t, []string{"artifact.tar"}, tar, 0644)
			r := DeployRequest{SchemaVersion: 1, Scope: DeployScope, ID: "66666666-6666-4666-8666-666666666666", Stage: stage, RunID: "789", ArtifactID: "890", Evidence: StagedEvidence{ZIPDigest: hash(uploaded), PayloadDigest: stage.PayloadDigest, TARDigest: hash(tar), TARBytes: len(tar)}}
			switch kind {
			case "scope":
				r.Scope = "wrong"
			case "reused-id":
				r.ID = stage.ID
			case "payload":
				r.Evidence.PayloadDigest = strings.Repeat("0", 64)
			case "candidate":
				candidate[0] ^= 1
			case "changed-tar":
				copyTar := append([]byte{}, tar...)
				copyTar[1] ^= 1
				uploaded = deployWrapper(t, []string{"artifact.tar"}, copyTar, 0644)
				r.Evidence.ZIPDigest = hash(uploaded)
			case "truncated":
				uploaded = uploaded[:len(uploaded)-1]
				r.Evidence.ZIPDigest = hash(uploaded)
			case "extra":
				uploaded = deployWrapper(t, []string{"artifact.tar", "other"}, tar, 0644)
				r.Evidence.ZIPDigest = hash(uploaded)
			case "traversal":
				uploaded = deployWrapper(t, []string{"../artifact.tar"}, tar, 0644)
				r.Evidence.ZIPDigest = hash(uploaded)
			case "link":
				uploaded = deployWrapper(t, []string{"artifact.tar"}, tar, os.ModeSymlink|0644)
				r.Evidence.ZIPDigest = hash(uploaded)
			case "trailer":
				uploaded = append(uploaded, 0)
				r.Evidence.ZIPDigest = hash(uploaded)
			case "pin":
				r.Evidence.ZIPDigest = strings.Repeat("0", 64)
			case "tar-pin":
				r.Evidence.TARDigest = strings.Repeat("0", 64)
			case "size":
				r.Evidence.TARBytes++
			}
			raw, _ := json.Marshal(r)
			if kind == "duplicate-json" {
				raw = append([]byte(`{"scope":"wrong",`), raw[1:]...)
			}
			if kind == "unknown" {
				raw = append([]byte(`{"callback":"https://invalid.example",`), raw[1:]...)
			}
			if kind == "numeric-id" {
				raw = bytes.Replace(raw, []byte(`"stage_run_id":"789"`), []byte(`"stage_run_id":789`), 1)
			}
			parsed, e := ParseDeployRequest(raw)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if kind == "cancelled" {
				cancel()
			}
			if e == nil {
				e = verifyDeploymentInput(ctx, parsed, candidate, uploaded)
			}
			if (kind == "valid") != (e == nil) {
				t.Fatal("receiver integrity", e)
			}
		})
	}
}

// TestCompiledDeployReceiver independently checks CMS-generated intent and real compiled site archives.
func TestCompiledDeployReceiver(t *testing.T) {
	root := os.Getenv("CMS_DEPLOY_RECEIVER_FIXTURE")
	if root == "" {
		t.Skip("explicit compiled fixture required")
	}
	raw, e := read(root+"/deployment-request.json", 12288)
	if e != nil {
		t.Fatal(e)
	}
	r, e := ParseDeployRequest(raw)
	if e != nil {
		t.Fatal(e)
	}
	candidate, e := read(root+"/candidate.zip", MaxArchive)
	if e != nil {
		t.Fatal(e)
	}
	uploaded, e := read(root+"/uploaded.zip", MaxDeployWrapper)
	if e != nil {
		t.Fatal(e)
	}
	if e = verifyDeploymentInput(context.Background(), r, candidate, uploaded); e != nil {
		t.Fatal(e)
	}
}
