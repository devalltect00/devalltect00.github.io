# Candidate archive proof (M5m)

Follows homepage **c01a528** and CMS **34e535b**. The private Owner flow now checks
the exact successful candidate ZIP against frozen approved content. The homepage
changes only its local proof script and documentation/diagrams; the candidate and
Pages workflows, authored content and dependencies remain unchanged.

The existing workflow name, dispatch UUID/code guard and one-day
`astro-cms-candidate-<UUID>` artifact contract remain required. Private Go checks
first-run/repository/workflow/ref/code identity, unique artifact metadata,
expiration, compressed size/digest, then actual bilingual HTML/marker facts.
The ZIP is inspected in memory, never extracted or executed. Successful evidence
is append-only under current Owner/source/observation fences (migration 00024),
and does not complete an Astro job or authorize deployment.

The signed download capability is isolated from GitHub authentication. Only fixed
GitHub API paths receive the token. A separate client accepts documented HTTPS
artifact service hosts, pins public DNS answers, and rejects credential forwarding,
proxies and secondary redirects. Complete compressed ZIP is capped at 8 MiB,
expanded content at 32 MiB and entries at 256. Plain Store/Deflate, portable paths,
unique names and all regular-file CRCs are checked. Unsupported archives/hosts
fail closed. No new permission, environment variable or package is introduced.

## Local compiled proof

The existing optional `--go <private backend>` path now invokes
`TestCompiledAstroArtifact` after every initial/update/empty dist build. The test
packs the real dist and exercises the real Go metadata/ZIP/content protocol using
synthetic provider authority. Independent public-source and signed-source fixture
commands are retained. Their original content sources stop before builds/browser
reads; the artifact fixture is a separate local HTTP transport, not the CMS API.
No real provider download or hosted provenance is claimed. Standalone public CI
continues without a private Go repository or credentials.

Retained content proof is historical. Controlled deployment of exact verified
bytes, independent deployed-output checks, bounded expired/failed-run recovery
and an atomic Astro receipt remain future gates. The existing marker/title/hash
alone is not authentication or publication success. See [security evidence](security/cms-artifact-verification-report.md)
and [native diagram](diagrams/cms-artifact-verification.drawio).
