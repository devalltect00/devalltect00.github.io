# Astro protected publishing and completion

Current implementation, 9 November 2026. This guide supersedes earlier local-only
deployment prerequisites; their dated evidence remains historical. Source and
local tests do not prove a hosted rollout. No repository variable, environment,
live credential or hosted workflow was enabled during this development work.

## Components and authority

The Owner uses Releases → selected Astro candidate → verified Pages upload →
deployment panel. Candidate building and artifact staging remain distinct from
deployment. `cms-pages-deploy.yml` accepts only the exact server-generated public
request and permanent UUID. A protected, separately reviewed policy pins repository
ID, workflow, source SHA, receiver origin, OIDC audience and issuance origin.
Unsigned JSON and artifact hashes are integrity data, not CMS permission.

The Go workflow tool downloads the exact candidate and uploaded wrapper with
read-only artifact authority, checks provenance before and after download and
independently verifies complete archive contents without extracting them. Only
then does it acquire workflow identity and ask the CMS for current permission.
CMS verifies fixed-issuer RSA identity, exact run/code/repository/environment,
current Owner, immutable approved source and target policy. Durable one-use
admission precedes a single Pages create call; no automatic write retry exists.

Provider I/O occurs outside SQL transactions. Admission, completion and retirement
repeat current authority under the shared workspace/Owner lock. Public file GETs
never receive private GitHub credentials. All machine and Owner responses retain
the existing bounded parsing, constant error codes and private no-store policy.

## Database and default configuration

Migration **00036** adds immutable completions, explicitly sourced reconciled
receipts, immutable no-write retirements and a guarded current-attempt pointer.
Existing schema-35 claims are backfilled into that pointer. All original dispatch,
claim, receipt and audit history stays retained. The API requires schema 36; it
does not apply migrations implicitly. Use the existing migration procedure with
the API stopped; do not modify migration history manually. A downgrade refuses
retained completion or retirement evidence.

Deployment stays disabled when both new settings are absent:

- `CMS_ASTRO_DEPLOY_WORKFLOW_ID`: canonical numeric ID of the distinct deployment workflow.
- `CMS_ASTRO_OIDC_KEY_ID`: operator-reviewed GitHub OIDC signing key ID. This is a public identifier, not a secret. Key rotation needs policy review; an unknown/rotated key fails closed.

The existing candidate/staging configuration and private token file are still
required. The fixed audience is `https://github.com/<repository-owner>`; the
subject is `repo:<owner>/<repo>:environment:github-pages`. Never put token contents
in dotenv, workflow dispatch inputs, Git, screenshots or documentation.
The current local production-like Compose profile does not expose the optional
publishing configuration or mount its credential file. Use the existing native
API environment loader for local configuration; production mounting is part of
the separately reviewed deployment handoff, not an implicit Compose enablement.

## Protected workflow setup for a future rollout

Keep `HOMEPAGE_CMS_DEPLOY_ENABLED` unset until review. The proposed environment
variables are read from the `github-pages` protected environment/repository:

| Variable                             | Reviewed value                                            |
| ------------------------------------ | --------------------------------------------------------- |
| `HOMEPAGE_CMS_DEPLOY_ENABLED`        | Literal `true` only during an authorized rollout          |
| `HOMEPAGE_CMS_RECEIVER_ORIGIN`       | Fixed public HTTPS CMS origin; no arbitrary path/query    |
| `HOMEPAGE_CMS_OIDC_AUDIENCE`         | Exact owner-scoped audience above                         |
| `HOMEPAGE_CMS_CANDIDATE_WORKFLOW_ID` | Exact candidate workflow numeric ID                       |
| `HOMEPAGE_CMS_STAGE_WORKFLOW_ID`     | Exact, distinct staging workflow numeric ID               |
| `HOMEPAGE_CMS_DEPLOY_POLICY`         | Reviewed JSON described below                             |
| `HOMEPAGE_CMS_DEPLOY_POLICY_SHA256`  | SHA-256 of the exact UTF-8 JSON bytes; whitespace matters |

Policy fields: `schema_version: 1`, `scope: "astro-protected-runtime-policy-v1"`,
`repository`, `repository_id` (decimal string), `workflow_ref` (exact
`<owner>/<repo>/.github/workflows/cms-pages-deploy.yml@refs/heads/main`), `head_sha`,
`receiver_origin`, `audience`, and `issuer_origin` (the reviewed GitHub Actions
issuance origin). Confirm the literal scope against the parser before rollout;
configuration is never copied from public dispatch input. Source SHA and policy
digest must be reviewed again when the consumer code changes.

The workflow requires main, first attempt, GitHub-hosted runner and `github-pages`.
Only its deployment job has Pages/OIDC write permission. Checkout credentials are
not persisted. The ordinary authored Pages workflow skips its build while CMS
publishing is enabled and uses a separate disabled concurrency group. Drain all
pre-existing authored Pages runs before switching modes; disabling a future run
cannot undo one already executing. Pages must use the GitHub Actions source.

## Owner workflow and recovery

1. Refresh the selected release's deployment panel. It shows durable SQL progress,
   permanent attempt/run IDs and the lifetime budget; it does not poll or grant
   provider authority merely by displaying a status.
2. Confirm one deployment request for the verified upload. Accepted means GitHub
   acknowledged dispatch; it does not mean the site has deployed. An uncertain
   response must not be blindly resent or replaced.
3. After the workflow terminates and the short write admission expires, choose
   **verify/complete**. The CMS independently verifies exact run provenance,
   deterministic Pages operation success, complete candidate/upload bytes and
   every expected public file, with stable metadata before and after inspection.
4. Only that final transaction retains verification and completion, marks the
   Astro job succeeded, updates aggregate release state and discharges that
   attempt's remote ordering hold. Other consumer jobs remain unchanged.

A workflow can fail after Pages successfully publishes, for example when its
callback is lost. Completion may reconcile that exact successful operation even
when the workflow conclusion is failed/cancelled. It must still be terminal and
all independent checks must pass. A missing receipt is stored explicitly as
`owner-reconciled`, never represented as a successful workflow callback.
Completed replay is read-only and performs no new provider write or verification.

**No-write recovery:** close a never-dispatched claim, or select its exact terminal
workflow run when dispatch acknowledgement was lost. CMS independently checks
repository/code/title/first-run provenance and repeats the absence of any Pages
write admission, report or completion under lock. This permanently retires old
permission and removes only the current pointer. A separate confirmation may
create a fresh UUID within the unchanged **three-claim lifetime budget**.
No history is deleted, run cancelled or old attempt resent. Live/ambiguous runs,
any admitted write, stale authority and exhausted budgets remain blocked.

Do not clear SQL rows, reset counters, fake success or mark other consumers done
to remove a hold. Artifact expiry (currently one-day workflow retention), GitHub
outages, unavailable public bytes, key rotation or exhausted budgets can require
operator investigation. This bounded recovery is not universal automatic repair.

## API and operational limits

Owner routes: `GET/POST /api/v1/admin/releases/{id}/astro-deployment`,
`POST .../complete`, and `POST .../retire`. Browser mutations require current
session, permitted Origin and CSRF. Editors cannot invoke them. Refer to OpenAPI
and `the private CMS REST Client examples` for strict request bodies.

One process-local slot bounds Owner deployment mutations; permanent SQL ledgers
and workspace locks protect cross-process authority. Dispatch is bounded to 40s,
completion to 65s and retirement to 20s. Flutter response limits are respectively
50s/75s/30s; API write timeout is 90s. Ordinary client requests keep their existing
15s limit. Do not configure a proxy with a shorter completion response deadline.

Safe logs use constant event names and identifiers: `astro.deployment_configured`,
`astro.deployment_accepted`, `astro.delivery_verified`. Transactional audit records
include completion, reconciled receipts and retirement. Never log raw provider
errors, request archives, policy blobs or token material. Use release/attempt and
request-reference IDs to correlate an incident with the preserved evidence.
