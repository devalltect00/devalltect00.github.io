# Protected sender security evidence

Measured 8â€“9 November 2026 (Asia/Jakarta).

Astro `go test ./... -count=1 -cover` passed in 3.451 seconds with **90.3%**
standalone-tool statement coverage. `go vet ./...` passed. This is not full CMS
coverage or a vulnerability-free guarantee; no race-detector result is claimed.

`TestProtectedSenderExchange` has twenty joined local HTTP cases covering happy
path, complete candidate/upload rejection, cancellation, authorization/admission
denial, expired/substituted permit, provider status/redirect/lost connection/MIME/
encoding/size errors, wrong ID/URL/duplicate acknowledgement, report-token failure,
report rejection and foreign report plan. Every second sender call is rejected.
Actual acknowledgement survives report failure; explicit report-only recovery
returns a receipt while the provider create count remains one. Fixtures check
exact order, method/path/body and separation of provider vs CMS bearer headers.

`TestPagesCreatorConfiguration` checks fixed production origin, rejects unsafe
origins and seven repository/token/plan substitutions, and requires a short caller
deadline. Earlier archive and machine-client negative cases remain in the suite.

Threats addressed include credential forwarding to the wrong service, unsigned
plan substitution, duplicate local sends, redirect/body replay, unbounded responses,
expired admission and accidental create retries after report failure. Trusted
configuration/adapters, process integrity, synchronized clocks, TLS/system trust
and durable CMS admission remain assumptions. A compromised privileged workflow
can bypass local code; environment protection and reviewed runtime are still needed.
No hosted workflow, live token, penetration test or joined RSA/SQL runtime is claimed.

Full `pnpm validate` passed formatting, ESLint, 53-file Astro type checks (zero
diagnostics), 25 unit tests, 23 repository tests and static build. Native diagrams
exported successfully. Current CMS schema remains 35; no database code changed,
and no migration/integration rerun or compiled public-output proof was needed.
Retained-history Gitleaks scans passed: CMS 99 commits / 8.24 MB; Astro 23
commits / 1.10 MB. Changed visible-source scan passed (283.87 KB); final evidence
and diagram changes were included in the follow-up source scan. Diff whitespace
checks passed. Both checkpoint messages passed the installed Custy validator.
No staging, commit, tag, push, credentials or hosted deployment was performed.
