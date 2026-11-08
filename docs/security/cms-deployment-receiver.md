# Deployment receiver integrity evidence

Measured 8 November 2026 (Asia/Jakarta).

Standalone public `go test ./... -count=1 -cover` and `go vet ./...` passed at
**87.2%** module statement coverage. Tests use canonical candidate/TAR/wrapper
fixtures and reject wrong scope, duplicate/unknown JSON, precision-losing numeric
IDs, reused UUID, altered payload/candidate/TAR, truncated/extra/traversal/link/
trailing wrapper members, wrong pins/size and cancellation. No extraction or
filesystem output occurs. Optional compiled fixture input is explicit and bounded;
without it the joint-only test is skipped rather than presented as evidence.

The private compiled-output fixture now creates exact public deployment request
and wrapper from independently verified compiled content. The public proof calls
`TestCompiledDeployReceiver` for each initial/update/empty site. That test parses
private-generated JSON and verifies both complete archives independently. This
separates cross-repository compatibility from synthetic unit fixtures.

No protected authority, CLI/runtime route, workflow or Pages/OIDC permission has
been added. The validator does not prove current CMS authorization, provider
provenance, private visitor semantics or deployment/completion. Source/target/Owner
freshness, replay ordering, protected receiver/identity and independent deployed
output proof remain gates. No live credential, hosted operation, editing database,
migration or production-copy edit occurred. New Linux race/device/load and full
backend coverage/vulnerability results are not claimed.

## Final executed checks

Source-matched Astro `pnpm validate` passed: format/lint, Astro diagnostics,
unit/repository tests and authored static build. `pnpm test:cms -- --go <backend>`
passed the genuine source-independent initial/update/empty compiled proof,
including all three new private-to-public deployment intent/archive checks and
existing browser/theme/accessibility checks. Signed-source mode was not rerun.
The first compiled attempt stopped because its child PATH lacked Go; the successful
run explicitly set HOMEPAGE_GO_BINARY. That failed attempt is not evidence.
Private provider package tests and vet passed with the expanded fixture. The
public standalone Go module tests/vet passed at 87.2% package coverage. No new
full backend coverage measurement or provider-authority proof is claimed.

Final redacted changed-source and retained-history scans passed in both repositories.
Both untagged Custy messages and whitespace checks passed. Native diagrams and
PNG previews are mirrored and were visually reviewed.
