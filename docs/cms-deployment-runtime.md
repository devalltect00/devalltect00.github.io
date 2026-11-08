# Guarded internal deployment runtime

Updated 9 November 2026 (Asia/Jakarta).

Astro `scripts/pages-tool/deployment_runtime.go` now composes configuration checks,
complete local archive verification, bounded workflow token acquisition and the
protected sender. This is an internal Go entry point, not an executable command
or enabled workflow. No migration (CMS schema 35), API or public-output change.

`deploymentRuntime.run` uses production fixed-origin constructors only. Reviewed
policy and issuer origin must be supplied independently of public request data.
`runWith` has a package-local factory seam solely for trusted local HTTP fixtures;
no command, dispatch field or environment variable selects a factory/loopback mode.

The instance is permanently spent at entry, including cancellation and malformed
input. A seventy-five-second outer context bounds configuration, content checks,
issuance and the existing sixty-second sender flow. Before credential access, it
checks the nonsecret policy, strictly parses the request, checks code/run separation,
bounds and copies input archives and verifies their complete content. Requests are
parsed to independent values; caller buffers are not retained. Callers must not
mutate buffers concurrently while they are being copied.

Only then does `runtimeClients` construct separate issuer, CMS receiver and provider
transports. OIDC runner credentials go only to the reviewed issuance origin;
workflow GITHUB_TOKEN belongs only to the Pages creator. The CMS dispatch token
is never requested by this runtime. No Owner/DB credentials are loaded.

One initial acquisition feeds authorization, permanent admission and one provider
create. The second acquisition happens only after a validated provider response,
then reporting follows. Issuance, admission, creation and runtime are never retried
automatically. A validated acknowledgement survives report failure for explicit
report-only recovery. All adapters are closed on success/failure without resetting
budgets. No raw errors, credentials or acknowledgement data are logged.

Tests join complete real archive fixtures to actual HTTP clients for issuer, CMS
and provider on loopback, with synthetic decisions/tokens and separately checked
credential headers. They are not joined issuer-signature/PostgreSQL or hosted proof.
Environment checks are observations, not signed authority. CMS signature/current
SQL admission remain mandatory and unchanged. Compromised privileged runtime code
can bypass helpers; protected environment/code review are still rollout gates.

Next: an explicitly reviewed executable entry and acquisition of exact local input
artifacts, followed by independent terminal/output verification joined to atomic
completion and ordering-hold release. No CLI, hosted workflow, permissions, live
credentials, production copies or editing database were changed. A returned report
is historical operation evidence, never deployment completion or hold discharge.
