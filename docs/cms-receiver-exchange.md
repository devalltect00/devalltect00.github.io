# Receiver admission and report HTTPS clients

Updated 8 November 2026 (Asia/Jakarta).

The Astro standalone Go tool now contains matching clients for the optional CMS
`begin-write` and `report` machine routes. They are internal, unregistered helpers:
no CLI, workflow, OIDC acquisition, provider credential or hosted publishing is
enabled. Normal CMS startup still leaves these routes absent. No migration;
current CMS schema remains 35. Public content and archive formats are unchanged.

## Admission

`receiverExchange.beginVerified` uses the existing reviewed `receiverClient`
configuration and transport. It requires the prior exact consumed authorization
decision and original workflow bearer. Complete candidate and uploaded archive
verification, code/run separation and token bounds precede egress. The instance
spends admission even on local failure; the server's permanent ledger supplies
cross-process replay protection. No automatic retry follows denial, rate limiting,
redirect, timeout or lost reply. A new instance cannot reset the server ledger.

Only `{}` is sent to the fixed deployment-specific `begin-write` path. The response
must be explicit non-null fields, exact scope/attempt/run/artifact/request hash and
independently derived build version, with a live deadline at most ten seconds and
no later than the original authorization expiry. `validFor` must be called again
immediately before a future protected write. Local response validation itself is
not durable authority and cannot revoke an already issued network side effect.

## Reporting

`reportVerified` checks complete archives before forwarding only a bounded JSON
acknowledgement object (8192 bytes; wrapper 8256 bytes). Duplicate fields are
rejected. Fresh same-workflow authentication and actual provider acknowledgement
binding belong to the server. Local code does not claim a synthetic acknowledgement
is authentic. Returned scope/plan/operation must match exactly; recorded time is
nonzero and no more than two minutes ahead. Old timestamps are valid for historical
replay; `replayed` must be explicitly present and boolean, never absent or null.

Report calls may be explicitly repeated, returning the original stored receipt.
They never invoke admission or provider writes. The transport does not retry
automatically. A lost report does not authorize repeating a Pages create operation.

## Transport and remaining work

Both paths reuse the reviewed fixed root HTTPS origin, no proxy/cookies/redirects/
compression and TLS 1.2 minimum. Each exchange has a twelve-second HTTP context,
bounded headers, JSON 200/no-store and a response maximum of 4096 bytes. Request
body replay is disabled. All failures return a constant error, never raw remote
bodies or token material. Archive verification has its own twenty-second context;
these are separate bounded phases, not a single twelve-second total operation.

The public client fixtures use complete locally built archives and synthetic
loopback responses. CMS real RSA/HTTP/PostgreSQL checks exist separately. This
checkpoint does not claim a joined hosted workflow proof or a live deployment.
DNS/system TLS and reviewed origin/configuration remain trust assumptions; clocks
must be synchronized. Rate limits are server-side process-local limits.

Next is reviewed protected runtime and credential ownership: the separate provider
API token must not silently become the CMS Actions credential. Only then can a
protected one-shot sender join admission to the actual Pages create/report flow.
Independent verification, atomic completion and ordering-hold release remain
separate gates. No editing database, live secret or production copy was changed.

See [security evidence](security/cms-receiver-exchange.md) and [diagram](diagrams/receiver-client-exchange.drawio).
