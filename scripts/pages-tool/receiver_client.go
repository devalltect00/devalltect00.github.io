// Receiver HTTPS exchange validates a consumed decision; it never performs Pages writes or runtime/OIDC setup.
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"errors"
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

var errReceiver = errors.New("PAGES_RECEIVER_UNVERIFIED")

// receiverConfig is reviewed server/workflow configuration, never content from deployment JSON or token claims.
type receiverConfig struct{ Origin, RunID, HeadSHA, Ref, Attempt string }

// receiverDecision is meaningful only as the response to this client's fixed authenticated HTTPS request.
type receiverDecision struct {
	AttemptID   string    `json:"attempt_id"`
	RunID       string    `json:"run_id"`
	RequestHash string    `json:"request_sha256"`
	ExpiresAt   time.Time `json:"expires_at"`
	Scope       string    `json:"scope"`
}

// receiverClient retains one fixed origin and permits at most one POST for its workflow instance.
type receiverClient struct {
	config receiverConfig
	http   *http.Client
	spent  atomic.Bool
}

// newReceiverClient accepts reviewed HTTPS configuration only and performs no network work.
func newReceiverClient(c receiverConfig) (*receiverClient, error) {
	return receiverClientWithOrigin(c, false)
}

// receiverClientWithOrigin permits literal-loopback HTTP only for package-local fixtures; no CLI exposes this switch.
func receiverClientWithOrigin(c receiverConfig, loopback bool) (*receiverClient, error) {
	u, e := url.Parse(c.Origin)
	run, f := strconv.ParseInt(c.RunID, 10, 64)
	if e != nil || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || c.Origin != u.Scheme+"://"+u.Host || f != nil || run < 1 || strconv.FormatInt(run, 10) != c.RunID || !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(c.HeadSHA) || c.Ref != "refs/heads/main" || c.Attempt != "1" {
		return nil, errReceiver
	}
	if !(u.Scheme == "https" && (u.Port() == "" || u.Port() == "443") && net.ParseIP(u.Hostname()) == nil && regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.-]*\.[A-Za-z0-9-]+$`).MatchString(u.Hostname())) && !(loopback && u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "::1")) {
		return nil, errReceiver
	}
	transport := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 15 * time.Second, MaxResponseHeaderBytes: 8192, DisableCompression: true, MaxIdleConnsPerHost: 1}
	return &receiverClient{config: c, http: &http.Client{Transport: transport, Timeout: 20 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

// close releases idle sockets without resetting spent permission or exposing token material.
func (c *receiverClient) close() { c.http.CloseIdleConnections() }

// validFor must be rechecked immediately before any future protected write; it grants no ordering/atomicity itself.
func (d receiverDecision) validFor(r DeployRequest, run string, now time.Time) bool {
	raw, _ := json.Marshal(r)
	return d.Scope == "astro-receiver-authorization-v1" && d.AttemptID == r.ID && d.RunID == run && d.RequestHash == hash(raw) && d.ExpiresAt.After(now) && !d.ExpiresAt.After(now.Add(30*time.Second))
}

// authorizeVerified reconstructs complete input before one fixed POST, then strictly validates the consumed decision.
// No retry, redirect, workflow dispatch, extraction, file write or Pages call occurs. Lost replies fail closed.
func (c *receiverClient) authorizeVerified(ctx context.Context, r DeployRequest, candidate, uploaded []byte, token string) (receiverDecision, error) {
	fail := func() (receiverDecision, error) { return receiverDecision{}, errReceiver }
	if len(token) == 0 || len(token) > 16384 || !regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`).MatchString(token) || r.Stage.HeadSHA != c.config.HeadSHA || c.config.RunID == r.RunID || c.config.RunID == r.Stage.RunID {
		return fail()
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if verifyDeploymentInput(ctx, r, candidate, uploaded) != nil || ctx.Err() != nil {
		return fail()
	}
	req, e := http.NewRequestWithContext(ctx, "POST", c.config.Origin+"/api/v1/machine/astro/deployments/"+r.ID+"/authorize", bytes.NewReader([]byte("{}")))
	if e != nil {
		return fail()
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if !c.spent.CompareAndSwap(false, true) {
		return fail()
	}
	response, e := c.http.Do(req)
	if e != nil {
		return fail()
	}
	defer response.Body.Close()
	typ, _, e := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if e != nil || typ != "application/json" || response.StatusCode != 200 || response.ContentLength > 4096 || response.Header.Get("Content-Encoding") != "" || response.Header.Get("Cache-Control") != "no-store" {
		return fail()
	}
	raw, e := io.ReadAll(io.LimitReader(response.Body, 4097))
	if e != nil || len(raw) == 0 || len(raw) > 4096 {
		return fail()
	}
	var envelope struct {
		Data *receiverDecision `json:"data"`
	}
	if jsonv2.Unmarshal(raw, &envelope, jsonv2.RejectUnknownMembers(true)) != nil || envelope.Data == nil || !envelope.Data.validFor(r, c.config.RunID, time.Now()) || ctx.Err() != nil {
		return fail()
	}
	return *envelope.Data, nil
}
