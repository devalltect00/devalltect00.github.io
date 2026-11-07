# Approved CMS builds for the Astro gateway

This is a **local static build preview**, following homepage `514a8b4` and CMS
`30e23b4`. It does not deploy Pages, trigger a provider or finish an Astro CMS job.
The portfolio consumer remains separate and unchanged in this checkpoint.

## What visitors see

The existing concise EN/ID gateway, navigation, themes and static root URLs remain.
CMS mode shows only projects with an explicitly enabled **Astro** placement, ordered
by that placement's position and then stable ID. Portfolio-only placement does not
select a project here. Titles and summaries use the requested locale; status labels
are translated. Technology labels remain distinct, including React and React JS.
Explicit approved repository, website, demo and documentation links become buttons.
Missing optional links, technology labels or version labels are simply omitted.

The approved Profile supplies the display name, public contacts and a concise
excerpt of current focus (or summary), at most 320 Unicode code points. Plain text
stays text. A removed Profile removes its personal contacts and name; an empty
project selection shows a translated empty state. CMS mode never restores old
authored personal data. The Devalltect brand and configured portfolio/docs gateway
destinations remain. CV links and project imagery are not rendered in this slice.
Skills, Journey, education and publications belong to the portfolio, not this gateway.

## Prepare and build (PowerShell)

Run from this repository with Node 24 and pnpm 12.5.1:

```powershell
pnpm install --frozen-lockfile
pnpm content:cms:prepare --origin https://your-approved-public-reader.example
# Copy the actual release UUID printed by preparation:
$env:HOMEPAGE_CMS_RELEASE = '<prepared-release-uuid>'
pnpm build
pnpm preview
```

The origin above is a placeholder, not a provisioned service. Preparation reads
exactly `/public/v1/content/en` anonymously. The package already contains both
languages; one prepared immutable release supplies `/` and `/id/` together.
For an independently running local public reader only:

```powershell
pnpm content:cms:prepare --origin http://127.0.0.1:8790 --allow-loopback
```

Use the public reader origin, **not** the private CMS API port 8090. Do not pass a
database URL, Owner session, API key, signing seed or export grant. Preparation
does not load CMS environment files. HTTPS is required except explicit literal
loopback/localhost development mode. Redirects, URL credentials and origin paths,
queries and fragments are rejected. The operator must choose a trusted reader.

`HOMEPAGE_CMS_RELEASE` is a server-side **process environment** selector. Set it in
the terminal or CI process running the build; placing it only in Vite's `.env`
file does not opt in. It is never a `PUBLIC_` variable. Preparation does not
silently change this selector. Empty, malformed or unavailable explicit pins fail.

To return to the existing authored source deliberately:

```powershell
Remove-Item Env:HOMEPAGE_CMS_RELEASE -ErrorAction SilentlyContinue
pnpm build
```

## Immutable local input and output

Preparation validates UTF-8, bounded bytes, strict JSON/envelope, the Go-compatible
content digest, project/Profile DTOs and cross-record references. It creates:

```text
.cms-cache/releases/<uuid>/
  package.json       # complete anonymous public package, never a private export
  complete.json      # raw-byte hash plus exact release/sequence/projection identity
dist/
  index.html
  id/index.html
  .well-known/homepage-content.json
```

A temporary task-owned directory is renamed only after both input files exist.
Retained releases are rechecked and never overwritten. Re-preparing the same
semantic identity is idempotent. Source outage or failed preparation leaves earlier
releases intact. Build validates the explicit pin **before** invoking Astro, so
an invalid pin does not erase the last successful local `dist`.

The source may be offline during the actual build and browsing. No browser CMS
request or database connection is added. Only the selected presentation props
reach Astro components. The cached full public package (including nonrendered
approved kinds), completion file and server adapters are not copied into `dist`.

The public marker and HTML metadata bind release ID, decimal sequence, complete
snapshot digest and selected projection digest. Scope is
`astro-project-profile-preview-v1`; consumed kinds are `project`, `profile`.
These are **build observations**, not signatures, approval authority, deployed
site evidence or a CMS success receipt. Authored builds identify themselves as
`authored`. An intentionally pinned older release can be previewed locally; this
command does not implement a monotonic deployment gate.

## Where to extend

| Responsibility                                              | Location                              |
| ----------------------------------------------------------- | ------------------------------------- |
| Strict public wire rules and duplicate-key/depth protection | `src/core/cms/`                       |
| Complete public integrity and selected gateway projection   | `src/server/cms/snapshot.ts`          |
| Trusted-operator anonymous fixed-route request              | `src/server/cms/http.ts`              |
| Local immutable preparation and build-time revalidation     | `src/server/cms/cache.ts`             |
| One centralized build selection, metadata and public marker | `src/data/homepage.ts`                |
| Preparation/build process entrypoints                       | `scripts/content/`                    |
| Pure/adversarial boundary tests                             | `tests/cms-*.test.ts`                 |
| Real static browser proof                                   | `scripts/testing/prove-cms-build.mjs` |

Routes remain thin. Keep Node/filesystem/network adapters out of browser scripts.
Do not add raw HTML, runtime CMS fetches or inferred personal links. New rendered
kinds require explicit schemas, projection and tests; the four other known approved
kinds are integrity-checked opaque records here, not validated UI models.

## Verify without live accounts

```powershell
pnpm validate
pnpm audit
pnpm exec playwright install chromium
pnpm test:cms
```

The fixture uses only public synthetic data. It prepares initial, updated and empty
releases, stops the source, builds each pin and verifies output identities/removals.
Updated output is exercised in EN/ID at desktop/mobile widths, light/dark themes,
with reduced motion, overflow checks and Axe WCAG A/AA checks. Invalid explicit
input must fail without replacing successful output. Local evidence is ignored in
`.cms-cache/proof/`. On Windows with installed Edge:

```powershell
$env:HOMEPAGE_BROWSER_CHANNEL = 'msedge'
pnpm test:cms
```

The read-only CI workflow installs Chromium and runs the standalone proof. Existing
Pages deployment remains authored until a separately reviewed dispatch/deployment
workflow explicitly prepares and pins a release. No credentials are required here.
See [verification and limits](security/cms-build-verification.md) and the native
[three-page diagram](diagrams/cms-builds.drawio).
