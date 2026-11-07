# Run observation boundary verification

Reviewed 7 November 2026. Homepage baseline 5ef4fb9; CMS baseline 331153d.
Homepage changes are documentation/diagrams only. Its candidate workflow and
Pages workflow remain unchanged; this checkpoint does not run them on GitHub.

Private-side tests exercise exact metadata correlation with a fixed GitHub GET,
finite string IDs/states, native/browser Owner/CSRF, source/Owner revocation,
append-only audit rollback, dedupe, timestamp normalization/terminal fencing,
resource limits and concurrent observation revisions. Flutter tests cover visible
errors, confirmation, pending disable and retained SQL reads without provider work.
Measured final results and explicit limits are recorded below.

Metadata is not dispatch-input attestation or artifact/deployment proof. A copied
run title is not an authentication credential. Original uncertain outcomes remain
historical; no check enables another dispatch or marks publication complete.
There is no archive fetch, cancellation or rerun. The reused private API token's
Actions-write authority remains broader than this read-only operation.

See [operations](../cms-run-observations.md) and CMS
`docs/security/astro-run-observations-verification.md` for source/test mapping.

## Final executed results

Homepage `pnpm validate` passed: **48 diagnostic files, zero errors/warnings/hints**,
**25 unit tests in 6 files**, **12 repository checks**, formatting, lint and authored
production build. Source/workflows/content/dependencies are unchanged. Existing
M5k independent/signed compiled/browser proofs were not rerun or counted as new
evidence for this documentation-only consumer checkpoint.

The CMS full native backend tests passed with **7995/9296 statements (86.00%)**
coverage and real PostgreSQL 18 tests, including six concurrent read preflights,
revision fencing, audit rollback, expiry and same-second completion. The final
test process exited zero; its scratch database stop wait exceeded 25 seconds.
A separate verified-directory check confirmed the server was no longer running,
then removed only that private cluster. The wrapper's cleanup timeout remains
documented; no editing database was touched.

Flutter analysis and **194 tests** passed at **8155/8631 lines (94.48%)** coverage.
Its JavaScript release web build passed. Wasm dry-run success is not a Wasm runtime
test; the existing Cupertino font warning remains. No new Android or Linux race
run is claimed. Final Go vet and reproducible SQL generation passed. Go vulnerability
analysis reports zero affected call paths/imported-package findings and one
required-module advisory without a detected called path.

Both complete JavaScript audits contain zero known vulnerability counts.
Redacted source/history scans passed: curated sources cover **153 homepage /
968 CMS** Git-visible UTF-8 files, excluding ignored private settings/caches.
Synthetic source/env secret and tracked/staged private-artifact canaries passed.
Both untagged message files passed Custy validation. The three-page draw.io source
and PNG previews were synchronized and visually checked.

No live provider read/dispatch, credential change, migration of editing data,
deployment, commit, tag or push occurred. Metadata is correlation only; independent
artifact/frozen-content and deployed-output proof plus durable completion are next.
