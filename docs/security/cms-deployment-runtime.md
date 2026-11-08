# Internal runtime verification evidence

9 November 2026 (Asia/Jakarta).

Implementation/tests: Astro `deployment_runtime.go` / `deployment_runtime_test.go`.
Twenty-one joined loopback scenarios exercise initial issuance, authorization,
admission, create, second issuance and report failures, invalid archives and
cancellation. Every second runtime call is refused. Five preflight scenarios
reject invalid policy/request/code/source-run/disabled operation before factory
or credential reads. Production constructor tests check success and policy,
issuer, receiver and provider-token failure/cleanup without network requests.
Report-only recovery never increments create count. No live token, hosted proof,
joined RSA/PostgreSQL runtime or race-detector result is claimed. Public archive
and writer contracts and database implementation are unchanged.

Standalone `go test ./... -count=1 -cover` and `go vet ./...` passed.
Final module coverage is 91.1% (11.227-second test run), not full CMS coverage.
Full `pnpm validate` passed formatting, ESLint, 53-file Astro check (no diagnostics),
25 unit tests, 23 repository tests and static build. Native diagrams exported and
rendered labels reviewed. Both checkpoint messages passed Custy validation.
Retained-history and changed-source Gitleaks scans and whitespace checks passed.
No staging, commits, tags, credentials or hosted deployments were performed.
No database code changed and no migration was required (schema 35).
