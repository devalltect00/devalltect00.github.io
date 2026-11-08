# Receiver-side fixed HTTPS authorization exchange

Updated 8 November 2026 (Asia/Jakarta).

The independent public Go tool now has a receiver client for the existing optional
CMS machine authorization contract. It is **not exposed through CLI, workflow or
runtime**. Normal CMS startup still leaves that route absent. The code does not
mint/request GitHub OIDC tokens, upload/deploy content, configure permissions or
complete jobs. Hosted deployment remains disabled.

The client accepts reviewed receiver configuration: a fixed root HTTPS CMS origin,
canonical deployment workflow run ID, exact reviewed code SHA, main ref and first
attempt. No destination URL comes from public intent or token claims. Production
configuration rejects HTTP, IP origins, custom ports, userinfo, paths, query and
fragments. Package-local loopback HTTP exists only for test fixtures. DNS/TLS trust
is the reviewed origin's system trust; this does not authorize untrusted arbitrary
configuration. Future runtime values must come from reviewed workflow/environment,
never dispatch input or an unsigned manifest.

Before any request, it validates complete candidate/upload bytes with the local
independent verifier, checks code/run separation and a bounded compact bearer.
Then it makes one fixed POST with `{}` and Authorization only to the exact machine
path. No browser Cookie/Origin, token-selected URL, query, policy or secret input
is sent. Transport has no proxy, cookies, compression or redirect following, TLS
1.2 minimum, finite connect/TLS/header/body limits and a twenty-second total bound.
Only JSON 200 with no-store and a response body at most 4096 bytes is accepted.

Strict response validation requires exact scope/attempt/current run and canonical
typed deployment-request SHA-256. Expiry must be live and at most thirty seconds
on the trusted local clock. Unknown/duplicate/null/foreign/malformed response fails
closed. Only the exact server-derived decision may return. A later protected write
must recheck it with `validFor` immediately before use; that check alone does not
prove current SQL authority after consumption or serialize remote writes.

Each client instance atomically spends its one POST opportunity before sending.
It never retries—even on denial, 429, redirect, outage or lost reply. Invalid local
content does not send. Across process restarts the server's permanent one-use ledger
still applies; the local spent bit is not durable replay protection. A lost consumed
reply cannot safely become permission to launch another attempt. These availability
tradeoffs are intentional until explicit recovery/ordering controls exist.

## Remaining gates

No actual OIDC acquisition, authenticated runtime configuration or live receiver
exchange is enabled. Tests use synthetic local authority responses; actual private
HTTP/RSA/SQL and issuer-key validation are tested separately. A fully joined live
workflow is not claimed. Both clocks must be synchronized; skew fails closed.
Protected receiver identity/review permissions, receiver-side stale-write ordering,
Pages side effects, independent deployed-output proof and atomic completion remain
unfinished. Keep all hosted deployment gated until those boundaries are verified.

No migration (schema 32), dependencies or application-version changes. CMS and Astro
development copies have documentation/diagram updates; only the public Go client
adds production source. Production copies, live credentials and hosted settings are
untouched. Existing writers, contracts and compiled-output semantics are unchanged.

See [evidence](security/cms-receiver-https.md) and [diagram](diagrams/cms-receiver-exchange.drawio).
