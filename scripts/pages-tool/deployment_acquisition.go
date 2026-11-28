// Deployment artifact acquisition uses exact run/artifact IDs under one reviewed repository and never extracts downloaded content.
package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// deploymentAcquirer reuses credential-separated GitHub/blob transports with distinct reviewed source workflow IDs.
type deploymentAcquirer struct {
	source        *client
	stageWorkflow int64
}

// newDeploymentAcquirer accepts workflow-read credentials only after reviewed configuration and request checks.
func newDeploymentAcquirer(settings deploymentSettings, get func(string) string) (*deploymentAcquirer, error) {
	candidate, e := positiveProviderID(get("CMS_CANDIDATE_WORKFLOW_ID"))
	stage, f := positiveProviderID(get("CMS_STAGE_WORKFLOW_ID"))
	repo, g := positiveProviderID(settings.RepositoryID)
	if e != nil || f != nil || g != nil || candidate == stage {
		return nil, errReceiver
	}
	token := get("GITHUB_TOKEN")
	if !validProviderToken(token) {
		return nil, errReceiver
	}
	return &deploymentAcquirer{newClient(runtime{settings.Repository, settings.Receiver.HeadSHA, token, repo, candidate}), stage}, nil
}

// close releases both credentialed API and credential-free signed-download sockets.
func (c *deploymentAcquirer) close() {
	c.source.api.CloseIdleConnections()
	c.source.blob.CloseIdleConnections()
}

// inspectUploaded checks exact successful first staging run, unexpired artifact and complete retained evidence pins.
func (c *deploymentAcquirer) inspectUploaded(ctx context.Context, r DeployRequest, budget *int) (run, artifact, error) {
	var x run
	var a artifact
	root := "https://api.github.com/repos/" + c.source.config.repo + "/actions/artifacts/" + r.ArtifactID
	if c.source.read(ctx, "/actions/runs/"+r.RunID, &x, budget) != nil || c.source.read(ctx, "/actions/artifacts/"+r.ArtifactID, &a, budget) != nil {
		return x, a, errReceiver
	}
	runID, _ := positiveProviderID(r.RunID)
	artifactID, _ := positiveProviderID(r.ArtifactID)
	now := time.Now()
	config := c.source.config
	if x.ID != runID || x.Attempt != 1 || x.Workflow != c.stageWorkflow || x.SHA != config.head || x.SHA != r.Stage.HeadSHA || x.Branch != "main" || x.Event != "workflow_dispatch" || x.Title != "CMS Astro staging "+r.Stage.ID || x.Path != ".github/workflows/cms-pages-stage.yml" || x.Status != "completed" || x.Conclusion != "success" || x.Repository.ID != config.repoID || x.HeadRepo.ID != config.repoID || !strings.EqualFold(x.Repository.Name, config.repo) || !strings.EqualFold(x.HeadRepo.Name, config.repo) || x.Pulls == nil || len(*x.Pulls) != 0 || x.Created.IsZero() || x.Updated.Before(x.Created) || x.Updated.After(now.Add(2*time.Minute)) {
		return x, a, errReceiver
	}
	if a.ID != artifactID || a.Name != "astro-cms-pages-"+r.Stage.ID || a.Size < 1 || a.Size > MaxDeployWrapper || a.Digest != "sha256:"+r.Evidence.ZIPDigest || a.Expired == nil || *a.Expired || a.URL != root || a.Archive != root+"/zip" || !a.Expires.After(now) || a.Run.ID != runID || a.Run.Repo != config.repoID || a.Run.HeadRepo != config.repoID || a.Run.Branch != "main" || a.Run.SHA != config.head || a.Created.Before(x.Created.Add(-2*time.Minute)) || a.Created.After(x.Updated.Add(2*time.Minute)) || a.Updated.Before(a.Created) || a.Updated.After(now.Add(2*time.Minute)) {
		return x, a, errReceiver
	}
	return x, a, nil
}

// acquire rechecks metadata around both exact raw downloads and complete independent content validation.
func (c *deploymentAcquirer) acquire(ctx context.Context, r DeployRequest) ([]byte, []byte, error) {
	fail := func() ([]byte, []byte, error) { return nil, nil, errReceiver }
	raw, _ := json.Marshal(r)
	if _, e := ParseDeployRequest(raw); e != nil || c == nil || c.source == nil || c.stageWorkflow < 1 || c.stageWorkflow == c.source.config.workflowID {
		return fail()
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if ctx.Err() != nil {
		return fail()
	}
	budget := 0
	cr, ca, e := c.source.inspect(ctx, r.Stage, &budget)
	if e != nil {
		return fail()
	}
	sr, sa, e := c.inspectUploaded(ctx, r, &budget)
	if e != nil {
		return fail()
	}
	candidate, e := c.source.download(ctx, "/actions/artifacts/"+r.Stage.ArtifactID+"/zip", ca.Size)
	if e != nil {
		return fail()
	}
	uploaded, e := c.source.download(ctx, "/actions/artifacts/"+r.ArtifactID+"/zip", sa.Size)
	if e != nil || verifyDeploymentInput(ctx, r, candidate, uploaded) != nil {
		return fail()
	}
	cr2, ca2, e := c.source.inspect(ctx, r.Stage, &budget)
	if e != nil {
		return fail()
	}
	sr2, sa2, e := c.inspectUploaded(ctx, r, &budget)
	if e != nil {
		return fail()
	}
	before, _ := json.Marshal([]any{cr, ca, sr, sa})
	after, _ := json.Marshal([]any{cr2, ca2, sr2, sa2})
	if string(before) != string(after) || ctx.Err() != nil {
		return fail()
	}
	return candidate, uploaded, nil
}
