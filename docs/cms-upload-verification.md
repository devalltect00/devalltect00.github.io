# Exact uploaded staging archive verification

M5r prerequisite, 8 November 2026. Baselines: CMS `af42fdc`, Astro `649953f`.

The CMS now has a local `astrobuild.VerifyStagedArchive` primitive for the ZIP
wrapper returned when downloading a staging artifact. It accepts the trusted
server preparation, original verified candidate, its frozen projection, uploaded
ZIP bytes and an independently selected uploaded-ZIP digest. It reconstructs the
canonical gzip/TAR from the candidate and checks the prepared payload hash before
comparing **every TAR byte** with the uploaded `artifact.tar`. Outer ZIP order or
compression is not content authority. Both ordinary Store and Deflate wrappers
are supported; directories, links, extra/duplicate/nested/traversal entries,
changed/truncated/trailing TAR bytes and corrupted ZIP CRCs are rejected.

Before allocating ZIP directory entries, the verifier checks a bounded ordinary
single-entry EOCD directory: no multi-disk/ZIP64/comment/trailing wrapper, central
directory at most 4096 bytes and outer ZIP at most `MaxPagesTar + 1 MiB`. Only
`artifact.tar` is allowed; declared expanded size must exactly match the bounded
canonical TAR. Reads check cancellation and ZIP integrity. No file is extracted,
executed or served; no network, credential or SQL operation occurs. Errors use
the existing constant `ErrInvalid`; evidence contains only whole ZIP/payload/TAR
hashes and TAR byte count. Handwritten code/types/functions explain invariants.

This is a **local integrity primitive**, not the authenticated reconciliation
endpoint. Its caller must authenticate the exact first-run staging workflow,
repository/code/artifact identity, current Owner/source/candidate/target and
trusted prepared hashes. A supplied digest or self-consistent public tuple is
not that authority. A successful byte comparison never writes an observation,
retires/retries a request, changes a job, deploys, or completes publication.

No migration, configuration, API, Flutter, public runtime/tool/workflow,
dependency or app version changes. The editing database was not read/migrated.
Hosted staging remains disabled. Next: authenticated staging observation and
artifact acquisition with final SQL fences and append-only evidence, then
conservative recovery, protected deployment and deployed-content proof.

See [native flow](diagrams/cms-upload-verification.drawio). The matching
public diagram documents the same boundary. Exact hosted wrapper acceptance,
resource exhaustion under production load, Linux race and Android execution are
not claimed. Parsing limits mitigate resource abuse; they are not hard real-time
CPU guarantees. Trusted repository/runner/provider compromise remains a limit.
