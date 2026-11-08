# Fixed staging-run reader verification

Reviewed 8 November 2026 after CMS 035f60a / Astro 4b20dad.

Native `go test -coverprofile=<external E: profile>` passed for astrogithub,
astrobuild and astrocandidate: **92.0%, 91.5%, 95.0%** package statement coverage,
respectively. `go vet` passed for those same packages. These are targeted
package results, not a new whole-backend/SQL/Flutter coverage measurement.
Candidate tests exercise the unchanged protocol through the refactored reader;
staging tests exercise the separate workflow/title, exact prepared work and
zero-I/O guards. Both use actual bounded HTTP parsing against loopback fixtures,
with no hosted credentials. Source is astrogithub/observation.go and
stage_observation.go; domain inputs/errors are in astrocandidate/stage_observation.go.

The fixed repository/run GET rejects swapped IDs, workflow paths, forks/PRs,
reruns, event/ref/code drift, bad correlation titles, malformed canonical URLs,
timestamp/state errors, duplicate JSON, redirects and unsafe HTTP responses.
Staging target, source-run equality, missing preparation and release/code mismatch
fail before any GET. Errors are constant; provider bodies/tokens are not exposed.

Metadata is not input attestation or uploaded artifact/content proof. Work must
be assembled under current SQL authority by a future service, then final checks
must repeat before immutable evidence is recorded. No new endpoint, background
poller, SQL row, retry, recovery or publication operation is introduced. The
previous local uploaded-byte verifier is not yet connected to acquisition here.
No migration/config/dependency/version change; current schema remains 26. No
editing database or production copy was touched. No commit, tag, hosted action,
deployment, cloud provision or credential operation was performed.

The native diagram and preview are mirrored and visually reviewed. Public
runtime/tool/workflows remain unchanged. Full SQL/device/Linux race/load/hosted
acceptance checks are not claimed for this provider-only checkpoint. The existing
whole-project results remain historical. See the operational boundary document
for next steps and the distinction between metadata and current Owner authority.

Public source-matched `pnpm validate` passed (format/lint, clean 52-file diagnostics, 25 unit tests, 23 repository tests and static build). Final Git-visible UTF-8 source scans passed with fully redacted output. Both separate Custy messages validate; indexes and checkpoint heads remain unchanged.
