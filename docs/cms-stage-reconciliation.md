# Owner staging upload reconciliation

M5s, 8 November 2026. CMS baseline `95485aa`; homepage baseline `c9e02cb`.
This public checkpoint updates documentation and mirrored diagrams only.

The private Owner API/UI now connects the fixed staging-run reader and uploaded
TAR verifier. For the current acknowledged **accepted** request, it checks current
SQL authority/source/proof/target, reacquires verified candidate bytes, selects
one exact uploaded staging artifact, compares all TAR bytes against the prepared
payload, repeats run/artifact checks and retains immutable proof/audit under a
final current-authority transaction. Provider work never holds SQL locks.

The existing source candidate cap stays 8 MiB; staged wrapper acquisition has its
own finite MaxPagesTar + 1 MiB cap. Metadata/API credentials are isolated from
signed blob downloads. No archive is extracted. Cached proof/status reads use
SQL only; upload verification never dispatches, retires, deploys or completes a
job. A retained proof is historical, not current availability/approval authority.

Private migration **00027** stores immutable `astro-stage-upload-v1` evidence.
This public repository needs no migration/new secret/runtime change. Original
candidate, staging and Pages workflows/tools, dependency locks and versions are
unchanged. Hosted staging remains off by default; no hosted setting/action was
performed. Actual environment reviewers remain operator configuration.

Uncertain acknowledgement, live/failed/skipped runs or missing/expired artifacts
remain blocked from invented success or automatic retries. Next: authenticated
conservative failure/uncertain recovery, protected exact deployment, deployed
output proof and atomic destination completion. SQL/provider revocation cannot
be atomic; future deployment must repeat authority/content checks. Trusted
admin/runner/provider compromise and production load remain outside local proof.

See [native flow](diagrams/cms-stage-reconciliation.drawio) and
[verification evidence](security/cms-stage-reconciliation.md). The optional
compiled proof now exercises genuine compiled initial/update/empty candidate
bytes through the complete private two-archive acquisition protocol using
synthetic provider authority. It does not prove hosted upload acceptance.
