# Exact Astro build handoff (M5j)

Follows CMS **5f09c5d** and homepage **3687134**. This checkpoint binds a static
candidate to the exact requested release and independently verifies its visitor
content with Go. It does **not** add an Owner API action, hosted deployment,
provider credential, database migration or durable Astro destination receipt.
The two development repositories change; portfolio and production copies do not.

## Why another check after the build marker?

The marker is readable public JSON. Copying its identity into unrelated/stale HTML
must not prove success. Go now derives the expected Astro projection from the real
frozen schema independently of TypeScript, then checks actual semantic elements in
both pages: selected project IDs/order, localized title/summary/status, version and
technology labels, optional link kinds/URLs, Profile name/focus and ordered public
contacts. Empty/removal states are authoritative, including absence of old email.

Scripts, comments and serialized props cannot substitute for visible semantic
elements. Explicit hidden/inert markup and inline styles are rejected, except the
known numeric project animation-delay custom property. The verifier is deliberately
for these fixed templates: it is not a general hostile-CSS renderer, continuous
monitor or proof that a hosting account is uncompromised.

## Public request contract

```json
{
  "schema_version": 1,
  "scope": "astro-project-profile-build-v1",
  "release_id": "11111111-1111-4111-8111-111111111111",
  "sequence": "1",
  "snapshot_digest": "0000000000000000000000000000000000000000000000000000000000000000",
  "projection_digest": "0000000000000000000000000000000000000000000000000000000000000000"
}
```

Values above illustrate shape only; generate a real request from the approved
public package. Sequence remains decimal text through the full signed int64 range.
Duplicate/unknown keys, malformed UTF-8, zero/leading-zero/overflow sequences and
noncanonical hashes fail. This intent contains **no authority, URL, provider token,
secret or database connection**. Matching public hashes do not authenticate approval.

The existing output marker keeps `astro-project-profile-preview-v1`; request scope
is a different contract. Old marker-only pages lack the semantic evidence required
by the new verifier and therefore fail; no old receipt is upgraded or deleted.

## Native local procedure (PowerShell)

Use an already prepared anonymous package from the Astro guide. All placeholder
paths/UUIDs below must be replaced with actual local public files. No private CMS
environment loader, Owner password or database connection is needed.

```powershell
cd D:\project\develop\devalltect-cms\backend
go run ./cmd/cms astro-build request --package D:/path/to/public-package.json --out E:/path/to/new-build-request.json

cd D:\project\develop\devalltect00.github.io
pnpm content:cms:prepare --origin https://your-approved-public-reader.example --request E:/path/to/new-build-request.json
$env:HOMEPAGE_CMS_RELEASE = '<exact-requested-release-uuid>'
$env:HOMEPAGE_CMS_REQUEST = 'E:/path/to/new-build-request.json'
pnpm build

cd D:\project\develop\devalltect-cms\backend
go run ./cmd/cms astro-build verify --package D:/path/to/public-package.json --request E:/path/to/new-build-request.json --dist D:/project/develop/devalltect00.github.io/dist
```

`request` exclusively creates a new file in an existing operator-selected directory;
it never replaces an existing file. `verify` is read-only and returns public JSON
with exact request, selected project/page counts and SHA-256 for the EN/ID HTML and
marker. It loads no CMS settings and makes no network call. There is no API or DB
restart requirement; use the updated native CLI rather than an older image binary.

The operator chooses trusted local package/request paths and workspace. Package
reads are capped at 2 MiB + 8192 bytes, request/marker at 8192, each HTML at 2 MiB.
Non-files, known symlinked input leaves and linked directories under the selected
dist root are refused. Only `index.html`, `id/index.html` and
`.well-known/homepage-content.json` are read from dist; no URL in content/HTML is
followed. The check does not enumerate or attest all CSS, scripts, images or other
files in a deployment artifact. It is not a race-proof local filesystem sandbox.

## Exact preflight and an advancing public source

Preparation with `--request` checks all four identity values **before** writing
cache files. If the reader has moved from requested release A to B, the command
fails; it must not silently build B while reporting A. Existing prepared releases
remain intact. A trusted operator can build an already retained A offline using
its explicit pin and matching intent. Arbitrary historical remote fetch/rollback
is not added. The existing reader only exposes its active public snapshot.

`HOMEPAGE_CMS_REQUEST` is a process-environment filename, alongside the existing
process `HOMEPAGE_CMS_RELEASE`. Blank/missing/corrupt/mismatched explicit intent or
intent without a CMS pin fails **before** Astro can replace dist. The central
frontmatter adapter repeats the check for direct Astro entrypoints. No `.env`
autoload or `PUBLIC_` selector is added. Preview pins without a request still work
for M5i compatibility; automated candidates must use the exact intent. Clear both
selectors deliberately to return to authored mode.

## Optional manual candidate workflow

The homepage adds `.github/workflows/cms-build-preview.yml`, separate from the
unchanged Pages deployment. If later enabled in the development repository, an
authorized repository operator may manually supply exact public ID, sequence and
digests. The only source origin is the operator-set repository variable
`HOMEPAGE_CMS_PUBLIC_ORIGIN`, a root HTTPS public reader. It is not a workflow input.
Shell steps receive inputs as environment data and never interpolate them as code.

The job creates a strict new public request, prepares only that exact release,
builds it and retains `dist/` as `astro-cms-candidate` for **one day**. Hidden-file
inclusion is limited to that public dist root so the `.well-known` marker survives.
No cache, request directory or private repository root is uploaded. Actions use
full commit pins. Repository permissions are `contents: read`, with a ten-minute
job limit; the runner's artifact service intentionally stores the candidate.
There is no Pages/OIDC/contents-write permission, deploy action, CMS callback or
Owner token. The private Go verifier is not shipped into the public repository:
the resulting artifact is a **candidate requiring independent verification**.

No workflow was dispatched, origin configured, token created, artifact uploaded
or site deployed during this local checkpoint. A future Owner dispatcher still
needs fixed repository/workflow/ref and run identity, durable attempts/leases,
active-source and Owner checks, verified deployment output, retry/recovery and an
atomic receipt/audit. No local report or GitHub artifact is sufficient to complete
an Astro job by itself. Jobs remain pending.

## Implementation and verification

| Responsibility                                                             | Source                                                          |
| -------------------------------------------------------------------------- | --------------------------------------------------------------- |
| Independent Go normalization, exact projection, Unicode excerpt and intent | CMS `backend/internal/platform/astrobuild/plan.go`              |
| HTML5 parsing, bounded raw guards and semantic output                      | CMS `backend/internal/platform/astrobuild/html.go`, `markup.go` |
| Bounded file reads and read-only CLI                                       | CMS `backend/internal/administration/astroverify/`              |
| Exact TS request parsing/comparison                                        | Homepage `src/server/cms/request.ts`                            |
| Early failure and repeated build selection                                 | Homepage `scripts/content/`, `src/data/homepage.ts`             |
| Semantic attributes on actual visitor elements                             | Homepage project/Hero/contact/footer components                 |
| Real Go requests plus output checks in both public-source proofs           | Homepage `scripts/testing/prove-cms-build.mjs`                  |
| Candidate authority/input policy gate                                      | Homepage `scripts/repository/tests/cms-workflow.test.mjs`       |

Run the two compiled proof commands from the earlier Astro preview guide. They now
generate requests with the real Go CLI and independently verify every initial,
updated and empty compiled output. Standalone public CI runs the same request
preflight using public synthetic intent; cross-repository Go verification runs
locally/private-side. Both sources are stopped before builds/browsing.

See [security evidence](security/cms-build-handoff-verification.md) and the native
[three-page diagram](diagrams/cms-build-handoff.drawio).
