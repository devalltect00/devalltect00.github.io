# Local CMS deployment receiver integrity

Updated 8 November 2026 (Asia/Jakarta).

The standalone Go tool now has an independent deployment-intent parser and local
archive validator in `scripts/pages-tool/deploy_request.go` and `deploy_verify.go`.
It pins a distinct deployment UUID, nested exact staging/build identity, canonical
lossless provider IDs, code/payload/wrapper/TAR hashes and TAR size. The request is
bounded to 12 KiB and rejects unknown/duplicate/type-invalid or inconsistent fields.
Candidate reconstruction is bounded by the existing 8 MiB / 32 MiB / 256-entry
policy. The uploaded wrapper has one regular `artifact.tar`, finite directory and
wrapper bounds, exact complete source-derived TAR bytes, CRC and digest checks.
Links, paths/extra members, ZIP64/multidisk/trailers and altered bytes fail closed.
There is no extraction, output writer or provider access in this validator.

This code is not exposed through CLI, workflow, HTTP or runtime. **No protected
receiver or Pages publication is enabled.** Public intent and matching hashes are
integrity data, not current CMS permission, provider-run provenance or private
approved visitor semantics. The future receiver needs authenticated fresh current
Owner/source/target authority, ordering/replay control and protected reviewer/OIDC/
Pages permission checks before any side effect. A queued stale request must never
be authorized by its unsigned manifest or environment name alone.

`go test`/`go vet` for the standalone module include synthetic archive cases.
Optional joint compiled proof:

```text
pnpm test:cms -- --go D:/project/develop/devalltect-cms/backend
```

That opt-in check writes unsigned synthetic public artifacts and validates actual
initial/update/empty compiled output against the private checker and independent
receiver. It never accesses CMS settings/database, live credentials, or deploys.
Existing staging/packing modes and candidate/Pages workflows remain unchanged.
No new dependencies, migration or application version. See
[security evidence](security/cms-deployment-receiver.md) and
[editable diagram](diagrams/cms-deployment-receiver.drawio).
