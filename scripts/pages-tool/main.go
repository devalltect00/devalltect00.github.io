// Command pages-tool prepares exact TARs and offers an explicitly gated operator deployment command.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

var errInvalid = errors.New("PAGES_STAGE_UNVERIFIED")

// read refuses links/non-files and allocation overruns in operator-selected local input.
func read(name string, max int64) ([]byte, error) {
	s, e := os.Lstat(name)
	if e != nil || !s.Mode().IsRegular() || s.Size() < 1 || s.Size() > max {
		return nil, errInvalid
	}
	f, e := os.Open(name)
	if e != nil {
		return nil, errInvalid
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, max+1))
	if e != nil || int64(len(raw)) > max {
		return nil, errInvalid
	}
	return raw, nil
}

// settings requires opt-in, first attempt, main branch and exact request/code before creating a network client.
func settings(get func(string) string, r StageRequest) (runtime, error) {
	repo := get("GITHUB_REPOSITORY")
	rid, e := strconv.ParseInt(get("GITHUB_REPOSITORY_ID"), 10, 64)
	wid, f := strconv.ParseInt(get("CMS_CANDIDATE_WORKFLOW_ID"), 10, 64)
	token := get("GITHUB_TOKEN")
	if get("CMS_STAGE_ENABLED") != "true" || get("GITHUB_REF") != "refs/heads/main" || get("GITHUB_RUN_ATTEMPT") != "1" || get("CMS_STAGE_ID") != r.ID || get("GITHUB_SHA") != r.HeadSHA || !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,99}/[A-Za-z0-9][A-Za-z0-9._-]{0,99}$`).MatchString(repo) || e != nil || f != nil || rid < 1 || wid < 1 || !regexp.MustCompile(`^[A-Za-z0-9_]{32,512}$`).MatchString(token) {
		return runtime{}, errInvalid
	}
	return runtime{repo, r.HeadSHA, token, rid, wid}, nil
}

// execute separates packing/staging from the explicitly opted-in deployment command.
func execute(args []string, get func(string) string, output io.Writer) error {
	if len(args) > 0 && args[0] == "deploy-local" {
		return runDeploymentCommand(args[1:], get, output)
	}
	if len(args) < 1 {
		return errInvalid
	}
	flags := flag.NewFlagSet("pages-tool", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	request := flags.String("request", "", "local staging request")
	zipFile := flags.String("zip", "", "local candidate ZIP")
	out := flags.String("out", "", "new artifact.tar")
	if flags.Parse(args[1:]) != nil || flags.NArg() != 0 || *out == "" {
		return errInvalid
	}
	var raw []byte
	var e error
	if args[0] == "pack" {
		raw, e = read(*request, 8192)
	} else if args[0] == "stage" && *request == "" && *zipFile == "" {
		raw = []byte(get("CMS_STAGE_REQUEST"))
	} else {
		return errInvalid
	}
	if e != nil {
		return errInvalid
	}
	r, e := ParseStageRequest(raw)
	if e != nil {
		return errInvalid
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var body []byte
	if args[0] == "pack" {
		raw, e = read(*zipFile, MaxArchive)
		if e != nil {
			return e
		}
		body, e = pack(ctx, r, raw)
	} else {
		cfg, x := settings(get, r)
		if x != nil {
			return x
		}
		c := newClient(cfg)
		defer c.api.CloseIdleConnections()
		defer c.blob.CloseIdleConnections()
		body, e = c.prepare(ctx, r)
	}
	if e != nil {
		return errInvalid
	}
	f, e := os.OpenFile(*out, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return errInvalid
	}
	n, e := f.Write(body)
	closed := f.Close()
	if e != nil || closed != nil || n != len(body) {
		return errInvalid
	}
	_, e = io.WriteString(output, "Exact Pages TAR prepared for a separate artifact upload. No deployment or CMS receipt.\n")
	if e != nil {
		return errInvalid
	}
	return nil
}

// main emits only a constant diagnostic, never paths, capabilities or upstream bodies.
func main() {
	if execute(os.Args[1:], os.Getenv, os.Stdout) != nil {
		diagnostic := errInvalid
		if len(os.Args) > 1 && os.Args[1] == "deploy-local" {
			diagnostic = errReceiver
		}
		fmt.Fprintln(os.Stderr, diagnostic)
		os.Exit(1)
	}
}
