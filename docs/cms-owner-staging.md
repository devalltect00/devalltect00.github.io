# CMS Owner staging integration

M5q, 8 November 2026. Homepage baseline `7e68be1`; private CMS baseline
`3b7ac32`. The public uploader/tool introduced in M5p is unchanged. This
checkpoint documents its new private caller; it does not enable or execute it.

The private Owner UI binds a permanent staging UUID to the explicitly selected
current candidate. It verifies current Owner/source/target/observation/content
proof, retains a claim **before provider reads**, re-downloads and checks the exact
ZIP, reconstructs the canonical payload and preflights the separate fixed staging
workflow/current main. Exact preparation and audit commit before the final
authority fence and **one** dispatch POST. No SQL lock spans network I/O.

The private database has three lifetime staging admissions per release, including
rejections and earlier candidates. Its accepted replay uses SQL without another
POST. Lost/malformed/non-success acknowledgement and expired claims remain
uncertain and block automatic fresh requests. Rejected means preparation ended
before any POST. The UI preserves retry identity, shows references and separates
acceptance from upload/publication. Latest evidence is scoped to the selected
current candidate; lifetime counts retain previous candidates.

The service uses an 18-second action budget, up to two detached seconds to record
outcome, and a 45-second permanent SQL lease. These are private coordination
rules, **not** guarantees enforced by the standalone public tool. Direct GitHub
workflow dispatch by a permitted operator bypasses the private ledger. Public
request JSON remains unsigned integrity data; never interpret it as current CMS
approval or a destination receipt.

Private migration 00026 and optional `CMS_ASTRO_STAGE_WORKFLOW_ID` are required
only in the CMS. This repository needs no migration/new dependency or runtime
CMS secret. Original candidate fingerprints are unchanged; the separate staging
target binds the uploader ID/path and candidate token incarnation. Browser
Origin/CSRF and current Owner checks remain private API responsibilities.

`cms-pages-stage.yml` stays off by default and retains main/first-attempt/exact
code guards, fixed source identities, credential-separated downloads and protected
environment. The operator must separately configure actual environment reviewers;
its name alone does not install protection. GitHub acknowledgement can precede a
skipped run if the opt-in variable is off. No Pages/OIDC-write permission is added.
Normal candidate and production Pages workflows are unchanged.

The private and provider worlds cannot share one atomic transaction. Authority
or branch code can change after private rechecks; the workflow repeats exact code
and payload checks. A successful staging run/upload is still not proof of deployed
content or permission to complete publication. Next gates are authenticated
staging reconciliation/recovery, protected deployment, independent deployed-output
verification and atomic Astro completion. No hosted run/configuration or new
Android/device proof was performed.

See the existing [public primitive contract](cms-pages-staging.md),
[new native diagram](diagrams/cms-owner-staging.drawio), and
[checkpoint evidence](security/cms-owner-staging-verification.md).
