# Deployment configuration guard

Updated 9 November 2026 (Asia/Jakarta).

Astro `scripts/pages-tool/deployment_settings.go` adds a side-effect-free guard
before future token acquisition and sender construction. Reviewed `deploymentPolicy`
is supplied separately from observed workflow environment; never construct it from
dispatch JSON, request bodies or the same untrusted environment it checks.

The guard requires explicit opt-in, GitHub Actions, workflow_dispatch, GitHub-hosted
runner, github.com/API origins, exact repository name and immutable numeric ID,
reviewed code and workflow SHA, the fixed cms-pages-deploy.yml main workflow ref,
main branch, first attempt, github-pages environment, pinned CMS HTTPS origin and
owner-scoped GitHub audience. Run IDs are canonical positive int64 strings.
Existing receiver configuration validation checks root HTTPS origin and code pins.
Returned settings contain no credential, grant or proof of CMS admission.

No token/password environment variable is read. There is no network request,
CLI registration, OIDC acquisition or workflow change. The helper does not prove
the runner environment is authentic: compromised code can forge observations.
Issuer-signed identity and current CMS authority remain mandatory. The environment
name is an observation, not proof of configured GitHub environment protection.
The selected audience is a reviewed policy constraint; actual issuer/provider
compatibility still needs verification before hosted rollout.

Tests mutate each of seventeen observations, seven policy/callback variants and
verify the happy path while failing on any credential-variable read. Standalone
Go tests/vet and full Astro validation supply local evidence; no hosted runtime,
penetration test or security guarantee is claimed. Public output/archive contract
is unchanged; no database code changed and no migration is needed (schema 35).

Next: bounded OIDC acquisition and verified runtime wiring. Deployment, completion
and ordering-hold release remain disabled. Production copies remain untouched.
