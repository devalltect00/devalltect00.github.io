// Workflow OIDC acquisition isolates runner credentials from CMS and provider API transports.
package main

import (
	"context"
	"crypto/tls"
	jsonv2 "encoding/json/v2"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

// workflowOIDC pins a reviewed issuance origin, runtime endpoint and explicit audience.
// Two lifetime attempts cover initial authorization and report acquisition; no retry or cache is implicit.
type workflowOIDC struct {
	endpoint, credential string
	http                 *http.Client
	calls                atomic.Int32
}

// newWorkflowOIDC checks nonsecret deployment context before reading runner-only OIDC variables.
// origin is protected reviewed configuration, never derived from the endpoint or public dispatch input.
func newWorkflowOIDC(p deploymentPolicy, origin string, get func(string) string) (*workflowOIDC, error) {
	settings, e := validateDeploymentSettings(p, get)
	if e != nil {
		return nil, errReceiver
	}
	return workflowOIDCAt(origin, get("ACTIONS_ID_TOKEN_REQUEST_URL"), get("ACTIONS_ID_TOKEN_REQUEST_TOKEN"), settings.Audience, false)
}

// workflowOIDCAt accepts loopback only in package fixtures; production pins HTTPS and exact reviewed origin.
func workflowOIDCAt(origin, endpoint, credential, audience string, loopback bool) (*workflowOIDC, error) {
	o, e := url.Parse(origin)
	u, f := url.Parse(endpoint)
	if e != nil || f != nil || o.User != nil || o.Path != "" || o.RawQuery != "" || o.ForceQuery || o.Fragment != "" || origin != o.Scheme+"://"+o.Host || u.User != nil || u.Fragment != "" || u.Scheme+"://"+u.Host != origin || u.Path == "" || len(endpoint) > 4096 || len(credential) < 1 || len(credential) > 16384 || !regexp.MustCompile(`^[A-Za-z0-9._~-]+$`).MatchString(credential) || !regexp.MustCompile(`^https://github\.com/[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`).MatchString(audience) {
		return nil, errReceiver
	}
	production := o.Scheme == "https" && (o.Port() == "" || o.Port() == "443") && strings.HasSuffix(o.Hostname(), ".actions.githubusercontent.com") && regexp.MustCompile(`^[a-z0-9.-]+$`).MatchString(o.Hostname())
	if !production && !(loopback && o.Scheme == "http" && (o.Hostname() == "127.0.0.1" || o.Hostname() == "::1")) {
		return nil, errReceiver
	}
	query, e := url.ParseQuery(u.RawQuery)
	if e != nil {
		return nil, errReceiver
	}
	if _, present := query["audience"]; present {
		return nil, errReceiver
	}
	query.Set("audience", audience)
	u.RawQuery = query.Encode()
	transport := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, TLSHandshakeTimeout: 3 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 8192, DisableCompression: true, DisableKeepAlives: true}
	return &workflowOIDC{endpoint: u.String(), credential: credential, http: &http.Client{Transport: transport, Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

// close releases idle connections without resetting the lifetime issuance budget.
func (c *workflowOIDC) close() { c.http.CloseIdleConnections() }

// acquire makes one bounded GET and returns only a compact token; claims here are not trusted or signature-verified.
// CMS must independently verify the issuer, signature, identity, audience, time and current durable authority.
func (c *workflowOIDC) acquire(ctx context.Context) (string, error) {
	if c == nil || c.calls.Add(1) > 2 || ctx.Err() != nil {
		return "", errReceiver
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	req, e := http.NewRequestWithContext(ctx, "GET", c.endpoint, nil)
	if e != nil {
		return "", errReceiver
	}
	req.Header.Set("Authorization", "Bearer "+c.credential)
	req.Header.Set("Accept", "application/json")
	response, e := c.http.Do(req)
	if e != nil {
		return "", errReceiver
	}
	defer response.Body.Close()
	typ, _, e := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if e != nil || typ != "application/json" || response.StatusCode != 200 || response.ContentLength > 32768 || response.Header.Get("Content-Encoding") != "" {
		return "", errReceiver
	}
	raw, e := io.ReadAll(io.LimitReader(response.Body, 32769))
	if e != nil || len(raw) > 32768 {
		return "", errReceiver
	}
	var token struct {
		Value string `json:"value"`
	}
	if jsonv2.Unmarshal(raw, &token, jsonv2.RejectUnknownMembers(true)) != nil || len(token.Value) > 16384 || !regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`).MatchString(token.Value) || ctx.Err() != nil {
		return "", errReceiver
	}
	return token.Value, nil
}
