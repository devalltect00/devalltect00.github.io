# Owner run metadata observation (M5l)

Follows homepage **5ef4fb9** and CMS **331153d**. This is a documentation-only
homepage checkpoint: workflow code, authored content, dependencies and Pages
deployment remain unchanged. The private CMS adds an explicit known-run metadata
check beside existing Owner candidate controls, using migration 00023.

Accepted runs use their retained numeric identity. After a lost dispatch response,
an Owner may supply a numeric run ID; the fixed GitHub reader verifies repository,
workflow, event, exact branch/code SHA, first run attempt, head repository,
correlation title, canonical URLs and bounded timestamp window. It follows no
returned URLs and fetches no archive or logs. Opening/refreshing the screen reads
retained SQL only. Explicit checks never dispatch, cancel or rerun workflows.

Run name `CMS Astro candidate <dispatch UUID>` remains the workflow correlation.
That title and a successful run are not proof of the exact content inputs or
artifact bytes: repository actors can copy titles. Independent Go verification
against the frozen expected release is still required before any deployment or
Astro destination completion. The existing public marker is not a credential.

Observations are append-only and audited, with source/current-Owner/target checks
before and after I/O and a prior-revision compare-and-set. Identical metadata
deduplicates; older/terminal-changing evidence is rejected. Same-second first
completion is supported; unresolved history never unlocks another dispatch.
There is a 64-distinct-observation budget per candidate. Original uncertain
submission outcomes are retained even after metadata association is recovered.

No new homepage environment variables, workflow inputs, permissions or packages
are introduced. The existing private API token still has broader repository
Actions-write scope, despite metadata checks being GET-only. No actual hosted
run, credential change or deployment was performed.

Next are bounded artifact retrieval/provenance, exact Go content proof,
deployed-output verification and an atomic Astro receipt. Publication remains
pending. See [security evidence](security/cms-run-observations-verification.md)
and [three-page diagram](diagrams/cms-run-observations.drawio).

## Follow-up: candidate ZIP proof (M5m)

[Archive verification](cms-artifact-verification.md) now documents the private Owner content proof and optional compiled ZIP test hook. Metadata-only endpoints above remain separate. Workflow/deployment/completion authority is unchanged.
