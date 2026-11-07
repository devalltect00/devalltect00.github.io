# Candidate correlation verification

Reviewed 7 November 2026. Homepage baseline 98f67a6; CMS baseline 38e29d4.
The workflow and request creator bind public dispatch UUID and exact code SHA.
The guard runs before checkout, and the pure TypeScript check runs before a
request file is written. Matching public IDs do not authenticate an Owner.

Workflow contents-read permissions, exact pinned checkout, disabled credential
persistence, pinned actions and dist-only one-day artifacts are policy-tested.
The private dispatch token's repository Actions-write permission is broader;
it must be protected independently. No real token or provider run was used.

Local validation passed **25 unit tests**, **12 repository checks**, strict Astro
diagnostics, lint, formatting and authored build. Additional compiled source/Go
and browser proof, audits and secret scans are recorded below.
The CMS suite passed **188 Flutter tests** and full real PostgreSQL integration
tests, with measured Go **85.94%** and Dart **94.39%** coverage. The external E:
workspace contains disposable evidence, not credentials or production authority.

No hosted workflow, Pages deployment, editing database migration, new Android
device run or Linux race run is claimed. The CMS requires migration 00022 and API/
Flutter restart before the optional screen can be used. Source role checks cannot
atomically revoke an already sent provider request; unresolved attempts block a
fresh dispatch rather than risk duplicate work. Independent run/artifact and
deployed-output verification plus a durable Astro receipt remain future work.

See [operations](../cms-owner-candidates.md) and the matching private CMS report
`docs/security/astro-owner-candidates-verification.md` for threat/control mappings.

## Final executed results

Final `pnpm validate` passed: **48 diagnostic files, zero errors/warnings/hints**,
**25 unit tests / 6 files**, **12 repository checks**, formatting, lint and authored
production build. The independent public-source and production-signed-worker
fixture proofs both passed with sources stopped before builds/browser reads.
Combined: **6 Go package checks, 6 Go requests, 6 pinned builds, 6 independent Go
HTML verifications, 8 bilingual responsive journeys and 16 themed Axe checks**.
Initial/update/empty content, optional links and safe invalid-pin failure remain
covered. Logs/screenshots/result.json remain outside source in the E: workspace.

Complete JavaScript audits report zero known vulnerability counts for both repos.
Redacted source/history scans passed for both repositories; curated sources cover
148 homepage and 948 CMS Git-visible UTF-8 files. Ignored private settings/caches
are excluded. Synthetic secret and binary private-artifact canaries passed.
CMS SQL generation reproduced exactly, OpenAPI/Compose validation passed and
the Flutter release web build passed. Its compiler emitted a Cupertino font-family
warning; this is not a new Android device/build claim.

The mirrored three-page draw.io diagrams and rendered previews were inspected.
Both untagged Custy messages passed validation. Dependency locks/application
versions, authored content, existing Pages workflow, Git heads/indexes and the
unrelated portfolio edit remain unchanged. No real provider action, token setup,
editing migration, commit, tag or deployment was performed. The next gate is
independent run observation/reconciliation before artifact/deployment completion.
