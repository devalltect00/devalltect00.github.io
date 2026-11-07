# Public exact Pages staging verification (M5p)

Reviewed 8 November 2026 (Asia/Jakarta), after CMS **67879ad** and homepage
**9bd1121**. This is local implementation, synthetic network and real compiled
output evidence. It is not hosted upload/Pages acceptance, an external penetration
test, a malware review or deployment approval.

## Controls, code and executed tests

| Boundary                         | Code                                                 | Meaningful checks                                                                                                                                                                                                            |
| -------------------------------- | ---------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Exact retained-proof request     | CMS `astroverify/stage.go` and `astrobuild/stage.go` | `TestStageRequestCommand`: frozen package/actual ZIP, proof scope/content/hash mismatch, canonical IDs, missing inputs, exclusive output and status writer error                                                             |
| Strict public wire/context       | Public `request.go` and `main.go`                    | `TestRequestAndSettings`: unknown/duplicate keys, numeric IDs, leading zeros, overflow, reused candidate UUID and every opt-in/ref/attempt/repository/code/token context field                                               |
| Complete portable reconstruction | Public `pack.go`                                     | `TestPackAndCLI`, `TestUnsafeArchives`, `TestPortablePolicy`: complete binary/empty files, directory metadata, source/gzip pins, cancellation, corruption, hostile/full inventories, metadata/path/device-name policy        |
| Fixed provenance                 | Public `github.go`                                   | `TestProviderProtocol`: successful exact first run, rerun/fork/code/workflow/title/PR/status rejection, expiry/size/digest/URL checks and changed run/artifact after packing                                                 |
| Credential separation and egress | Public `download.go`                                 | Synthetic transports assert API credentials stay at fixed GitHub paths and blob headers have no auth/cookie/referer; signed-host/IP/port policy and no-proxy/no-redirect transport tests                                     |
| Local output and failures        | Public `main.go`                                     | `TestLocalFailures`: bad request/ZIP/payload/output parent, failed status output and disabled valid staging input without egress                                                                                             |
| Workflow authority               | `cms-pages-stage.yml` and Node policy tests          | Manual-only/off-by-default/main/first-attempt/environment declarations, pinned actions, no Pages/OIDC-write permission, one TAR path/retention; execute the actual before-checkout Python guard on valid and hostile context |
| Exact real handoff               | Optional compiled proof in both repos                | Private CLI reconstructs actual candidate/retained proof, public Go independently repacks, Node compares the exact upload TAR with verified private gzip/manifest after every initial/update/empty build                     |

The COM0/LPT0 portability range now matches the strict Go policy in Node too.
A dedicated forged-but-consistent Node fixture tests this validation rather than
relying only on an unrelated checksum failure. New Node TAR comparison tests
reject changed upload bytes without extracting any member.

## Actual authority and residual risk

The workflow has explicit `contents: read` and `actions: read` permissions.
That does **not** make it side-effect-free: the official upload action uses the
run-scoped Actions runtime capability to stage one artifact. It does not have
Pages/OIDC-write permission and does not deploy the website. The private Go
command is offline; the public Go `stage` mode makes bounded GETs and creates a
new local TAR; the subsequent action is the provider upload operation.

The request and retained proof files are unsigned public integrity data. Trusted
local inputs and explicit GitHub operator/environment review are required. The
environment name in YAML is not evidence that required reviewer/branch rules are
configured. No environment, reviewer, variable, live token or provider action was
configured or exercised here; absent opt-in leaves the workflow disabled.
Existing normal candidate and Pages workflows remain unchanged.

This boundary has no current CMS Owner/source/target checks, durable staging
claim, cross-run idempotence/lifetime stage budget or CMS success receipt.
First-run rejection and serialized workflow queueing are narrower controls.
A future authenticated Owner coordinator must add those fences and audit before
automated UI dispatch/deployment, then independently verify deployed output
before atomic destination completion. Local request creation never marks a job
delivered. Old proof/code cannot silently authorize a different reviewed workflow.

Fixed API paths, strict canonical IDs and request schema prevent caller-selected
repository/origin/path fetches. API JSON is capped at 256 KiB each / 1 MiB total;
ZIP at 8 MiB, expanded data at 32 MiB, entries at 256 and central directory at
128 KiB. Requests are at most 8192 bytes. TLS/header/body timeouts and a 30-second
operation context bound reads; compression remains bounded CPU work rather than
a hard real-time interrupt. Unsupported hosts/proxies/archives fail closed.

Signed capability URLs never receive private headers, enter SQL, logs or output.
The blob transport validates all DNS answers before literal-IP dialing with
normal original-host TLS verification. Classification and transport policy are
tested locally; no real DNS/TLS/GitHub/Azure integration is claimed. GitHub/TLS,
reviewed code, the host/runner and repository administrators remain trust
boundaries. Provider state can change after final reads. Unsigned hashes are not
an approval signature or arbitrary asset malware certification.

Only trusted local CLI paths are supported. Leaf links/non-files/oversized inputs
fail, and output uses exclusive creation; concurrent hostile mutation of an
operator's ancestor directory is not sandboxed. A failed local write may leave a
partial new output, which must be discarded. Generated data stays in ignored
CMS `.tmp/` or Astro `.cms-cache/` rather than tracked source.

## Executed native and public results

The full native private `go test -count=1 -coverpkg=<unchanged handwritten CI
scope> -coverprofile=<external profile> ./...` run passed on a disposable
PostgreSQL 18 cluster. Its integration package passed in **177.627 seconds**.
Coverage measured **8617/10011 statements, 86.08%**, excluding generated SQL and
entrypoints as before. The fixture stopped and its exact validated directory
was removed normally; no editing database was read/migrated.

The standalone public Go tool passed its complete native tests at **254/292
statements, 86.99%** (Go's summary rounds to 87.0%). Scope includes the command
entrypoint and production DNS/dial branches not exercised against real networks.
Both `go vet ./...` runs passed. Pinned `govulncheck@v1.8.0` reports no
findings for the public standard-library-only module. The private backend scan
reports **0 affected call paths / 0 imported-package vulnerabilities / 1
required-module advisory without a detected called path**. That advisory remains
documented; these tools do not certify all configuration or unknown flaws.

Source-matched Astro `pnpm validate` passed: formatting/lint, **52 diagnostic
files / 0 errors / 0 warnings / 0 hints**, **25 unit tests in 6 files**, **23
repository checks** and authored static build. The repository tests include two
new TAR/device-name cases and three staging workflow cases, with the actual
pre-checkout guard executed through Python. Earlier fixture redirect, identity,
directory-expectation and CRLF/format failures were corrected and are not counted
as successful runs.

Both optional real compiled proof modes passed. Together: **6 frozen-package
Go checks, 6 exact requests, 6 pinned builds, 6 independent semantic HTML checks,
6 complete ZIP protocol/content checks, 6 Go Pages preparations and 6 Node payload
checks**, plus **6 new private staging-request generations, 6 public Go repacks
and 6 independent Node upload-TAR comparisons**. The existing **8 bilingual
desktop/mobile journeys and 16 themed Axe checks** passed. Sources stop before
build/browser reads; provider authority remains synthetic around real compiled
initial/update/empty output. No hosted wrapping/upload acceptance is claimed.

SQL/generated bindings, API/OpenAPI, Flutter, application versions and dependency
locks remain unchanged. A new public Go `go.mod` pins only the existing 1.27.3
toolchain with **no third-party library requirements**. Flutter/Android/APK and a
new Linux race run were not executed here. The new public CI job declares native
Go formatting/vet/test-race/vulnerability checks with read-only repository
permissions and no CMS/deployment secret. Declared CI is not claimed as hosted
execution. Existing private Linux race gates remain unchanged.

No new migration or app restart is needed; schema stays **25**. No commit, tag,
push, live credential/setup, hosted action, deployment or production-copy edit
was performed. Other consumers remain outside this checkpoint. See
[contract](../cms-pages-staging.md) and
[native diagram](../diagrams/cms-pages-staging.drawio).

## Final repository/security review

Both complete JavaScript dependency audits reported zero known vulnerability
counts. Fully redacted retained-history and curated Git-visible UTF-8 scans passed
for both repositories, with **1026 CMS / 185 homepage files**
in the final source inventory. Ignored dotenv, private keys, caches, generated
archives and binary fixtures are excluded. Synthetic source/env secret canaries
and tracked/staged private-binary path canaries passed. The input-placement guard
also rejects workflow input expressions outside explicitly allowed data fields,
checking all lines rather than only inline run headers. Final reports were
formatted, rescanned and checked against exact current source bytes.

Three native draw.io pages and their PNGs are mirrored byte-for-byte and visually
reviewed. Both separate untagged Custy messages passed the installed validator.
Only the CMS and Astro development repositories were edited. Original checkpoint
heads/indexes, app versions/dependency locks, schema/generated SQL/API/Flutter,
normal candidate/Pages workflows and unrelated portfolio work remain preserved.
The new public standard-library module, validation CI job and disabled staging
workflow are intentional additions. No credential, hosted dispatch/upload,
deployment, editing migration, commit, tag or push was performed. Tests and scans
provide bounded evidence; no guarantee against unknown vulnerabilities is made.
