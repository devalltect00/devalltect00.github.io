# Public Pages payload verification (M5o)

Reviewed 8 November 2026 (Asia/Jakarta), after CMS **6c650e8** and homepage
**685a49e**. Commands ran across 7–8 November. This is local preparation,
cross-language integrity and compiled fixture evidence. It is not a hosted Pages
acceptance test, external penetration test, malware review or deployment approval.

The private Go implementation and tests live in the CMS development repository;
this public report documents their cross-repository contract without credentials.

## Implemented controls and meaningful checks

| Boundary                                  | Source                                                            | Evidence                                                                                                                                                                                                         |
| ----------------------------------------- | ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Trusted exact source ZIP and frozen facts | Go `astrobuild/pages.go` plus existing `archive.go`/HTML verifier | Wrong/malformed source pin, altered ZIP, cancelled context, stale request and wrong visitor content reject without payload output; existing hostile ZIP tests remain                                             |
| Complete regular-file preservation        | `PreparePages`, `PagesEvidence`/`PagesFile`                       | Binary bytes 0–255, empty file, CSS, long USTAR prefix, omitted empty directory metadata, all body hashes/sizes and original ZIP identity checked                                                                |
| Deterministic wrapper                     | Canonical Go gzip/USTAR writer                                    | Reordered source ZIPs produce identical payload bytes while retaining distinct source ZIP hashes; fixed mode/UID/GID/time and no link/PAX metadata                                                               |
| Exclusive operator output                 | `administration/astroverify/pages.go`, command flags              | New-file-only output, existing output unchanged, wrong/missing ZIP/request/pin/flags/parent, unrelated verify flags and stdout write failure                                                                     |
| Independent public parser                 | Astro `scripts/content/pages-payload.mjs`                         | Six built-in Node cases: strict/coercible/unknown manifest rejection, runtime asset corruption, link/hardlink/PAX/header substitution, CRC/checksum/padding/ending/multigzip rejection and bounded decompression |
| Actual interoperation                     | Optional `prove-cms-build.mjs` and `TestCompiledAstroArtifact`    | Each real initial/update/empty candidate goes through Go frozen ZIP/content, canonical payload preparation and independent Node complete inventory validation                                                    |

All regular member bytes remain exact; wrapper/header bytes are normalized.
No member is extracted, executed, served or fetched by either payload checker.
Go independently derives expected visitor content; Node checks the transport
container and all files, rather than pretending to independently derive those
semantics from its input manifest. No new package, Go module, runtime config,
database table, API, Flutter UI or workflow permission is introduced.

The source ZIP stays capped at 8 MiB, expanded data at 32 MiB and entries at 256.
USTAR is at most 32 MiB + 256 KiB + 1024 bytes; gzip is at most that + 1 MiB;
the local manifest is at most 1 MiB. Strict headers and single-member gzip reject
hidden alternate interpretation, extra tails, links, aliases and PAX/base-256
extensions. Unrepresentable USTAR paths fail instead of silently changing names.
Compression is bounded CPU work; the 15-second local context is checked between
bounded reads/files and is not a hard real-time interrupt guarantee.

## Authority, local filesystem and security limits

The **unsigned public manifest is not approval or deployment authority**. A
malicious party able to replace both payload and manifest can create a mutually
consistent pair. The CLI must receive the expected source ZIP digest from trusted
retained artifact proof, and independently checks frozen content. A future
coordinator must repeat current locked Owner/source/target/attempt/proof checks
before any external mutation and bind that exact evidence to the deployed bytes.

Operator-selected input/output paths are trusted local administration paths,
not content-supplied paths or an upload API. Leaf symlinks/non-files and oversized
input fail; O_EXCL prevents overwriting an existing file/link. An I/O failure can
leave a partial new local output, which must be discarded before retrying.
These paths are not a sandbox against another local user changing an ancestor
directory concurrently. Keep generated archives/manifests in the existing ignored
CMS `.tmp/` or Astro `.cms-cache/`, outside source/deployment inventories.

Determinism is tested under the pinned Go compressor, not guaranteed across
toolchain versions. Existing candidate/Pages workflows remain unchanged. GitHub
documents the gzip/TAR format and link restriction in its
[official Pages guidance](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages#configuring-the-upload-pages-artifact-action).
Hosted Actions wrapping/upload/Pages acceptance and provider deployment are not
exercised here. Integrity checks are not general CSS/JavaScript/image malware
review; reviewed code/host/provider trust remains a separate concern.

## Executed native backend results

The complete `go test -count=1 -coverpkg=<unchanged handwritten CI scope>
-coverprofile=<external profile> ./...` passed against a new private PostgreSQL 18
fixture. The integration package passed in **312.704 seconds**. Coverage measured
**8564/9950 statements, 86.07%**; generated SQL and process entrypoints remain
outside the coverage scope. Targeted new archive/CLI and existing artifact tests
also passed. An existing Windows symlink-permission skip remains explicit.

The fixture's 30-second fast-stop wait timed out after passing tests. The runner
then stopped **only that verified disposable cluster** with immediate mode,
verified its postmaster PID disappeared and removed only its private validated
directory. This is safe disposal of a throwaway fixture, not a recovery procedure
for the editing database. No fixture remains. No editing schema was changed.
Read-only verification before implementation confirmed latest migration **25**,
its applied row, recovery table and enabled immutable trigger.

`go vet ./...` passed. `govulncheck@v1.8.0 ./...` reports **0 affected call paths,
0 imported-package vulnerabilities and 1 required-module advisory without a
detected called path**. That advisory remains documented; unknown vulnerabilities
and configuration weaknesses are not ruled out. SQL/schema/generated bindings
did not change, so sqlc regeneration is not needed for this checkpoint.

Flutter runtime/API/OpenAPI and dependencies are unchanged; Flutter tests and web
build were not rerun here. The prior M5n 206-test/web evidence remains historical.
No new Android/device/APK or Linux race run is claimed. The unchanged CI race gate
still covers these existing handwritten Go packages. Native Windows does not
exercise every POSIX branch.

## Remaining delivery work

Controlled Actions artifact wrapping/upload, protected exact-byte deployment,
independent deployed-output verification and atomic attempt-fenced Astro job
completion remain separate implementation gates. A passing local manifest does
not mark a release delivered. Existing candidate recovery and its three lifetime
admissions remain unchanged. Other consumer integrations remain separate.
See [public contract](../cms-pages-payload.md) and
[native diagram](../diagrams/cms-pages-payload.drawio). No migration beyond 25,
API/Flutter restart, commit/tag/push, live credential or deployment is required or
performed by this local checkpoint.

## Final public and cross-repository verification

Source-matched Astro `pnpm validate` passed: formatting/lint, **51 diagnostic
files with 0 errors/warnings/hints**, **25 unit tests in 6 files**, **18 repository
checks** (including six new Node cases) and authored static build.

Both optional compiled proof commands passed in the external workspace, for
independent public source and the signed local worker fixture. Together they
performed **6 frozen-package Go checks, 6 exact requests, 6 pinned builds,
6 standalone semantic HTML checks, 6 complete ZIP protocol/content checks,
6 new Go Pages preparations and 6 new independent Node complete payload checks**.
The existing **8 bilingual desktop/mobile journeys and 16 themed Axe checks**
also passed. Initial/update/empty content, removals, optional absence and invalid
pin preservation remain tested. Sources stopped before build/browser reads.
Synthetic provider metadata around actual compiled output is not hosted provenance.

An initial Windows proof launch could not resolve the Go executable. It was
rerun successfully using the existing `HOMEPAGE_GO_BINARY` test setting with the
installed executable. The failed launch is not counted as passing verification.
On Windows, an invoking terminal can set
`$env:HOMEPAGE_GO_BINARY = (Get-Command go).Source` if Node cannot find `go`.

Complete JavaScript audits report zero known vulnerability counts in both repos.
Fully redacted retained-history and curated Git-visible UTF-8 source scans passed;
final source covers **1019 CMS files / 171 homepage files**.
Ignored dotenv, private keys, caches, generated archives and binary fixtures are
excluded. Synthetic source/env secret canaries and tracked/staged private binary
path canaries passed. Final reports were formatted and rescanned with exact
curated/source byte parity checked. These bounded checks cannot certify the
absence of unknown flaws or weaknesses in the host/provider configuration.

Three native draw.io pages and PNG previews are mirrored byte-for-byte and
visually reviewed. Separate untagged Custy messages passed the installed validator.
Git heads/indexes, app/dependency versions/locks, API/OpenAPI/Flutter/schema,
both hosted workflows remain preserved. The unrelated portfolio file bytes are
preserved; its existing staged `next-env.d.ts` stays outside this checkpoint. No temporary
PostgreSQL fixture remains. No commit, tag, push, real provider action, credential
setup/rotation, editing migration or deployment was performed.
