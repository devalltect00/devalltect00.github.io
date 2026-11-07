// Fixed GETs inspect exact first-run provenance around acquisition and immutable payload reconstruction.
package main

import (
	"context"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"io"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// runtime is selected by trusted workflow context, never by caller repository/URL fields.
type runtime struct {
	repo, head, token  string
	repoID, workflowID int64
}

// client isolates API authentication and signed download policy; tests inject synthetic transports.
type client struct {
	config    runtime
	api, blob *http.Client
}

// newClient disables ambient proxy, response compression and redirects.
func newClient(r runtime) *client {
	return &client{r, &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 5 * time.Second}).DialContext, TLSHandshakeTimeout: 5 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 8192, DisableCompression: true}}, newArtifactHTTP()}
}

// request authenticates only a server-constructed api.github.com path.
func (c *client) request(ctx context.Context, path string) (*http.Request, error) {
	r, e := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/"+c.config.repo+path, nil)
	if e != nil {
		return nil, errInvalid
	}
	r.Header.Set("Authorization", "Bearer "+c.config.token)
	r.Header.Set("Accept", "application/vnd.github+json")
	r.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	r.Header.Set("User-Agent", "Devalltect-Pages-Staging")
	return r, nil
}

// read bounds each metadata body and aggregate JSON; all remote diagnostics are masked.
func (c *client) read(ctx context.Context, path string, out any, budget *int) error {
	r, e := c.request(ctx, path)
	if e != nil {
		return errInvalid
	}
	res, e := c.api.Do(r)
	if e != nil {
		return errInvalid
	}
	defer res.Body.Close()
	kind, _, e := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if e != nil || kind != "application/json" || res.StatusCode != 200 || res.ContentLength > 262144 {
		return errInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(res.Body, 262145))
	*budget += len(raw)
	if e != nil || len(raw) < 1 || len(raw) > 262144 || *budget > 1048576 || jsonv2.Unmarshal(raw, out) != nil {
		return errInvalid
	}
	return nil
}

// run retains the finite fields needed for exact first-run provenance and stable rechecks.
type run struct {
	ID         int64     `json:"id"`
	Attempt    int       `json:"run_attempt"`
	Workflow   int64     `json:"workflow_id"`
	SHA        string    `json:"head_sha"`
	Branch     string    `json:"head_branch"`
	Event      string    `json:"event"`
	Title      string    `json:"display_title"`
	Path       string    `json:"path"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	Created    time.Time `json:"created_at"`
	Updated    time.Time `json:"updated_at"`
	Repository struct {
		ID   int64  `json:"id"`
		Name string `json:"full_name"`
	} `json:"repository"`
	HeadRepo struct {
		ID   int64  `json:"id"`
		Name string `json:"full_name"`
	} `json:"head_repository"`
	Pulls *[]json.RawMessage `json:"pull_requests"`
}

// artifact identifies one exact unexpired candidate ZIP under the fixed repository/run/code.
type artifact struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Size    int64     `json:"size_in_bytes"`
	Digest  string    `json:"digest"`
	Expired *bool     `json:"expired"`
	URL     string    `json:"url"`
	Archive string    `json:"archive_download_url"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
	Expires time.Time `json:"expires_at"`
	Run     struct {
		ID       int64  `json:"id"`
		Repo     int64  `json:"repository_id"`
		HeadRepo int64  `json:"head_repository_id"`
		Branch   string `json:"head_branch"`
		SHA      string `json:"head_sha"`
	} `json:"workflow_run"`
}

// inspect refuses reruns/forks/changed code and never follows provider metadata URLs.
func (c *client) inspect(ctx context.Context, r StageRequest, budget *int) (run, artifact, error) {
	var x run
	var a artifact
	root := "https://api.github.com/repos/" + c.config.repo + "/actions/artifacts/" + r.ArtifactID
	now := time.Now()
	if c.read(ctx, "/actions/runs/"+r.RunID, &x, budget) != nil || c.read(ctx, "/actions/artifacts/"+r.ArtifactID, &a, budget) != nil {
		return x, a, errInvalid
	}
	runID, _ := strconv.ParseInt(r.RunID, 10, 64)
	artID, _ := strconv.ParseInt(r.ArtifactID, 10, 64)
	if x.ID != runID || x.Attempt != 1 || x.Workflow != c.config.workflowID || x.SHA != r.HeadSHA || x.SHA != c.config.head || x.Branch != "main" || x.Event != "workflow_dispatch" || x.Title != "CMS Astro candidate "+r.CandidateID || x.Path != ".github/workflows/cms-build-preview.yml" || x.Status != "completed" || x.Conclusion != "success" || x.Repository.ID != c.config.repoID || x.HeadRepo.ID != c.config.repoID || !strings.EqualFold(x.Repository.Name, c.config.repo) || !strings.EqualFold(x.HeadRepo.Name, c.config.repo) || x.Pulls == nil || len(*x.Pulls) != 0 || x.Created.IsZero() || x.Updated.Before(x.Created) || x.Updated.After(now.Add(2*time.Minute)) {
		return x, a, errInvalid
	}
	if a.ID != artID || a.Name != "astro-cms-candidate-"+r.CandidateID || a.Size < 1 || a.Size > MaxArchive || a.Digest != "sha256:"+r.ZIPDigest || a.Expired == nil || *a.Expired || a.URL != root || a.Archive != root+"/zip" || !a.Expires.After(now) || a.Run.ID != runID || a.Run.Repo != c.config.repoID || a.Run.HeadRepo != c.config.repoID || a.Run.Branch != "main" || a.Run.SHA != r.HeadSHA || a.Created.Before(x.Created.Add(-2*time.Minute)) || a.Created.After(x.Updated.Add(2*time.Minute)) || a.Updated.Before(a.Created) || a.Updated.After(now.Add(2*time.Minute)) {
		return x, a, errInvalid
	}
	return x, a, nil
}

// prepare performs GET-only acquisition and exact reconstruction before stable final provenance checks.
func (c *client) prepare(ctx context.Context, r StageRequest) ([]byte, error) {
	budget := 0
	before, a, e := c.inspect(ctx, r, &budget)
	if e != nil {
		return nil, e
	}
	raw, e := c.download(ctx, "/actions/artifacts/"+r.ArtifactID+"/zip", a.Size)
	if e != nil {
		return nil, e
	}
	out, e := pack(ctx, r, raw)
	if e != nil {
		return nil, e
	}
	after, b, e := c.inspect(ctx, r, &budget)
	x, _ := json.Marshal(before)
	y, _ := json.Marshal(after)
	z, _ := json.Marshal(a)
	w, _ := json.Marshal(b)
	if e != nil || string(x) != string(y) || string(z) != string(w) {
		return nil, errInvalid
	}
	return out, nil
}
