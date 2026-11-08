# Uploaded staging integrity evidence

Reviewed 8 November 2026. This is a local prerequisite, not authenticated hosted
reconciliation or deployment. Source: the private CMS `astrobuild/staged.go`.
`TestStagedArchive` compares real canonical payloads in Store/Deflate wrappers.
`TestStagedArchiveRejects` checks 22 hostile/invalid cases with independently pinned
outer bytes, including traversal, duplicate/extra members, links/directories,
CRC mutation, ZIP64/multi-disk/comment/trailing wrapper, changed/truncated/trailing
TAR, wrong candidate/payload/request/pin, size bound and cancellation.

Targeted native Go tests cover astrobuild, astrogithub and astrocandidate. The
astrobuild package reports **91.5% statement coverage**; this is package coverage,
not a newly measured whole-backend result. Go vet covers the same three packages.
No database/API/UI/schema behavior changed; the previous complete SQL/Flutter
results remain historical, not rerun claims. The fixture's initial directory
writer error was corrected to create an empty directory before the passing run.

The verifier bounds the outer directory before ZIP parsing, reconstructs expected
bytes from the frozen candidate, checks the prepared hash, then checks exact
member type/name/size/CRC and all canonical TAR bytes without extraction. Digest
evidence alone does not authenticate run/source/current authority. No network,
SQL write, retry, retirement, publication or deployment is performed by this code.
Hosted wrapping acceptance and production load/race/device execution remain
untested. The next caller must acquire exact authenticated staging run/artifact
bytes and persist evidence under current Owner/source/target/proof fences.

The native diagram and PNG are mirrored and visually reviewed. Both checkpoints
remain untagged; no commits, tags, hosted actions or editing migrations performed.
Public runtime/workflows/dependencies remain unchanged. Schema stays **26**.

Public source-matched `pnpm validate` passed: clean formatting/lint, 52 diagnostic files with no issues, 25 unit tests, 23 repository checks and authored static build. Both separate Custy messages validate. The CMS retained-history scan passed. Final Git-visible UTF-8 source scans passed for both repositories; ignored live credentials and binary artifacts are excluded. The initial public scanner invocation referenced a nonexistent custom config and was corrected to use its default policy. No JavaScript or Go dependencies changed; whole-backend/database/Flutter suites were not rerun for this unwired pure verifier.
