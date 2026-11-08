# Workflow-owned protected sender boundary

Updated 9 November 2026 (Asia/Jakarta).

The Astro Go tool now has an internal `protectedSender` joining real fixed-origin
HTTP clients: complete archive verification, consumed authorization, permanent
write admission, one Pages create, then authenticated operation reporting.
There is no CLI command, runtime environment loader, workflow or hosted write.
Default CMS startup still leaves machine routes absent. No migration (schema 35).

## Credential ownership decision

The intended Pages creator runs in the protected Astro workflow. Its ephemeral
repository `GITHUB_TOKEN` belongs only to the provider API adapter; it never goes
to the CMS. The CMS Actions dispatch credential stays private and is not reused
for Pages writes or passed to a workflow. Workflow OIDC authenticates the CMS
exchanges and accompanies the provider create body. No database, Owner password,
TOTP seed or long-lived CMS session enters the public repository/workflow.

Future runtime must pin repository identity, reviewed code, main/first attempt,
the `github-pages` environment, exact workflow and CMS HTTPS origin independently
of unsigned request data. It must acquire OIDC under an audience policy compatible
with both GitHub Pages and CMS verification; this is not implemented or assumed
proven here. Configure protected environment review and minimum job permissions
(`pages: write`, `id-token: write`, plus reviewed artifact-read needs) only during
explicit hosted rollout. Fork/PR code must never receive these permissions.

This choice follows [GitHub's custom workflow guidance](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages).
The CMS prototype creator remains unconstructed; this checkpoint selects the
workflow-local implementation for future runtime rather than importing private
Go internals into the public consumer. The two transport implementations currently
share protocol expectations but not code; future protocol changes must test both.

## Implemented sender

`scripts/pages-tool/protected_sender.go` owns a sixty-second total context and
permanently spends the instance at entry, including invalid/cancelled calls.
It checks the creator's pre-pinned plan before authorization, acquires the original
consumed authorization and exact single-use admission, then rechecks live permission
immediately before a create request bounded by that permission's deadline.

`pages_create.go` pins reviewed repository segments, provider token and exact
deterministic operation plan. Its sole production origin is `https://api.github.com`;
loopback override exists only in package-local fixtures. It requires a caller
deadline within ten seconds and has its own eight-second HTTP timeout. TLS 1.2+,
no proxy/cookies/redirects/compression, bounded headers and response, disabled body
replay, and one-shot instance consumption constrain the mutation.

The [Pages REST contract](https://docs.github.com/en/rest/pages/pages#create-a-github-pages-deployment)
is used with exactly integer artifact ID, `github-pages`, deterministic build version
and OIDC token. Only JSON 200 at most 8192 bytes is accepted. Duplicate JSON,
wrong operation ID, foreign status URL or malformed responses fail closed. The
expected status URL is constructed locally, never followed from response data.

After a validated create response, a trusted callback supplies a report bearer.
Only then is report attempted. The result retains an in-memory copy of the actual
validated acknowledgement even if reporting fails, so callers can explicitly retry
only `reportVerified`. Never log this result or serialize it into public artifacts.
No automatic retries occur. Invalid/lost create replies retain no acknowledgement
and require separate operator reconciliation; a new process cannot reset SQL admission.
Timeout/cancellation cannot undo a request already accepted by GitHub.

## Limits and next gates

Joined tests use synthetic CMS decisions and real loopback HTTP on both origins,
complete real archive fixtures and fake credentials. They do not prove issuer/RSA/
SQL integration in one hosted runtime. Existing private authentication/storage tests
remain separate. No public-output contract changed, so compiled content proof is
unchanged. No current editing DB or production copy was touched.

Next: bounded protected runtime configuration and OIDC acquisition, followed by
independent terminal/public-output proof joined to atomic completion and hold release.
No successful sender result certifies deployment completion or releases ordering.

[Configuration guard](cms-deployment-settings.md) is implemented internally; token acquisition remains separate.
