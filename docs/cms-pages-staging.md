# Public exact Pages staging contract (M5p)

Follows CMS **67879ad** and homepage **9bd1121**. This checkpoint implements the
**staging upload boundary**, before deployment: an optional manual workflow
acquires an exact candidate, reconstructs the pinned private Go payload and
uploads a single prepared `artifact.tar` through the official Actions uploader.
The workflow is disabled unless explicitly configured and has no Pages or
OIDC-write permission. No hosted workflow or provider mutation was run here.

## What is implemented

- Private `cms astro-build stage-request` checks a frozen approved package,
  exact request, retained public artifact proof and actual candidate ZIP. It
  reconstructs the Pages payload itself rather than accepting a caller's payload
  hash. Proof/archive/frozen-content evidence must match exactly.
- A compact public `astro-pages-stage-v1` request binds a new correlation UUID,
  candidate UUID, lossless run/artifact IDs, reviewed code SHA, whole ZIP/gzip
  digests and exact build identity. Strict unknown/duplicate/type/size rejection
  is shared as a cross-repository contract.
- The public standard-library-only Go tool under `scripts/pages-tool/` validates
  workflow opt-in/context, fixed repository/workflow/code provenance, completed
  successful first-run metadata, unexpired exact artifact and stable rechecks.
- GitHub API credentials and signed blob downloads have separate clients. The
  latter uses canonical HTTPS artifact hosts, public-IP pinning and no forwarded
  Authorization/Cookie/Referer, proxies or secondary redirects.
- Canonical in-memory ZIP-to-USTAR/gzip reconstruction must match **both** pins.
  All regular-file bytes, including `.well-known` and binary assets, survive.
  The output TAR is exclusively created after final provenance checks.
- `cms-pages-stage.yml` uploads that **one TAR file**, not a reconstructed site
  directory. The public Node checker also compares the actual upload TAR with
  the independently verified Go payload. No member is extracted or executed.

## Local request generation

This is an offline operator command, not a new admin API or UI action. Save the
direct `data.latest` proof from the existing Owner artifact GET into a private
local file. Do not include the surrounding response envelope or any session token.
Use the trusted approved package/exact request and corresponding candidate ZIP;
keep inputs/outputs in the ignored CMS `.tmp/` or Astro `.cms-cache/` directory.

From the private `backend/` directory, substituting local file paths:

```powershell
$stageID = [guid]::NewGuid().ToString()
go run ./cmd/cms astro-build stage-request --package approved-package.json --request exact-request.json --zip candidate.zip --artifact-proof artifact-proof.json --stage-id $stageID --out stage-request.json
```

Existing output files/links are never replaced. Invalid proof, source pin,
frozen facts or IDs fails before output. A partial new file after an I/O failure
must be discarded before retrying. These are trusted operator-selected local
paths, not a sandbox against concurrent host-user directory changes.
No CMS configuration, database, credential or network request is loaded by this
command. The resulting public JSON has no private token, URL or actor profile.

From the Astro repository, the offline compatibility check is:

```powershell
go -C scripts/pages-tool run . pack --request stage-request.json --zip candidate.zip --out artifact.tar
node scripts/content/verify-pages-payload.mjs candidate.tar.gz pages-evidence.json artifact.tar
```

Use absolute input/output paths when invoking Go with `-C`: relative paths are
resolved from the tool directory. The optional third Node input checks complete
TAR equality, without extraction. Pinned Go **1.27.3** is used by both writers;
cross-version compressor identity is not promised. No external Go dependency or
private repository is needed by the public tool.

## Workflow boundary and future hosted setup

`.github/workflows/cms-pages-stage.yml` has only `workflow_dispatch`. It requires
all of: `HOMEPAGE_CMS_STAGE_ENABLED == 'true'`, `main`, first workflow attempt,
fresh syntactically valid stage UUID, exact request/code correlation and the
`cms-pages-staging` environment. The before-checkout Python guard rejects
duplicate JSON keys and wrong ID/code/ref/attempt without echoing supplied data.
Checkout pins `github.sha` with no persisted credentials; actions are SHA-pinned.

Future operator setup, **not performed by this checkpoint**, requires a reviewed
public commit containing the workflow/tool, protected `main`, the
`cms-pages-staging` environment with appropriate reviewer/branch rules, and
repository variables `HOMEPAGE_CMS_STAGE_ENABLED` and
`HOMEPAGE_CMS_CANDIDATE_WORKFLOW_ID`. The latter is the numeric ID of the fixed
candidate workflow. Declare reviewer rules explicitly: naming an environment
does not configure or prove those protections. A sole administrator reviewing
their own run is not four-eyes approval. No environment, reviewer or variable was
created here. The feature stays off until an operator separately authorizes setup.

The request's reviewed SHA must equal the staging workflow commit; old candidate
proofs cannot silently authorize new workflow code. A future hosted candidate
must be built at that same reviewed commit. No normal Pages/candidate workflow
is changed or dispatched by this implementation. The private permanent token is
never sent to the public workflow; its job token is scoped to the current repo
with `contents: read` and `actions: read`.

**Read permissions do not mean the workflow is side-effect-free.** The official
artifact action uses the run-scoped Actions runtime capability to upload
`astro-cms-pages-<stage UUID>`, containing only `artifact.tar`, compressed at
level 6 and retained one day. It has no `pages: write`/`id-token: write` and
cannot use this boundary to deploy Pages. GitHub's official
[upload-pages-artifact implementation](https://github.com/actions/upload-pages-artifact/blob/v4/action.yml)
also uploads a TAR file through `upload-artifact`. We prepare the checked TAR
ourselves to preserve all approved members, not follow links or exclude the
public marker. Hosted service wrapping/acceptance is not claimed as tested.

## Resource and trust limits

Request JSON is capped at 8192 bytes. Source ZIP is capped at 8 MiB, expanded
content at 32 MiB, inventory at 256 and central directory at 128 KiB. Metadata is
at most 256 KiB per body / 1 MiB total. The GET/repack operation has a 30-second
context, finite connection/header/body deadlines and an overall five-minute job
limit. Cancellation is checked between bounded reads/files; compression is not
a hard real-time interrupt. Only ordinary portable regular files and empty
directory metadata are accepted; no links, special files, aliases or PAX output.

API paths are fixed to `api.github.com` and the workflow's own repository.
Signed download origins are restricted to canonical Azure artifact storage or
the exact documented Actions receiver host; all resolved addresses must be
public, then one checked literal is dialled with normal hostname TLS validation.
Unsupported hosts/corporate proxies fail closed. URLs remain transient secrets
and are never printed. Safe diagnostics are constant `PAGES_STAGE_UNVERIFIED`.

The public request/proof files are unsigned integrity data, not authenticated
approval. Trusted local inputs and GitHub environment/operator review are
required. This boundary does **not** check current CMS Owner/source/target state,
persist a deployment claim, enforce a cross-run staging lifetime budget, or
provide cross-run idempotence. First-attempt rejection and serialized queueing
are not those guarantees. GitHub state can change after rechecks; SQL and provider
state are not one transaction. A compromised repository administrator/runner,
reviewed code or trusted provider remains outside this local proof.

## Verification, upgrades and next work

No API, Flutter screen, OpenAPI route, database schema or app version changes.
Migration **25** remains current; no migration or app restart is needed. Existing
candidate recovery/three-admission history remains untouched. Public CI adds a
credential-free Go vet/test/race/vulnerability job for the standalone tool; normal
public validation does not need the private repository or CMS credentials.

The optional genuine initial/update/empty proof now runs the private request
generator, public Go repacker and independent Node upload-TAR check in both
independent and signed-source modes. Provider authority is synthetic. See
[security evidence](security/cms-pages-staging-verification.md) and
[native diagram](diagrams/cms-pages-staging.drawio).
Next: the authenticated CMS Owner coordinator with durable claims/audit and
current source/target/proof fences, then protected deployment, independent
deployed-output verification and atomic Astro destination completion.
