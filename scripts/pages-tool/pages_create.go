// Workflow-owned Pages creation isolates the ephemeral provider token from the CMS dispatch credential.
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

// pagesCreator pins one reviewed repository, provider token and exact admitted operation plan.
// It is not registered at startup and cannot establish SQL admission itself.
type pagesCreator struct {
	owner, repo, token, origin string
	plan                       receiverPlan
	http                       *http.Client
	spent                      atomic.Bool
}

// newPagesCreator creates one fixed HTTPS sender; token is separate from the workflow OIDC bearer.
func newPagesCreator(owner, repo, token string, plan receiverPlan) (*pagesCreator, error) {
	return pagesCreatorAt(owner, repo, token, plan, "https://api.github.com")
}

// pagesCreatorAt allows literal loopback HTTP only for package-local verification fixtures.
func pagesCreatorAt(owner, repo, token string, plan receiverPlan, origin string) (*pagesCreator, error) {
	u, e := url.Parse(origin)
	run, f := positiveProviderID(plan.RunID)
	_, g := positiveProviderID(plan.ArtifactID)
	sum := sha256.Sum256([]byte("devalltect.astro.pages-operation.v1\n" + strconv.FormatInt(run, 10) + "\n" + plan.RequestHash))
	if !providerName.MatchString(owner) || !providerName.MatchString(repo) || !providerToken.MatchString(token) || !uuid.MatchString(plan.AttemptID) || f != nil || g != nil || !digest.MatchString(plan.RequestHash) || plan.BuildVersion != hex.EncodeToString(sum[:20]) || e != nil || u.User != nil || u.Path != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || origin != u.Scheme+"://"+u.Host || origin != "https://api.github.com" && (u.Scheme != "http" || (u.Hostname() != "127.0.0.1" && u.Hostname() != "::1")) {
		return nil, errReceiver
	}
	t := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 8192, DisableCompression: true, MaxIdleConnsPerHost: 1}
	return &pagesCreator{owner: owner, repo: repo, token: token, origin: origin, plan: plan, http: &http.Client{Transport: t, Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

// Close releases sockets without resetting the spent operation.
func (c *pagesCreator) Close() { c.http.CloseIdleConnections() }

// create requires a live short admission deadline and invokes one exact nonretrying POST.
// Every error is uncertain; a lost reply or new client must never bypass durable admission.
func (c *pagesCreator) create(ctx context.Context, plan receiverPlan, oidc string) ([]byte, error) {
	fail := func() ([]byte, error) { return nil, errReceiver }
	if c == nil || !c.spent.CompareAndSwap(false, true) {
		return fail()
	}
	deadline, ok := ctx.Deadline()
	now := time.Now()
	if plan != c.plan || !ok || !deadline.After(now) || deadline.After(now.Add(10*time.Second)) || ctx.Err() != nil || len(oidc) > 16384 || !regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`).MatchString(oidc) {
		return fail()
	}
	artifact, _ := strconv.ParseInt(plan.ArtifactID, 10, 64)
	body, _ := json.Marshal(struct {
		ArtifactID  int64  `json:"artifact_id"`
		Environment string `json:"environment"`
		Version     string `json:"pages_build_version"`
		OIDC        string `json:"oidc_token"`
	}{artifact, "github-pages", plan.BuildVersion, oidc})
	req, e := http.NewRequestWithContext(ctx, "POST", c.origin+"/repos/"+c.owner+"/"+c.repo+"/pages/deployments", bytes.NewReader(body))
	if e != nil {
		return fail()
	}
	req.GetBody = nil // A consumed mutation body must never be automatically replayed.
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	req.Header.Set("User-Agent", "Devalltect-Homepage-Pages-Create")
	response, e := c.http.Do(req)
	if e != nil {
		return fail()
	}
	defer response.Body.Close()
	typ, _, e := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if e != nil || typ != "application/json" || response.StatusCode != 200 || response.ContentLength > 8192 || response.Header.Get("Content-Encoding") != "" {
		return fail()
	}
	raw, e := io.ReadAll(io.LimitReader(response.Body, 8193))
	if e != nil || len(raw) == 0 || len(raw) > 8192 {
		return fail()
	}
	var ack struct {
		ID        string `json:"id"`
		StatusURL string `json:"status_url"`
	}
	fixed := "https://api.github.com/repos/" + c.owner + "/" + c.repo + "/pages/deployments/" + plan.BuildVersion
	if jsonv2.Unmarshal(raw, &ack) != nil || ack.ID != plan.BuildVersion || (ack.StatusURL != fixed && ack.StatusURL != fixed+"/status") || ctx.Err() != nil {
		return fail()
	}
	return raw, nil
}

// Reviewed repository segments and ephemeral provider credentials never come from deployment JSON.
var providerName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`)
var providerToken = regexp.MustCompile(`^[A-Za-z0-9_]{32,512}$`)

// positiveProviderID rejects noncanonical decimal identifiers before numeric JSON conversion.
func positiveProviderID(s string) (int64, error) {
	n, e := strconv.ParseInt(s, 10, 64)
	if e != nil || n < 1 || strconv.FormatInt(n, 10) != s {
		return 0, errReceiver
	}
	return n, nil
}
