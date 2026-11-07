# Owner candidate correlation (M5k)

Follows homepage **98f67a6** and CMS **38e29d4**. The separate candidate workflow
now accepts six required string inputs: `release_id`, `sequence`, `snapshot_digest`,
`projection_digest`, `dispatch_id` (public attempt UUID), and `expected_head_sha`
(40 lowercase hex characters identifying reviewed consumer code).

Before checkout, a Bash guard validates correlation and compares `GITHUB_SHA`
with the expected commit. Checkout uses that exact commit and does not persist
credentials. `src/core/cms/dispatch.ts` repeats the context check before the
request creator writes files. Both custom variables absent preserves standalone
local previews; a partial explicit context fails. These values are public data,
not authentication. The existing four-field exact content request is unchanged.

Run name is `CMS Astro candidate <dispatch UUID>`. Only public dist is retained
for one day under `astro-cms-candidate-<dispatch UUID>`, including its well-known
marker. The workflow retains contents-read permissions, pinned actions and a
ten-minute limit. Existing Pages deployment and authored content are unchanged.
There is no CMS callback, Owner token, private database access or deployment here.

## Optional CMS operator setup

The private CMS optionally owns a fixed repository/workflow/branch/code target,
private token-file path and append-only SQL attempt ledger (migration 00022).
The Owner screen requests a build from an active distributed approval. Callers
cannot choose an origin, workflow or code revision. GitHub repository/workflow IDs
and exact branch head are checked before one dispatch POST. Accepted run ID is
not successful build, verified artifact, deployed site or completed publication.
Three lifetime admissions, no-send rejection retry and unresolved blocking are
enforced by the CMS. Owner/source/lease checks repeat around external I/O.

The updated workflow must be reviewed, committed/pushed and available on the
repository default branch before real opt-in. Pin the new reviewed revision,
not baseline 98f67a6. Configure existing `HOMEPAGE_CMS_PUBLIC_ORIGIN` separately
as a trusted root HTTPS public reader. No real configuration or dispatch was
performed during development.

The private dispatcher credential requires repository Actions write and Contents
read. That provider permission is broader than one workflow: credential theft
could affect other workflows. Review repository protections and Pages environment
approval before opt-in. The candidate's own contents-read runtime token is
separate. See [GitHub's dispatch contract](https://docs.github.com/en/rest/actions/workflows#create-a-workflow-dispatch-event).

## Verification and remaining work

`pnpm validate` checks formatting, lint, strict Astro diagnostics, unit/repository
tests and authored build. The dispatch module has three tests for optional,
matching and malformed/partial contexts. Workflow tests enforce the guard before
checkout, exact checkout, no persisted credentials and correlated artifact names.
Existing independent Go and signed-source compiled proofs verify actual bilingual
visitor content. Read [security evidence](security/cms-owner-candidates-verification.md).

Independent run observation/reconciliation, verified artifact download,
deployed-output verification and an atomic Astro receipt remain future work.
The CMS Astro job stays pending. See [previous exact build handoff](cms-build-handoff.md)
and [editable diagram](diagrams/cms-owner-candidates.drawio).

## Follow-up: explicit run checks (M5l)

[Run observation](cms-run-observations.md) now documents private Owner metadata checks and immutable association after response loss. Accepted/uncertain dispatch history stays intact. Metadata never certifies content, deployment or destination completion; the workflow itself is unchanged.
