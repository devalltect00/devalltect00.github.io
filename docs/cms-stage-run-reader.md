# Fixed staging-run metadata reads

M5r provider boundary, 8 November 2026. Baselines CMS `035f60a`, Astro `4b20dad`.

The private `StageClient.ObserveStage` adapter now reads one exact staging run
using the existing fixed GitHub client. `StageRunWork` distinguishes uploader
run identity/admission times from the source candidate run. Before I/O it checks
both target incarnations, prepared tuple, matching release/code and distinct
positive staging run ID. No caller URL, repository or workflow is accepted.

The shared internal `observeRun` routine checks repository/head-repository IDs
and names, no pull requests, exact workflow path/ID, main/code, first attempt,
workflow_dispatch event and `CMS Astro staging <stage UUID>` title. Exact run
URLs, bounded creation/updated timestamps and finite state/conclusion are
required. The ten-second GET uses existing no-proxy/no-redirect, per-response
256 KiB and finite body-budget controls. Duplicate JSON, wrong MIME/status,
oversized bodies and transport errors fail with constant errors. No response
body, signed location or token is exposed. Candidate observations still use
their original separate workflow/title through the same implementation.

This returns **metadata only**. Display-title correlation is not input
attestation or proof of uploaded bytes. The previous local uploaded-archive
verifier remains a distinct boundary. No metadata result writes SQL, retires
or retries staging, deploys, or completes a destination. A future Owner service
must assemble work under current SQL authority, authenticate/acquire the exact
artifact, verify its bytes, and repeat source/proof/Owner/target checks before
retaining append-only evidence. No reconciliation route/UI is exposed yet.

`TestObserveRunMetadata` and `TestObserveStageMetadata` exercise both real
adapters with the shared hostile loopback protocol: swapped IDs/workflows,
reruns, forks, PRs, foreign refs/code, bad titles/URLs/timestamps, nonfinite
state, duplicate JSON and unsafe transport. Staging-only target/source-run/
missing-preparation/release/code mutations send zero GETs. Real HTTP remains
synthetic: no hosted credentials, runs or configuration were used.

No schema/API/Flutter/config/dependency/version change. Schema remains **26**;
no new migration is needed or applied. Hosted uploader stays disabled. Next:
authenticated artifact acquisition and uploaded-byte verification, durable
Owner observations and UI, then conservative recovery and exact deployment.
