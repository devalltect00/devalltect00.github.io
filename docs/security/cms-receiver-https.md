# Receiver HTTPS client evidence

Measured 8 November 2026 (Asia/Jakarta).

Standalone public `go test ./... -count=1 -cover` passed at **88.2%** module
statement coverage; `go vet ./...` passed. `TestReceiverExchange` verifies complete
candidate/wrapper bytes and fixed POST/body/bearer, matching lossless run/hash/
intent/scope/expiry and no-store, repeat-call refusal, invalid code/source-run,
cancellation, malformed token, changed candidate/upload, wrong/expired/future/
unknown/duplicate/null decisions, status/redirect/MIME/cache/encoding/size errors.
Invalid local data makes zero calls; each other path makes at most one.
`TestReceiverConfiguration` checks production transport and rejects unsafe origin,
foreign ref, rerun and noncanonical run configuration. Expired decisions fail the
final point-in-time checker. These are synthetic loopback responses, not a live
protected deployment or joined public-client/private-OIDC authorization proof.

Implementation: `scripts/pages-tool/receiver_client.go` and its tests in the Astro
development copy. Existing local canonical packing and request parsing run before
network I/O. Reviewed configuration—not token/intent fields—selects the fixed CMS
origin and expected run/code. Cookies/proxy/redirects/compression are disabled;
header/body/time limits constrain remote data. Strict response matching prevents
mixing another request/run or accepting a self-consistent foreign payload. No token,
URLs or upstream bodies are exposed in errors/logs. A per-instance atomic spent bit
prevents repeated calls; durable replay control remains in the server ledger.

Limits: no runtime/CLI/workflow entry point, OIDC token acquisition, live TLS/public
issuer proof, receiver-side ordering or Pages permission/write. Instance-local
state does not survive restart, nor guarantee exactly-once cross-provider effects.
Lost replies/429 are not automatically retried. System TLS/DNS trust, reviewed
configuration, synchronized clocks and trusted runner remain assumptions. A stale
SQL decision can still be within its TTL; `validFor` is not a stale-write ordering
protocol. No new race/device/load/vulnerability/full-backend coverage is claimed.
No editing database, migration, credentials, production copy or hosted configuration
was changed. Old joint compiled proofs are historical; no writer/schema/build-output
contract changed, and those optional modes were not newly measured here.

Final source-matched Astro `pnpm validate` passed (format/lint, diagnostics,
unit/repository checks and authored static build). Redacted changed-source and
retained-history secret scans passed in both repositories. Both Custy messages
and whitespace checks passed. Native diagrams/PNGs are mirrored and visually
reviewed. No commit, tag, credential or hosted operation occurred.
