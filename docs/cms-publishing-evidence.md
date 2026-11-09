# Astro publishing acceptance and security evidence

Measured 9 November 2026 (Asia/Jakarta). This is local release-candidate evidence,
not hosted rollout, penetration testing or a claim of zero vulnerabilities.

## Implemented controls and evidence

| Boundary                              | Source                                                                    | Direct checks                                                                                                                                                                                                                                                                                                |
| ------------------------------------- | ------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Current Owner/source/target authority | `postgres/astro_pages_completion.go`, `astro_deployment_retirement.go`    | Real PostgreSQL tests change Owner/source/policy during work and reject stale completion/retirement.                                                                                                                                                                                                         |
| Atomic completion and audit           | `postgres/astro_pages_completion.go`, migration 36                        | Completion, verification, optional reconciled receipt, Astro job and aggregate state roll back together on audit failure; other consumer jobs stay unchanged. Concurrent completion retains one result.                                                                                                      |
| No duplicate Pages write              | Existing immutable admission plus protected sender                        | One permanent admission precedes one fixed create; lost acknowledgement never permits automatic create retry. Retired attempts cannot dispatch or acquire write permission.                                                                                                                                  |
| Lost callback recovery                | `astrocandidate/pages_completion.go`, `astrogithub/pages_verification.go` | Exact terminal failed workflow may settle only with independently successful Pages operation and complete output. A missing callback is explicitly `owner-reconciled`; completed replay performs no provider work.                                                                                           |
| Safe no-write retirement              | `astrocandidate/deployment_retirement.go`, PostgreSQL adapter             | Thirteen real SQL cases include accepted/uncertain/lost-ack/pre-send, live/wrong run, existing admission/report, Owner/source changes, audit rollback, concurrency and lifetime budget.                                                                                                                      |
| Preserved history and finite retries  | Migration 36, current pointer trigger                                     | Claims, retirement and completion remain immutable. Pointer cannot update/delete without retained retirement. Three historical claims maximum; closure does not reset counts.                                                                                                                                |
| Upgrade safety                        | `astro_completion_migration_test.go`                                      | Downgrade to 35 before settlement, upgrade/backfill existing claim, unchanged original dispatch, migration replay, and refusal to downgrade retained retirement all pass in random test schemas.                                                                                                             |
| Browser permissions and input         | `httpapi/astro_deployment.go`, `astro_completion_http_test.go`            | Actual stored-session middleware checks unauthenticated access, CSRF, foreign Origin, unknown proof fields, Owner-only access and completion replay. Sessions are synthetic fixture inserts; this test does not exercise password/TOTP enrollment.                                                           |
| Archive/credential isolation          | Public `deployment_acquisition.go` and existing private verifier          | Exact repository/workflow/code/run/artifact pins before and after downloads; complete candidate/TAR/wrapper verification, no extraction, no private API headers sent to signed-download/public-file clients. Corruption/provenance/cancellation unit fixtures pass.                                          |
| Workflow authority                    | `cms-pages-deploy.yml`, Node guard tests                                  | Off-by-default, main first attempt, protected environment, reviewed policy byte digest, correlation before checkout, exact source checkout, minimal job permissions and no input-to-shell interpolation. Eleven actual Python guard substitutions reject invalid contexts. No hosted workflow was triggered. |
| UI/error handling                     | Flutter deployment client/panel tests                                     | Central API paths, finite strict status, explicit confirmation, disabled pending action, safe error guidance, completion, closure/budget display and narrow dark layout are tested.                                                                                                                          |

## Executed checks

- Full backend `go test ./...` discovery and `go vet ./...` passed with no editing-database connection. Database-dependent cases skip in that discovery; the separately measured real PostgreSQL
  completion/recovery/claim/retention/verification/admission/ordering suites passed
  in **274.532s**. Disposable evidence: `astro-completion-sql-sq6ljn3c/tests.log` in
  the operator's local E-drive build cache. PostgreSQL was stopped afterward.
- The two additional upgrade/downgrade tests passed in **7.930s**; disposable
  evidence: `astro-completion-sql-mfgnc1ho/tests.log`. Editing databases were not
  migrated or reset by these tests.
- Flutter analysis passed and the complete **227-test** suite passed. These are
  unit/widget/transport tests, not a fresh Android device acceptance session.
- Public Go `go test ./... -count=1 -cover` passed at **85.9% module statement
  coverage** and `go vet ./...` passed. This is the standalone Pages tool's
  coverage, not the backend aggregate coverage.
- Astro `pnpm validate` passed formatting, lint, 54-file diagnostics (zero errors,
  warnings or hints), 25 unit tests, 25 repository tests and authored static build.
- `pnpm test:cms -- --go D:/project/develop/devalltect-cms/backend` passed actual
  initial/update/empty builds, EN/ID selected content/removals, independent Go
  package/HTML/archive/deployment-receiver checks, browser themes/accessibility and
  invalid-pin preservation. It creates no CMS destination completion receipt.
- OpenAPI lint passed. SQL was regenerated from source with sqlc 1.31.1.
- Redacted Gitleaks retained-history and changed-source scans passed for both
  development repositories; whitespace checks passed. No secrets were rotated.
- Five current diagram pages were exported and visually checked: the two-page
  normal/recovery overview and the three focused ordering/verification/receiver
  diagrams. Source diagrams and PNG previews are synchronized.

## Limits and rollout gates

No live token, OIDC issuance, GitHub Pages deployment, protected environment or
provider settings were exercised. Synthetic adapters prove local behavior, not
provider availability, quotas or production proxy configuration. Hosted rollout
still needs reviewed credentials/policy and a separate explicit authorization.

Historical observation capacity reserves one final completion proof (16 ordinary
observations plus one final record). Expired artifacts, unavailable public bytes,
live/ambiguous work, key rotation, a spent admission without successful Pages, or
exhausted attempt budgets remain visible operator stops. Do not delete history,
reset counters or invent target success to recover them.

Migration **00036 must be applied explicitly before the new API starts**; it has
not been applied to the editing database here. Optional publishing configuration
remains empty by default. The local Compose production-like preview currently
does not mount optional publisher credentials/configuration; final production
wiring is part of milestone 6. Drain old authored Pages jobs before enabling the
CMS writer, and set proxy completion deadlines above the 65s service limit.

No new Linux race-detector, fresh full-database aggregate coverage, dependency
vulnerability database scan, Android emulator smoke or backup/restore acceptance
is claimed in this report. Final whole-system checks belong to milestones 4–5.
No commit, tag, push, cloud provisioning, production-copy edit or hosted action
was performed. The goal remains active for the other five outcomes.
