# Bounded candidate recovery evidence (M5n)

Reviewed 7 November 2026 (Asia/Jakarta), following homepage **fa81d01** and private
CMS **fbe9010**. The homepage changes documentation/diagrams only. Candidate and
Pages workflows, authored data, public source, dependencies and proof hook remain
unchanged. The private CMS's full control-to-code/test mapping lives in its
`docs/security/astro-candidate-recovery-verification.md`; this public report
describes the cross-repository contract without publishing private settings.

## Local checks and authority

Private tests exercise real authenticated Owner/browser/native HTTP operations,
CSRF/Origin rejection, strict UUID-only requests and safe reference-bearing
errors. Synthetic fixed GitHub transport tests check exact before/after first-run
correlation, terminal failure, two successful empty lists or two identical valid
expired-artifact lists. Active/unknown/rerun/usable state, malformed/null metadata,
404/timeouts, wrong identity, duplicate artifacts and changing lists fail closed.
No blob download or provider mutation occurs in this method.

Real PostgreSQL tests verify current Owner/source/target/revision fences, audit
rollback, immutable recovery/history, matching concurrent writes, migration 25
readiness, the unchanged three-admission lifetime budget and new UUID isolation.
Retired candidates cannot acquire new observation/artifact work. Current status
reads cannot show old evidence as evidence for a new candidate. Six Flutter
tests verify strict safe wire types, mobile/dark pending controls, explicit
confirmation, reference-bearing errors and separate new-build confirmation.

Full native backend tests passed (**8494/9870 statements, 86.06%**); the
integration package took **394.153 seconds**. Its fixture teardown exceeded the
60-second stop wait; follow-up verified the cluster stopped and removed only its
private validated directory. The targeted fixture cleaned normally. No editing
database was migrated by these tests. Vet and pinned sqlc generation passed;
**37** generated files reproduced byte-for-byte. `govulncheck@v1.8.0` reports
**0 affected call paths / 0 imported-package vulnerabilities / 1 required-module
advisory without a detected called path**. Flutter analysis and **206 tests**
passed at **8462/8941 executable lines, 94.64%** coverage. Earlier scroll-fixture
failures were fixed and are not counted as successful checks.

## Honest limits and remaining gates

Recovery grants only a separate fresh candidate admission inside the original
budget. Original acknowledgements/proofs remain retained; it never resets,
dispatches, cancels, reruns, deletes, deploys or completes publication. Its
existing repository Actions-write/Contents-read credential is broader than this
GET-only operation. Local fixed paths do not narrow that provider grant. The
homepage receives no such credential or new capability.

GitHub/TLS, reviewed code and administrators remain trusted. Two immediate reads
are bounded evidence rather than guaranteed lasting absence; SQL cannot lock
provider state or prevent later external reruns. Missing output does not prove
absence of arbitrary side effects, so preserve the read-only dist-only workflow.
No external penetration test, hosted action, device/APK, new Linux race or
security certification is claimed. M5m compiled ZIP/browser/Axe evidence remains
historical; those unchanged proof paths are not newly claimed here.

No production copy, unrelated portfolio edit, dependency/version, real token,
provider mutation, commit, tag, push, deployment or editing migration changed.
Private operators must stop old API, apply **00025** and restart API/Flutter.
See [contract](../cms-candidate-recovery.md) and
[editable diagram](../diagrams/cms-candidate-recovery.drawio). Controlled exact
verified-byte deployment, independent deployed-output proof and an atomic Astro
receipt remain future work.

## Final cross-stack and repository checks

The JavaScript Flutter release web build passed in **187.0 seconds**, using
`flutter build web --release --base-href /admin/ --no-web-resources-cdn --no-pub`
with external output. The compiler's Wasm dry run succeeded; a runtime Wasm test
is not claimed. The existing Cupertino font-family warning remains, with no new
icon package or handwritten CupertinoIcons use introduced.

CMS whole-repository formatting and Redocly OpenAPI validation passed. A source-
matched external Astro copy passed **pnpm validate**: format/lint, **48 diagnostic
files with 0 errors/warnings/hints**, **25 unit tests in 6 files**, **12 repository
checks** and authored static build. Only documentation and diagrams changed in
the homepage; compiled private proof and hosted workflows were not rerun.

Both complete JavaScript dependency audits reported zero known vulnerability
counts. Fully redacted Git-visible UTF-8 source and retained-history scans passed
for both repositories. Final curated source covers **1011 CMS files /
163 homepage files**; ignored dotenv, private keys, caches,
binaries and fixture DBs are excluded. Synthetic source/env secret canaries and
tracked/staged private-artifact canaries/path policy passed. Final reports were
formatted and source-rescanned; curated bytes were checked against current
source. These are bounded checks, not a guarantee against unknown flaws.

All three native draw.io pages and PNG previews are mirrored byte-for-byte and
visually reviewed. Both separate untagged Custy messages passed the installed
validator. Git heads/indexes, dependency/app versions, unrelated portfolio edit
and production copies remain preserved. No temporary PostgreSQL cluster remains.
No commit, tag, push, real provider action, credential or editing migration was
performed. Apply **00025** with the old API stopped and restart API/Flutter.
