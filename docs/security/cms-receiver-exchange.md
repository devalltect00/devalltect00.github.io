# Receiver exchange client security evidence

Measured 8 November 2026 (Asia/Jakarta).

Implementation: Astro `scripts/pages-tool/receiver_exchange.go`; fixtures in
`receiver_exchange_test.go`. Each admission/report fixture checks fixed machine
method/path/bearer and body. Cases cover exact success, wrong scope/attempt/run/
artifact/hash/version, expired/future/null/missing fields, unknown/duplicate JSON,
status/redirect/MIME/cache/encoding/size/lost-connection failures, changed candidate
or uploaded archive, malformed bearer, code/source-run mismatch and cancellation.
Report also tests invalid/duplicate/oversize acknowledgements and exact replay.
Admission is refused on every second call, including after a failed first call.

Fixtures are synthetic local authority responses, not issuer authentication. The
CMS independently tests RSA, current SQL authority, immutable one-use admission
and operation receipts. No joined hosted deployment, penetration test, global
security guarantee, fleet rate limit or completion/hold discharge is claimed.
The prior authorize helper retains its existing behavior; this report concerns
the new admission/report paths. No dependency, public writer or schema change.

## Verification results

- Standalone `go test ./... -count=1 -cover`: passed, 89.1% module statement
  coverage (not full CMS coverage). Final run 4.622 seconds. The 60 transport
  subcases and eight-caller concurrency test pass; `go vet ./...` passes.
- Astro `pnpm validate`: passed formatting, ESLint, 53-file Astro check with zero
  diagnostics, 25 unit tests, 23 repository tests and static build.
- Both native diagrams export successfully; rendered PNG reviewed for readable
  labels and explicit disabled-runtime/completion gates.
- Public writers and archive format are unchanged; no joint compiled-output proof
  rerun is required for this transport-only change. No database integration code
  changed; the preceding private HTTP/RSA/PostgreSQL checkpoint remains separate.
- No race-detector result is claimed: this Windows environment has no `gcc` in
  PATH. The concurrent admission test runs in the normal Go test suite.
- Gitleaks retained-history scans: CMS 98 commits / 8.18 MB; Astro 22 commits /
  1.07 MB; no leaks found. Changed visible-source scan also passed (250.97 KB)
  before final diagram-label synchronization; final source scan and diff checks also passed.
- Both checkpoint messages passed the installed Custy validator (`docs` / `feat`).
  No tag, commit, staging, push or deployment was performed.
