# Owner staging upload reconciliation evidence

M5s review, 8 November 2026; CMS 95485aa / Astro c9e02cb.

| Boundary                           | Code                                                             | Meaningful tests                                                                                                                                                                                                    |
| ---------------------------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Fixed staging run and source proof | astrogithub/stage_artifact.go, stage_observation.go, artifact.go | TestStageArtifactAcquisition: real source/upload ZIPs and canonical TAR; failed/rerun/code/repo/expiry/name/size/digest and final run/artifact drift rejection; source rerun during later uploaded work is rejected |
| Credential/egress separation       | download.go                                                      | Existing hostile artifact download tests plus dual-archive fixture assert token-free signed requests; only internal staged downloads have a separate finite maximum; candidate cap is unchanged                     |
| Complete bytes                     | astrobuild/staged.go                                             | Existing hostile wrapper/content tests plus actual candidate-derived uploaded TAR in the real protocol, without extraction                                                                                          |
| Current authority and atomic audit | postgres/astro_stage_proof.go, migration 00027                   | TestStageUploadProof, TestStageUploadFences: SQL-only read/replay, immutable row, Owner/source/observation changes, different selected stage, malformed evidence, audit rollback and missing-schema readiness       |
| Owner HTTP                         | httpapi/astro_stage_proof.go                                     | TestStageUploadOwnerHTTP: real sessions, no-session/Editor/Origin/CSRF rejection, strict request body, visible safe 409, exact replay and no provider calls on GET                                                  |
| Strict responsive UI               | astro_stage_proof*.dart, StageProofPanel                         | Strict lossless IDs/scope/hash/size/time; actual narrow dark confirmation/pending controls, visible failure/reference and no false success, SQL-only refresh                                                        |

## Executed checks

The final native Go suite (`go test -count=1 -timeout=30m -coverpkg=<unchanged
handwritten CI scope> -coverprofile=<external E: profile> ./...`) passed against
disposable loopback PostgreSQL 18; its integration package passed in **590.862
seconds**. Coverage is **9082/10589 statements, 85.77%**, with the same generated
SQL/entrypoint exclusions as CI. Additional focused current SQL/HTTP/concurrent
tests passed in **35.971 seconds**. The simultaneous equal-result test forces
provider overlap and proves one proof/one audit with no provider I/O under locks.
All task clusters stopped and their exact validated temporary directories were
removed; no editing database was read/migrated. `go vet ./...` passed; pinned
sqlc v1.31.1 generation reproduces generated bindings. Native Windows execution
does not claim Linux race-detector results; the existing Linux race CI remains.

Flutter 3.47.9 / Dart 3.13.7 analysis is clean; **215 tests** pass at
**8794/9282 executable lines, 94.74%**. The existing coverage exclusion is only
the trivial entrypoint. Final narrow-card/staging tests were rerun after the UI
wording update. Final release Web build passed in **257.1 seconds**. A source
wording edit during an earlier build left a stale release kernel; only the
generated kernel/JS/stamp mapped to this task's external output were invalidated.
The rebuilt JS was checked to contain the final retained-proof wording. The
existing Cupertino font-family warning remains; Wasm dry run is not runtime proof.
No new Android build or live authenticated Admin browser journey is claimed.

Source-matched public `pnpm validate` passed: format/lint, **52 diagnostic files /
0 errors / 0 warnings / 0 hints**, **25 unit tests**, **23 repository checks** and
authored static build. Both optional compiled proof modes passed. They cover
six genuine initial/update/empty compiled archives, including **six new complete
two-archive private staging-acquisition checks** with synthetic provider authority.
Existing exact request/ZIP/Pages/private request/public Go repack/Node upload-TAR
checks, **eight bilingual desktop/mobile journeys and 16 themed Axe checks** pass.
Sources stop before build/browser reads; no hosted wrapping acceptance is claimed.

Redocly/OpenAPI, formatting and gofmt checks pass. Both separate untagged Custy
messages pass the installed validator. Fully redacted retained-history and final
Git-visible UTF-8 source scans pass for both repositories; ignored live dotenv,
keys, caches and binary archives are excluded. Synthetic source/env secret
canaries and private-binary index canaries pass. Backend govulncheck v1.8.0 has
**zero affected call paths / zero imported-package findings**, with the documented
one uncalled required-module advisory. Dependencies/locks are unchanged; JS audits
are not newly measured here. These checks do not certify unknown vulnerabilities.

Three native draw.io pages/PNGs are mirrored byte-for-byte and visually reviewed.
Checkpoint heads stay **95485aa / c9e02cb** with unchanged indexes. Other consumer
implementation and production directories are untouched. The portfolio's existing
staged next-env.d.ts bytes remain preserved. Hosted enablement/deployment and
version/credential/Git mutations were not performed.

The initial fixture/script wiring errors and missing Dart member documentation
were fixed before passing results. Failed runs are not successful evidence.

## Actual authority and limits

Only acknowledged accepted current requests can acquire success. The provider
reads the successful first staging run, reacquires verified source bytes, selects
one exact unexpired staging artifact, checks provider ZIP pin/size and the whole
candidate-derived TAR, then repeats uploader/artifact **and source run** checks.
No provider operation occurs in a SQL transaction. Final SQL fences require the
same current Owner, active source/pending job, selected candidate, source proof
revision, both target incarnations, prepared tuple and original accepted run.
Proof and safe identifier-only audit commit together. Cached replay avoids I/O.
Original acknowledgements, budgets, jobs and receipts are never updated.

API/Flutter evidence is historical uploaded-byte proof, not current availability,
deployment approval or publication success. Provider and SQL state cannot commit
atomically; later provider changes/revocation remain possible. Future deployment
must recheck authority and exact bytes. Direct permitted GitHub dispatch bypasses
the private ledger; repository/admin/runner/provider compromise remains a limit.
Finite byte/count/deadline limits and one per-process verification slot reduce
resource abuse but do not certify peak memory or production-load behavior.

The service/provider budget is 18 seconds. Candidate downloads keep the existing
8 MiB cap; staged wrappers have MaxPagesTar + 1 MiB. Signed origins/public-IP
pinning, no proxy/redirect forwarding and token/cookie/referer separation remain.
No archive extraction, hosted run/upload, configuration, deployment, credential
operation or editing-database migration was performed. No commits/tags/pushes.
Production copies and unrelated consumers remain untouched. No dependencies or
app versions change. Current migration requirement is **27**; stop old API,
apply it and restart API/Flutter using the operational guide.

Native Windows results are not Linux race/device/hosted acceptance evidence.
The existing compiler font warning/Wasm dry run are not device/runtime proof.
Unknown/live/failed/skipped/expired work gets no success or permission to retry.
Conservative failure/uncertain recovery, protected deployment, deployed output
proof and atomic destination completion remain next gates.
