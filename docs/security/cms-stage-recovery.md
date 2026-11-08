# Conservative staging recovery verification

M5t review, 8 November 2026; private CMS cdd97c5 / homepage 9734df7.

| Control                           | Source                                                             | Executed tests                                                                                                                                                                                                                                      |
| --------------------------------- | ------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Known first-run terminal evidence | astrogithub/stage_recovery.go, stage_observation.go                | TestStagingRecoveryProtocol repeats candidate recovery's hostile real HTTP cases for the separate staging workflow: stable failure/missing/expiry, usable/live/rerun/changed/null/duplicate/foreign identity rejection; private GET-only assertions |
| Finite evidence/current authority | astrocandidate/stage_recovery.go, postgres/astro_stage_recovery.go | TestStagingRecoveryFences: Owner disable, source removal, source observation revision, bad result, wrong stage, unknown acknowledgement, audit rollback and readiness                                                                               |
| History and budgets               | migration 00028, astro_stage.sql / astro_stage.go                  | TestStagingRetirement: exact SQL-only replay, three successive fresh UUIDs, original accepted outcomes retained, immutable UPDATE/DELETE denial and exhausted lifetime budget                                                                       |
| Late and stale proof              | astro_stage_proof.go / stageRecoveryWork                           | TestStagingRecoveryFences/late-proof and TestStagingExpiredProofRecovery: retired work cannot accept late proof or expose old uploaded proof as current; original uploaded proof remains stored                                                     |
| Owner HTTP                        | astro_stage_recovery.go                                            | TestStageRecoveryOwnerHTTP: genuine session/role/Origin/CSRF, strict body, safe 409, SQL-only GET and cached replay                                                                                                                                 |
| Explicit UI consent               | astro_stage_recovery*.dart and StageRecoveryPanel                  | Strict scope/reason/identity/time/conclusion decoding, projected retired state, narrow dark confirmation/pending/error/reference; full parent flow proves retirement makes no dispatch and later confirmation uses a different UUID                 |

## Executed measurements

The complete native `go test -count=1 -timeout=30m -coverpkg=<unchanged CI
handwritten scope> -coverprofile=<external E: profile> ./...` suite passed using
disposable loopback PostgreSQL **18**. Integration passed in **273.591 seconds**;
coverage is **9257/10807 statements, 85.66%** with unchanged generated/entrypoint
exclusions. Focused real recovery SQL/HTTP tests passed in **35.458 seconds**.
Task clusters were stopped and their exact validated directories removed; no
editing database was read/migrated. `go vet ./...`, gofmt and pinned sqlc v1.31.1
reproduction passed. Native tests do not claim Linux race execution; existing
read-only CI retains its race/coverage gates.

Flutter analysis is clean. All **219 tests** pass with **8957/9450 executable
lines, 94.78%**, excluding only the existing trivial entrypoint. Release Web
build passed in **105.4 seconds**. No Dart source was changed during that build.
The existing Cupertino font-family warning remains; the Wasm dry run is not
Wasm/browser/device runtime proof. New narrow dark and whole-parent flow tests
verify explicit consent, disabled pending actions, visible errors/references,
preserved original run/budget and a different UUID only after separate consent.

Public source-matched `pnpm validate` passed: formatting/lint, **52 diagnostic
files with no errors/warnings/hints**, **25 unit tests**, **23 repository checks**
and authored static build. Public runtime/workflows/tools and all dependency
locks/versions remain unchanged. Optional genuine compiled modes were not rerun
because no archive writer/public handoff contract changed. New recovery protocol
tests reuse the complete hostile HTTP suite with separate staging identities.
Redocly/OpenAPI validation and both separate untagged Custy validators passed.

Fully redacted retained-history and final Git-visible UTF-8 source scans passed
in both repositories. Ignored live dotenv, keys, caches and binary archives are
not copied to those inventories. Synthetic source/env secret canaries and
tracked/staged private-binary canaries passed. Backend govulncheck v1.8.0 reports
**zero affected call paths / zero imported-package vulnerabilities**, with the
documented one uncalled required-module advisory. No new dependency/JS audit is
claimed; locks are unchanged. Tests/scans do not certify unknown vulnerabilities.

The three native draw.io pages and PNG previews are mirrored byte-for-byte and
visually reviewed. Heads stay **cdd97c5 / 9734df7**, indexes are untouched, and
the portfolio's pre-existing staged next-env.d.ts is unchanged. No production
copy or other consumer implementation was edited. No credentials were loaded
for hosted calls, rotated or written into source. No hosted action, deployment,
editing migration, commit, tag or push occurred. Migration **28** still needs to
be applied after stopping the old API, followed by API/Flutter restart.

Only acknowledged known work can retire. Unknown/malformed/lost acknowledgement,
local lease expiry, unverified/404/live/rerun/usable/ambiguous metadata never
permits a fresh request. There is no run discovery/binding, workflow mutation,
cancellation, rerun or artifact deletion. Recovery is not deployment/publication.
The original outcomes/proofs/counters remain immutable; current status projects
retired from a separate ledger. A new request is a distinct explicit action and
uses the unchanged three-admission lifetime limit including all previous history.

Before and after provider GETs, SQL rechecks current Owner, active source/pending
job, selected candidate/source proof/revision, both target incarnations, original
staging Owner, prepared request, accepted run and known artifact identity. A
newly added uploaded proof during recovery invalidates the frozen work. Retired
work cannot write/replay a current upload proof. No network work holds SQL locks;
retirement and safe identifier-only audit commit together. Duplicate replay
returns original recorded evidence without external reads. Scope/identity/reason
checks reject corrupted retained JSON rather than inventing success.

The service context is 14 seconds; a per-process recovery slot and existing
fixed API response/header/body budgets bound work. Source downloads and existing
archive verification are unchanged and are not used by recovery. Stable metadata
observations are point-in-time provider evidence, not cancellation or an atomic
availability guarantee. SQL/provider state may change after rechecks. Direct
permitted GitHub operators can bypass the private admission ledger; trusted
admin/runner/provider compromise and production load/race/device behavior remain
limits. Unknown identity needs a separate trustworthy reconciliation design.

Apply migration **28** after stopping the previous API; restart API and Flutter.
No editing database was read/migrated and no hosted operation/configuration,
credential, deployment, commit, tag, push, dependency or app version was changed.
Public workflows/tools/runtime are unchanged. No new archive writer/contract was
introduced; prior genuine compiled handoff evidence remains historical and those
optional modes are not newly measured here. Linux race/Android/live hosted
acceptance and production peak-memory/availability tests are not claimed.
The fixture variable-redeclaration/setup mistakes were corrected before passing
checks; failed runs do not count as successful evidence.
