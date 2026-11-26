# Hosted CMS identity preflight

## Immutable subject correction — 27 November 2026

The repository's read-only OIDC policy reports `use_immutable_subject:true`.
The preflight therefore pins this exact public identity:
`repo:devalltect00@62181469/devalltect00.github.io@1389151718:environment:github-pages`.
Owner/repository IDs and environment/runner/branch metadata are also compared
independently. Legacy or foreign identities are rejected. Public metadata now
includes the reviewed subject and owner ID; JWTs and capabilities stay in memory.

Failed runs report only fixed reason codes, such as `CLAIM_SUB`,
`ISSUANCE_REQUEST_FAILED` or `UNEXPECTED`, after `OIDC_PREFLIGHT_UNVERIFIED:`.
No token, claim value, URL or upstream diagnostic is echoed. The original generic
failure cannot establish which guard ran last, but its legacy subject expectation
is incompatible with the verified GitHub policy. This corrects the reader policy;
it does not change GitHub identity settings or grant deployment permission.

Commit/push this correction, confirm new-source CI, then create a **new** manual
first-attempt preflight with the **new full main SHA**. Do not rerun a prior attempt
or use historical SHA `4098e40` for corrected code. The private CMS must install the
corresponding reviewed `CMS_ASTRO_OIDC_SUBJECT` and patched API through its later
operator handoff before protected publication. No migration is required.
See [GitHub immutable subject documentation](https://docs.github.com/en/actions/reference/security/oidc#immutable-subject-claims).

The homepage keeps its existing static English `/` and Indonesian `/id/` routes.
The protected [publishing workflow](cms-publishing.md) remains the sole CMS Pages
writer. This guide adds an operator bootstrap check, not a new publishing authority.

## Before enabling publishing

1. Commit, review and publish the complete development source to `main`. Candidate,
   staging, deployment and identity-preflight workflows must exist on that exact
   default-branch commit. A local commit alone is not a remote workflow.
2. Preserve the existing live Pages deployment. Before pushing setup changes,
   drain/pause the ordinary authored Pages workflow so it cannot overwrite the
   site during the switch. Preserve CI and record how to restore the prior state.
3. Use a separate repository-scoped CMS fine-grained token: Actions read/write,
   Contents read-only and short expiry. Keep it outside Git and outside workflow
   inputs. Actions-write is repository-wide, not limited to these workflows.
4. Run the deployment-free identity preflight below. Review its public signing-key
   ID and issuance origin against the fixed GitHub issuer/public keys. Never guess
   the active key merely from a list containing multiple keys.
5. Pin the real workflow numeric IDs, repository ID, exact reviewed source SHA,
   fixed CMS origin, audience and observed issuance origin in the existing policy.
   Provision runtime settings/private token through the trusted operator channel.
6. Only then enable the separately reviewed candidate/stage/deployment controls.
   Build, verify, stage, authorize, deploy and independently complete through the
   existing Owner flow. A successful workflow or public marker alone is no receipt.

The Pages tool now pins Go 1.27.4. Its existing Linux CI job uses that module file
and already runs race/coverage, vet and pinned govulncheck. The Go patch includes
[GO-2026-6859 remediation](https://pkg.go.dev/vuln/GO-2026-6859); that advisory is
server-specific, so a patch pin is not evidence that this command was exploited.

## Run the identity-only workflow

In GitHub Actions choose **Inspect CMS workflow identity (no deployment)** and
run `cms-oidc-preflight.yml` on `main`, passing the exact full reviewed main commit
as `expected_head_sha`. The main/first-run/GitHub-hosted checks execute before
credentials are read. The job uses the existing `github-pages` environment and
preserves its branch protections. It has Contents read and id-token write; it has
**no Pages-write, repository-write or CMS credential**.

The script requests one ephemeral OIDC token from the runner's validated GitHub
issuance endpoint, reads the fixed public GitHub JWKS without forwarding the
request capability, and records only allowlisted public metadata. It prints no JWT,
request capability, signature, arbitrary claims or upstream error text. JSON is
finite/bounded/duplicate-free; redirects/proxies/cookies are disabled. The public
metadata artifact expires after one day.

The report includes the observed `oidc_key_id`, query-free `issuer_origin`, exact
code/repository/run identity and explicit `authorization:false` and
`signature_verified:false`. This helper parses observations; it is not a JWT
signature verifier. The existing CMS must independently validate the real signed
identity, deployment workflow, current Owner, release, lease and one-use admission.
The preflight workflow identity cannot satisfy the deployment-workflow pin.
An origin/key change on a later run fails closed and requires operator review.

Never rerun a publishing attempt merely because this bootstrap metadata exists.
The preflight changes no CMS job, consumes no CMS admission and never calls Pages.

## Local checks

```text
pnpm validate
python -B -m unittest discover -s scripts/testing -p test_oidc_preflight.py
cd scripts/pages-tool
go test ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
```

`pnpm validate` already discovers the Node wrapper that runs the Python security
cases, so the second command is an optional focused check. Tests cover rejecting
stale/foreign contexts before credential/network access, unsafe issuance URLs,
ambiguous/oversized JSON, wrong claims/time/header/key metadata, capability isolation
and the explicit absence of authorization/signature claims. They use synthetic
identities and no external credentials/network.

The private CMS operations journal owns actual hosted setup/dispatch/receipt
evidence. Source instructions do not claim those actions ran. No production copy
under `D:/project/publish` is part of this development workflow. Editable overview:
[cms-hosted-preflight.drawio](diagrams/cms-hosted-preflight.drawio),
[PNG preview](diagrams/previews/cms-hosted-preflight.png).

Upstream references: [GitHub OIDC claims and permissions](https://docs.github.com/en/actions/reference/security/oidc),
[fine-grained Actions permissions](https://docs.github.com/en/rest/authentication/permissions-required-for-fine-grained-personal-access-tokens).
