# GitHub provider-token compatibility

Updated 28 November 2026. This is a compatibility repair for the existing Pages
tool, not a new publisher, permission, endpoint or authentication scheme.

GitHub now issues opaque installation credentials in a longer format containing
dots and hyphens. The earlier 32–512 alphanumeric/underscore check rejected valid
modern credentials before provider I/O. The restriction existed in staging and
the shared acquisition/Pages creator paths.

[GitHub's format announcement](https://github.blog/changelog/2026-04-24-notice-about-upcoming-new-format-for-github-app-installation-tokens/)
and [completed rollout](https://github.blog/changelog/2026-10-02-stateless-github-app-installation-tokens-rolled-out/)
describe approximately 520-byte, variable-length installation credentials and
advise treating them as opaque strings. No real hosted token was retained to
inspect its claims or exact shape.

`scripts/pages-tool/provider_token.go` owns one validation function used by all
three clients. It accepts 32–4096-byte bearer-safe ASCII, including modern dots,
hyphens and padding characters; it rejects whitespace, controls, non-ASCII and
oversized values. The fixed GitHub API receives the credential unchanged and
authenticates it. The code does not decode or verify the credential as a JWT.
The separate workflow OIDC signature/issuer/subject/audience and source/run/SQL
authorization checks remain intact. No token or raw provider error is printed.

Synthetic tests cover legacy and longer formats, exact minimum/maximum, overrun,
CRLF, whitespace, NUL, DEL and Unicode; all three entry points and unchanged fixed
origin/header forwarding are exercised without provider requests or real secrets.
The isolated legacy-policy negative control fails those regressions; the fixed
complete tool suite and vet pass. Full `pnpm validate` also passed (25 Vitest
tests, 26 repository tests, static build), followed by the joint compiled
Go/Astro fixture proof including both locale builds and protected contract guards.
Exact-source CI remains required before changing any hosted source pin.

```text
go -C scripts/pages-tool test ./...
go -C scripts/pages-tool vet ./...
pnpm validate
```

Hosted rollout remains paused: source8f still owns the current candidate, both
writer flags are false, and two failed staging runs remain immutable. No Pages
admission or website replacement occurred. A new commit cannot reuse an old
source proof: finish legitimate old-work recovery before a reviewed source repin
and fresh candidate. Do not rerun old workflow attempts, delete artifacts to force
recovery or reset the three-admission lifetime budget. The private CMS operation
journal holds actual release IDs/receipts and remaining gates.

The source fix can be committed and pushed for CI while the writer flags stay
false. Keep the old server runtime/policy source pins until legitimate recovery:
historical-run observation checks that exact run, not the current main tip. Both
failed staging runs have been retired; candidate retirement remains separate.
