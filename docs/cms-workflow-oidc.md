# Bounded workflow OIDC acquisition

Updated 9 November 2026 (Asia/Jakarta).

Astro `scripts/pages-tool/workflow_oidc.go` implements an internal issuer client.
`newWorkflowOIDC` first checks the reviewed deployment policy and observations,
then reads only the runner OIDC endpoint and request bearer. Its separate reviewed
origin must exactly match the endpoint origin; production accepts root HTTPS on
a DNS host beneath `.actions.githubusercontent.com`, default/443 port only.
This origin is protected configuration, never derived from unsigned request JSON
or from the endpoint being validated. Live runner origin selection remains unwired.

GitHub documents `ACTIONS_ID_TOKEN_REQUEST_URL` and
`ACTIONS_ID_TOKEN_REQUEST_TOKEN` as the issuance interface in its
[OIDC reference](https://docs.github.com/en/actions/reference/security/oidc).
The client preserves runner query parameters and adds exactly one explicitly
encoded owner-scoped audience. Existing audience parameters, fragments, userinfo,
substituted origins and malformed credentials fail closed. The request bearer
goes only to the pinned issuance endpoint, never the CMS or Pages API.

Each instance has two lifetime issuance attempts: initial token and report token.
There is no implicit retry, caching or third acquisition. Cancelled/failed calls
spend budget. A newly requested token is not guaranteed to have a different JTI;
the CMS decides whether its signed identity is acceptable. TLS 1.2 minimum,
eight-second total timeout, finite headers/body, no proxy/cookies/redirects/
compression or connection reuse constrain transport and prevent reused-socket GET
replay. Only JSON 200 with one `value` field and a compact token up to 16 KiB is
accepted; response maximum is 32 KiB. No bodies, endpoints or credentials are logged
or returned in errors. Local shape checks do not verify signatures or trust claims.

CMS independently verifies issuer/RSA signature, pinned audience/workflow/code/run,
live time and durable current authority. This acquisition helper supplies no grant
and does not change CMS verification policy. Before hosted rollout, reviewed CMS
audience must match the selected audience and actual GitHub Pages acceptance must
be checked. No live audience/provider compatibility is claimed here.

Synthetic loopback tests cover status, redirect, lost reply, MIME, encoding, body
size, duplicate/unknown/null/missing/oversize/malformed tokens, cancellation,
attempt budget, endpoint substitution and credential injection. Bad context never
reads credential variables. No issuer-signed live token or joined hosted workflow
is used. Current production configuration, secrets and editing DB are untouched.

No CLI/runtime registration, workflow permissions, public output contract or schema
changes (CMS remains 35). Next: connect guarded configuration and token acquisition
to the sender in an explicitly reviewed runtime, then completion and hold release.
Hosted publishing remains disabled.
