# Candidate recovery contract (M5n)

Follows homepage **fa81d01** and private CMS **fbe9010**. This checkpoint changes
only homepage documentation and native diagrams. Candidate/Pages workflows,
public source, authored records, optional compiled proof hook, dependencies and
versions remain unchanged. Private CMS adds migration **00025**, explicit Owner
recovery endpoints, immutable evidence and Flutter controls.

## What recovery permits

An Owner first records the exact completed first-run observation. An explicit
confirmation asks the CMS to inspect that same repository/workflow/ref/code,
dispatch UUID and current source/target. Known terminal failure, or a successful
run with a stably missing/expired unique artifact, can authorize retirement.
A separate new UUID build request still needs another confirmation and one of
the original **three lifetime admissions**. Recovery does not dispatch, cancel,
rerun, delete, reset budget/history, deploy or complete an Astro job.

| Retained reason    | Required state                                                                                                                    |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------------- |
| `RUN_FAILED`       | Completed finite non-success first run, checked before and after.                                                                 |
| `ARTIFACT_MISSING` | Successful first run around two explicit successful empty artifact lists.                                                         |
| `ARTIFACT_EXPIRED` | Successful first run around two identical unique expired-artifact lists, correct identity/metadata and actual expiry in the past. |

Omitted/null/malformed lists, 404/timeouts, usable artifacts, live/unknown runs,
duplicates, changed metadata, changed source/Owner/target and reruns fail closed.
If content proof exists, expiry must match its artifact ID. There is no provider
blob download in recovery. The previous dist proof and dispatch outcome remain
historical. Current status endpoints show evidence for the newest candidate only.

## Workflow assumptions

Keep `.github/workflows/cms-build-preview.yml` read-only and dist-only, with the
existing dispatch/code guard, exact candidate title/UUID and one-day
`astro-cms-candidate-<UUID>` artifact contract. The normal Pages workflow remains
separate. Missing output does not prove absence of arbitrary side effects; the
reviewed workflow policy is part of the trusted boundary. Recovery rejects an
externally rerun `run_attempt > 1`; it cannot prevent a later administrator rerun.
GitHub inspection and SQL cannot be one distributed transaction.

## Private/public separation

`GET /api/v1/admin/releases/{id}/astro-candidate/recovery` reads retained private
evidence only. Owner/CSRF-protected `POST` accepts only an attempt UUID, performs
bounded fixed GETs and atomically retains `astro-candidate-recovery-v1` evidence
with audit. It accepts no URL, caller reason or workflow choice. No public route,
CMS secret or provider capability is added to the homepage. The private operator
must stop the old API, apply migration 25 and restart API/Flutter before use.

The existing token's repository Actions-write/Contents-read grant is broader than
this GET-only operation; local allowlists do not narrow the provider grant. The
homepage never receives this token. Trust in reviewed code, GitHub and TLS remains.
Two immediate empty reads are bounded evidence, not a long-term absence promise.

See [security/verification](security/cms-candidate-recovery-verification.md) and
[editable diagram](diagrams/cms-candidate-recovery.drawio). Prior
[content proof](cms-artifact-verification.md) remains required for a new usable
candidate. Controlled exact-byte deployment, independent deployed-output checks
and an atomic CMS receipt are still future gates. No hosted run or deployment is
performed by this documentation checkpoint.
