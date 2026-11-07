# Exact Astro build handoff: security verification

**7 November 2026 (Asia/Jakarta)**. Baselines: CMS `5f09c5d`, homepage `3687134`.
Local checkpoint, not a production penetration-test certificate or zero-risk claim.
The portfolio development copy and production repositories remain outside changes.
Paths below explicitly distinguish CMS from homepage.

## Implemented controls and proof locations

| Risk                                                 | Implemented control                                                                                                              | Real verification                                                                                                                     |
| ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Build reports A while consuming newer B              | Strict exact intent checked before cache writes and before Astro; repeated frontmatter check                                     | Homepage `cms-request.test.ts`: real child CLI stale-source, retained-byte checks; compiled proof mismatched-intent preserves dist    |
| Sequence rounding / ambiguous JSON                   | Decimal canonical positive int64 strings, duplicate/unknown field rejection and bounded UTF-8                                    | Go request/package tests, homepage malformed/overflow/duplicate tests                                                                 |
| Self-certified public marker                         | Go derives projection independently from frozen approval, computes full expected identity/digest and checks semantic output      | Both real Go request/verify compiled proofs; Go `TestSemanticOutput` and ordering/removal tests                                       |
| HTML/JSON string coincidence mistaken for UI         | HTML5 parsed nodes; literal text and explicit link kinds/URLs; scripts/comments/templates excluded                               | Mutations retain the marker while changing text/tech/status/version/link/order or moving evidence to scripts/comments                 |
| Hidden or duplicate evidence                         | Explicit hidden/inert/aria-hidden rejection, only known numeric animation-delay inline style, exact counts and optional absence  | Hidden main/row, opacity/font-size, duplicate field/attribute and extra article/email tests                                           |
| Parser repairs ambiguous raw markup                  | Bounded raw-token attribute scanner before HTML5 duplicate normalization                                                         | Duplicate/case-duplicate attributes, quoted delimiters, raw-tag/token/depth/resource tests                                            |
| Local path/file abuse                                | Fixed dist members, bounded regular-file reads, known link checks, no content-derived URL/path reads                             | CLI missing/directory/oversize/corrupt input, exclusive request creation; symlink regression is present with Windows limitation below |
| Workflow expression injection / destination choice   | Inputs only in environment, strict schema, operator-fixed root HTTPS origin, no direct expression in shell                       | Repository workflow policy tests plus real creator rejection/exclusive creation tests                                                 |
| Excess provider permissions / accidental publication | Pinned actions, repository contents-read only, fixed dist artifact, one-day retention, ten-minute job; no deploy/callback/secret | Workflow policy tests, unchanged Pages workflow and code review; no hosted run claimed                                                |
| Private data in diagnostics                          | Constant `ASTRO_BUILD_UNVERIFIED` / safe JS errors; public identity/hash/count report only                                       | CLI output/error tests and redacted source/history scans                                                                              |

Build requests and reports are public data. Neither authenticates approval, caller,
hosting or job completion. The configured public source, checked-out code and local
operator/workspace are trusted. A future dispatcher must independently obtain its
expected package/intent from active approved CMS state; it must not trust an
attacker-supplied request just because its checksums match. This code adds no
database mutation, HTTP endpoint, account capability or CMS success path.

HTML checks are for the fixed known templates and three named files. They do not
compute arbitrary CSS, execute arbitrary JavaScript, follow content links, inspect
every artifact member, certify all CDN replicas or continuously monitor a deployed
site. A compromised site can manufacture matching output, so authenticated provider
run/target/deployment provenance remains a separate required control. Filesystem
lstat checks do not eliminate races by a malicious local process replacing files.

## Dependency and verification scope

The new direct Go dependency is **golang.org/x/net v0.62.0**, using its maintained
[HTML5 parser](https://pkg.go.dev/golang.org/x/net/html). It replaces the need to
guess rendered structure with regex/substrings. `go mod tidy` leaves only this
intended module addition/checksums. No JavaScript dependency or application version
changes. Artifact hidden-file behavior was checked against the official
[upload-artifact documentation](https://github.com/actions/upload-artifact): the
upload path is only public dist, including its intentional `.well-known` marker.

Go vulnerability analysis (`govulncheck@v1.8.0 ./...`) found **0 affected code paths**
and **0 imported-package vulnerabilities**, with **1 required-module advisory**
without a detected call path. This preserves the prior known module limitation;
it is not a clean module-inventory claim. No newly detected HTML parser path issue.
Full JavaScript audits remain required; dependency audits are point-in-time results.

Targeted handwritten Go coverage measured **376/417 statements (90.17%)**;
final package breakdown and verification counts are recorded below. Core verifier and CLI are both
included in the existing GitHub/GitLab Go coverage scopes and documented checks.
The CLI process is exercised by real `go run` proofs, not hidden inside a synthetic
coverage claim. No full backend/database or Flutter coverage percentage is claimed.

The Windows host could not create Go symlinks, so `TestFilesystemLinks` explicitly
skipped. Other filesystem cases ran; the homepage tests exercise Windows junction
rejection. Linux CI retains the real Go link test and race detector. Docker responded
but its cached `golang:1.27.3-bookworm` image was absent; no image pull, CMS-service
restart or new Linux race result is claimed. No editing database was accessed.

Early adversarial testing found that the HTML5 tokenizer silently discards duplicate
attributes. The verifier now rejects duplicates in raw start-tag tokens before
parsing; the failing mutation and all assertions remain. Inline opacity/font-size
tampering is also explicitly rejected rather than relying only on display/visibility.

Local ignored evidence belongs under
`E:/portfolio-cms-y35sunj8/.portfolio-cache/astro-handoff-independent`,
`astro-handoff-signed`, the targeted `astro-request.cover`, validation/vulnerability
logs and handoff security scan directory. Only public fixtures and ephemeral Worker
test authority are used. Exact commands are in the operational guide.

## Remaining authority and deployment work

No provider workflow was run or cloud origin configured. The manual candidate
workflow definition is locally checked, not certified by a hosted execution.
Owner-authenticated dispatch, fixed repository/ref/run identity, durable target/
attempt/lease fencing, verified deployed outputs, rollback/recovery, migration and
atomic success/audit remain separate implementation. All Astro jobs stay pending.
No commit/tag/push/deployment or production-copy change is performed here.

## Final measured checkpoint results

- Targeted Go tests passed at **376/417 statements (90.17%)** combined:
  independent `astrobuild` **90.8%**, local `astroverify` **87.7%**. One Windows
  symlink test explicitly skipped as documented above. Go publishing tests,
  CLI compilation and `go vet ./...` passed. No full SQL/Flutter suite is claimed.
- Astro full validation passed with **45 files / 0 errors / 0 warnings / 0 hints**,
  **22 unit tests / 5 files**, **11 repository checks**, lint/formatting and authored
  production build. Tests include real stale-source preparation and strict creator
  child processes, plus workflow authority/input policy gates.
- Both final compiled fixture suites passed: **6 Go frozen package validations**,
  **6 independently generated Go requests**, **6 real pinned builds**, **6 independent
  Go semantic HTML verifications** and **8 EN/ID desktop/mobile browser journeys /
  16 themed Axe checks**. Sources stopped before builds/browser reads; selected
  initial/update/empty states, optional fields, literal text, exact markers and
  stale-request/invalid-pin dist preservation stayed required.
- Both complete JavaScript audits reported **no known vulnerabilities**. Redacted
  Git-visible UTF-8 source and retained-history scans passed for both changed
  repositories. Existing source/env secret and private-artifact binary/staged
  canaries passed. No private settings were collected into source evidence.
- Two separate untagged Custy messages passed the installed validator. Three-page
  native diagrams/previews are synchronized and inspected. Git heads/indexes,
  application versions, JavaScript lockfiles, Pages deployment workflow, authored
  records and the unrelated portfolio edit remain preserved. No commit/tag, hosted
  workflow run, deployment, credential change or migration was performed.
