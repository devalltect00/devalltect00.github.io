// Signed capabilities never receive GitHub credentials or enter diagnostics.
package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var azureArtifactHost = regexp.MustCompile(`^[a-z0-9]{3,24}\.blob\.core\.windows\.net$`)

// artifactLocation permits only canonical HTTPS artifact service hosts, no userinfo/ports/fragments or nested redirects.
func artifactLocation(raw string) (*url.URL, error) {
	u, e := url.Parse(raw)
	if e != nil || len(raw) > 8192 || len(raw) < 1 || u.Scheme != "https" || u.User != nil || u.Port() != "" || u.Fragment != "" || u.Path == "" || u.Path == "/" || u.RawQuery == "" || u.Host != u.Hostname() || (!azureArtifactHost.MatchString(u.Host) && u.Host != "results-receiver.actions.githubusercontent.com") {
		return nil, errInvalid
	}
	return u, nil
}

// publicArtifactIP prevents private, loopback, link-local, carrier NAT, reserved and translation/tunnel destinations.
func publicArtifactIP(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsValid() || !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return false
	}
	for _, raw := range []string{"0.0.0.0/8", "100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "2001::/32", "2001:db8::/32", "2002::/16"} {
		if netip.MustParsePrefix(raw).Contains(ip) {
			return false
		}
	}
	return ip.Is4() || netip.MustParsePrefix("2000::/3").Contains(ip)
}

// dialArtifact resolves once, rejects the complete answer if any address is private, then dials a checked literal IP.
func dialArtifact(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, e := net.SplitHostPort(address)
	if e != nil || port != "443" {
		return nil, errInvalid
	}
	ips, e := net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	if e != nil || len(ips) < 1 {
		return nil, errInvalid
	}
	for _, ip := range ips {
		if !publicArtifactIP(ip) {
			return nil, errInvalid
		}
	}
	// The HTTP transport still performs TLS verification against the original hostname.
	return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ips[0].String(), "443"))
}

// newArtifactHTTP has no auth jar, ambient proxy, automatic redirects or compressed response decoding.
func newArtifactHTTP() *http.Client {
	return &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }, Transport: &http.Transport{Proxy: nil, DialContext: dialArtifact, TLSHandshakeTimeout: 5 * time.Second, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 8192, DisableCompression: true, MaxIdleConnsPerHost: 1}}
}

// download isolates one authenticated redirect from the credential-free, public-IP-pinned blob client.
func (c *client) download(ctx context.Context, path string, size int64) ([]byte, error) {
	req, e := c.request(ctx, path)
	if e != nil {
		return nil, errInvalid
	}
	res, e := c.api.Do(req)
	if e != nil {
		return nil, errInvalid
	}
	location := res.Header.Get("Location")
	res.Body.Close()
	if res.StatusCode != 302 {
		return nil, errInvalid
	}
	if _, e = artifactLocation(location); e != nil {
		return nil, e
	}
	req, e = http.NewRequestWithContext(ctx, "GET", location, nil)
	if e != nil {
		return nil, errInvalid
	}
	req.Header.Set("Accept", "application/zip, application/octet-stream")
	res, e = c.blob.Do(req)
	if e != nil {
		return nil, errInvalid
	}
	defer res.Body.Close()
	if res.StatusCode != 200 || res.ContentLength > size || res.Header.Get("Content-Encoding") != "" && res.Header.Get("Content-Encoding") != "identity" {
		return nil, errInvalid
	}
	kind := strings.SplitN(res.Header.Get("Content-Type"), ";", 2)[0]
	if kind != "application/zip" && kind != "application/octet-stream" {
		return nil, errInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(res.Body, size+1))
	if e != nil || int64(len(raw)) != size {
		return nil, errInvalid
	}
	return raw, nil
}
