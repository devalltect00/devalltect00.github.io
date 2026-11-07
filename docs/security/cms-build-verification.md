# Astro CMS build security verification

Date: **7 November 2026 (Asia/Jakarta)**. Local development verification, not a
penetration-test certification or production security guarantee. Homepage baseline
`514a8b4`; CMS baseline `30e23b4`. The portfolio development copy and production
repositories are outside this change. No private account/record is used as a fixture.

Paths in the controls table refer to the Astro repository unless explicitly labeled CMS.

## Controls and evidence map

| Threat / invariant                                | Implemented boundary                                                                                              | Verification                                                                            |
| ------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| Draft/credential exposure                         | Anonymous fixed public route; no CMS env/auth/cookie loading                                                      | HTTP tests; public-source request assertions; signed fixture uses ephemeral test keys   |
| Redirect/content-URL fetch abuse                  | Trusted operator root origin, HTTPS, no redirects/proxy/asset fetch; loopback HTTP explicitly enabled             | `cms-http.test.ts`: origin, credentials/path/query, redirects and MIME                  |
| Oversize/stalled/truncated input                  | 2 MiB + 8192 byte package limit; 10s whole request deadline; bounded allocation and UTF-8 decode                  | HTTP oversize/partial/stall tests and snapshot bounds                                   |
| Ambiguous JSON/prototype pollution                | Reject duplicate and prototype-related keys, depth >64, malformed UTF-8                                           | `src/core/cms/strict-json.ts`; adversarial snapshot tests                               |
| Tampered/stale or mixed prepared identity         | Go-compatible complete content SHA-256, decimal int64 sequence, raw package hash and projection hash              | Snapshot/cache tests; real Go frozen normalization for all proof stages                 |
| Path traversal/known symlink cache input          | Canonical UUID selector; fixed filenames; lstat ancestors/members                                                 | Cache invalid-pin, junction/symlink, changed receipt and corrupted byte tests           |
| Partial preparation/retained input loss           | Exclusive task-owned temporary directory, complete rename; existing release revalidation without overwrite        | Cache failure, idempotency and immutable-conflict tests                                 |
| XSS/unwanted disclosure in page                   | Strict project/Profile models, selected plain-text props, explicit HTTPS URLs, no raw HTML or CV/media projection | Literal `<gateway>` browser text, selected-only/removal/optional-link output assertions |
| Disabled/unselected data or old personal fallback | Enabled Astro placement only; missing CMS fields stay absent                                                      | Initial/update/empty compiled output and authoritative empty tests                      |
| Runtime API dependency                            | No CMS browser code; sources stopped before builds and browser reads                                              | Both compiled proofs; static-only local server; no cache/package file in dist           |
| False consumer completion                         | Dist marker is a public build observation only                                                                    | No production CMS mutation code; no receipt endpoint/SQL change                         |

The two rendered kinds are strict; other known approved record bodies are checked
as part of the complete content hash but are not semantic UI models here. Integrity
hashes detect change; they do not authenticate approval. Source HTTPS/configuration
and the existing publisher trust boundary supply provenance. The operator/workspace
is trusted. This is not a DNS-rebinding/IP egress firewall or a sandbox against a
malicious local process racing filesystem swaps. Input labels/IDs are not executed.

The ignored local package contains only **already public** records, including kinds
not rendered here. It is not copied to `dist`; CI must not upload local caches as
public artifacts. Public marker identity is intentionally readable. An operator can
preview an older valid pin; a monotonic production deployment/receipt gate is future
work. Build-input failure preserves local dist, but later unrelated Astro failures
are not claimed to be filesystem-atomic deployments.

## Dependency remediation

The initial locked homepage toolchain audit found **24 advisories** (3 low,
14 moderate, 7 high). Updating compatible Astro/Vitest ranges and targeted patched
transitives removed the detected advisories. Final resolved direct versions are
**Astro 7.4.7**, **Vitest 5.1.0**, **Zod 4.6.6**, **Playwright 1.63.0** and
**Axe Playwright 4.13.0**. App metadata remains version **1.0.0**; no tag is created.
The frozen lockfile is authoritative; compatible ranges do not imply a major upgrade.

`pnpm-workspace.yaml` documents targeted overrides for affected transitive ranges:
fast-uri 3.1.9, minimatch 11.0.1, sharp 0.35.5, smol-toml 1.9.0,
source-map-js 1.2.2, undici 8.12.2 and Vite 8.3.3. Keep the overrides until their
parents resolve safe versions naturally; rerun lint/build/tests before removing
them, especially the minimatch compatibility override. Source advisory references:
[Vite](https://github.com/advisories/GHSA-587h-7wmj-h5pm),
[Vitest](https://github.com/advisories/GHSA-qxg6-89mg-rqv8),
[Undici](https://github.com/advisories/GHSA-2m5v-pch6-74hx),
[sharp](https://github.com/advisories/GHSA-wq5f-xc86-pv6w),
[minimatch](https://github.com/advisories/GHSA-g3hv-75wg-m249),
[smol-toml](https://github.com/advisories/GHSA-r4xh-jqrq-34v2),
[source-map-js](https://github.com/advisories/GHSA-68fv-2mgg-jv7q),
[fast-uri](https://github.com/advisories/GHSA-w74x-v6f9-j3xg).

The complete dependency audit reports **no known vulnerabilities** after remediation.
This is a timestamped advisory result, not proof of zero vulnerabilities. The new
Playwright/Axe dependencies support real static verification; Zod enforces public
DTOs. No third-party runtime service or analytics is added.

## Measured checks and local evidence

Final measurements are recorded below and in CMS checkpoint progress.
Commands: `pnpm install --frozen-lockfile`, `pnpm validate`, `pnpm audit`, standalone
compiled proof and production signed Worker compiled proof (both with optional Go
normalization). CI runs the public standalone proof with Chromium/read-only permissions.
Local verification uses Edge, Node 24.21 and pnpm 12.5.1. E-drive evidence is retained
outside Git under `E:/portfolio-cms-y35sunj8/.portfolio-cache/astro-independent` and
`astro-signed`, with stage logs, result JSON and EN/ID desktop/mobile light/dark PNGs.

Early checks caught invalid synthetic disabled-channel normalization, an accessibility
runner requiring an explicit browser context, and insufficient existing sticky-nav
light-theme contrast plus an unnamed mobile footer link. The fixture, runner and navigation color were corrected;
accessibility rules were not disabled. A temporary audit JSON accidentally included
in the build workspace's formatting scan was moved into the ignored evidence cache.

CMS changes are test-only fixtures plus an optional Go test. No full Flutter suite,
production database migration or new full backend coverage percentage is claimed for
this boundary. Prior required Go-module advisory and Docker/Linux race availability
limits remain as recorded in the M5h report; this checkpoint does not clear them.
No production deploy, live provider test or authenticated Astro delivery completion
has occurred. Continue periodic audits, review and separately authorized deployment.

## Final checkpoint measurements

- Frozen dependency install and full `pnpm validate` passed: formatting, zero-warning
  ESLint, Astro/TypeScript **42 files / 0 errors / 0 warnings / 0 hints**, **17 unit
  tests / 4 files**, **9 repository checks**, authored production static build.
- Both compiled proof sources passed: **6 real Go package normalizations**, **6
  pinned static builds** (initial/update/empty for each), **8 bilingual responsive
  browser journeys / 16 light-dark Axe checks**. Sources were stopped before builds
  and browsing; removals, optional fields, plain text and invalid-pin preservation
  remained required. No destination receipt was recorded.
- CMS publishing package tests and vet, Worker TypeScript, **Worker regression
  suite** and original six-kind portfolio fixture smoke passed. No CMS production
  application logic or database schema changed.
- Full JavaScript dependency audits and redacted Git-visible source/retained-history
  scans passed for both changed repositories. Existing source/env token and
  private-artifact binary/staged canaries passed. Two separate untagged message
  files passed installed Custy validation. The previous Go required-module advisory
  is not reclassified as fixed by this frontend/tooling checkpoint.
- Mirrored three-page native draw.io XML and PNG previews were inspected. Responsive
  screenshots were inspected; automated checks cover both languages and themes.
  Git indexes/heads, application versions, authored data and the portfolio's
  pre-existing next-env.d.ts edit were preserved. No commit/tag/deployment performed.

The patched install initially selected an incompatible `@emnapi/runtime` 2.0 beta
peer beside the checker/compiler's stable 1.x core. The manifest now narrowly
overrides that beta resolution to **1.11.3**, within the parent's accepted stable
peer range. Real type checks, frozen installation, builds and audit were rerun;
the checker was not bypassed and no prerelease runtime is required for this fix.
