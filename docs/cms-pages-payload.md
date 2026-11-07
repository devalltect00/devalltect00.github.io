# Exact Pages payload handoff (M5o)

Follows homepage **685a49e** and CMS **6c650e8**. Private Go tooling now converts
an exact digest-pinned, frozen-content-verified candidate ZIP to a deterministic
gzip/USTAR payload. This public repository independently verifies every payload
file and extends its optional private compiled proof. There is no new deployment
workflow, runtime setting, dependency, public UI or app version.

## Local independent check

```powershell
node scripts/content/verify-pages-payload.mjs candidate.tar.gz pages-evidence.json
```

The manifest must come from a **trusted Go candidate check**, not from an
untrusted source supplying a self-consistent archive and hashes. The private CLI
is `cms astro-build pages --package ... --request ... --zip ... --zip-sha256 ...
--out ...`; it loads no CMS settings and exclusively creates a new archive.
The expected source ZIP digest must originate from the retained trusted artifact
proof. Go checks complete ZIP integrity and independently derives actual EN/ID
visitor facts from the frozen approved package before preparing a payload.

The public checker reads only bounded regular local files. It pins the exact
gzip hash/length, validates a single canonical member with bounded raw inflate,
CRC/ISIZE and no trailing/concatenated member, then checks complete USTAR headers,
checksums, sorted names, file bytes, padding and endings. It never extracts,
executes, serves or fetches members. Unknown fields, unsafe paths, links, PAX,
changed runtime assets, missing/extra files and resource overruns fail with a
constant diagnostic. Semantic interpretation remains Go's separate responsibility.

## Public contract and limits

Scope is `astro-pages-payload-v1`, schema version 1. `source` retains the whole
candidate ZIP hash, file count, expanded bytes and original frozen content/request
evidence. `payload_sha256` and `payload_bytes` bind the exact output container;
`entries` records every ordered `{name,size,sha256}` regular file, including binary
assets and empty files. Empty directories are omitted. Metadata becomes mode
0644, UID/GID/time 0, empty owners/links, ordinary USTAR and portable gzip OS 255.
Unsupported USTAR lengths fail. File bytes are exact; the archive wrapper changes.
Determinism is scoped to the pinned Go compressor, not promised across versions.

Source bounds remain 8 MiB ZIP, 32 MiB expanded and 256 entries. USTAR is capped at
32 MiB + 256 KiB + 1024 bytes; gzip at that bound + 1 MiB; manifest at 1 MiB.
The Go CLI uses a 15-second context checked between reads/files; compression is
bounded CPU work rather than a hard real-time interruption guarantee. Trusted
local operator paths are not a sandbox against concurrent host-user mutation.

The public manifest is unsigned integrity data and **cannot authorize deployment
or complete a CMS job**. The future coordinator must bind it to retained current
Owner/source/target/attempt/proof state before external actions. GitHub documents
the Pages gzip/TAR requirement and link restriction in its
[official workflow guidance](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages#configuring-the-upload-pages-artifact-action).
Hosted artifact wrapping/upload/acceptance is not tested by this checkpoint.

## Real compiled proof

Standalone public CI needs no private repository or credential. The existing
optional private command adds Go preparation and this independent Node inventory
check after all initial/update/empty compiled candidates:

```powershell
node scripts/testing/prove-cms-build.mjs --go D:/project/develop/devalltect-cms/backend
node scripts/testing/prove-cms-build.mjs --worker D:/project/develop/devalltect-cms/edge/scripts/portfolio-fixture.mjs --go D:/project/develop/devalltect-cms/backend
```

Provider metadata remains synthetic. Sources stop before build/browser checks.
`CMS_ASTRO_PAGES_OUTPUT` is an explicit test-process output directory, not a
public runtime or production credential setting. The fixture requires a new
directory and writes public payload/manifest there only when explicitly supplied.
No token or live provider enters this boundary. Existing Pages/candidate workflows
and authored records stay unchanged. No migration beyond 25 is needed.

See [verification](security/cms-pages-payload-verification.md) and
[editable diagram](diagrams/cms-pages-payload.drawio). Controlled exact-byte
upload/deployment, independent deployed-output proof and an atomic CMS Astro
receipt remain future work. A valid local payload is not hosted publication.

Keep generated ZIP/payload/evidence files in the existing ignored CMS `.tmp/` or
Astro `.cms-cache/` directory. They are local outputs, not repository source.

## Subsequent staging boundary

M5p now implements [exact optional artifact staging](cms-pages-staging.md). Earlier upload TODOs describe the historical M5o checkpoint. Owner coordination and protected deployment/deployed-output completion remain separate.
