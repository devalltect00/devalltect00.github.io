# Development Journey

This file records the architectural decisions and repeatable setup steps for
the Devalltect homepage. It is an internal development record, not the public
portfolio narrative.

## 2026-09-21 — Initial Astro foundation

### Objective

Create a distinct root homepage for `devalltect00.github.io` that introduces
the wider Devalltect ecosystem without duplicating the portfolio or
documentation portal.

### Decisions

- Selected Astro static output because the site is primarily content and
  navigation with a small interaction budget.
- Selected strict TypeScript and pnpm on Node.js 24.
- Kept English at `/` and Indonesian at `/id/`.
- Used Astro components instead of a UI framework to minimize shipped
  JavaScript.
- Reused the Devalltect palette and display/body font pairing while creating a
  calmer visual identity than the interactive portfolio.
- Centralized project records and translated content outside page components.
- Added accessible theme, navigation, focus, and reduced-motion behavior.
- Configured root GitHub Pages deployment through GitHub Actions.

### Setup commands

```bash
pnpm create astro@latest devalltect00.github.io --template minimal
pnpm add @astrojs/sitemap @fontsource-variable/dm-sans @fontsource-variable/space-grotesk
pnpm add -D @astrojs/check prettier prettier-plugin-astro eslint eslint-plugin-astro typescript vitest jsdom @testing-library/dom
```

### Validation contract

```bash
pnpm validate
```

### Next visual assets

- Add a social preview image before the first public release.
- Add project imagery only when it improves navigation instead of repeating the
  portfolio gallery.
- Validate the production Pages URL on mobile and desktop after deployment.

## 2026-09-21 — Configurable public destinations

- Added optional `PUBLIC_PORTFOLIO_URL` and `PUBLIC_DOCUMENTATION_URL` build
  variables.
- Preserved the production portfolio and documentation URLs as defaults.
- Added URL validation so blank, malformed, and non-HTTP overrides cannot
  produce unsafe or broken navigation.
- Exposed optional GitHub repository variables to validation and Pages builds.

## 2026-09-24 — First stable release preparation

- Promoted the planned release baseline from the temporary `0.1.0` development
  version to the first stable Semantic Versioning release, `1.0.0`.
- Replaced the untagged development-checkpoint template with distinct Custy
  commit and annotated-tag messages for v1.0.0.
- Added repository-scoped VS Code formatting, extension, launch, and Windows
  Command Prompt defaults without changing other operating systems.
- Added bilingual public documentation to Devalltect Docs for setup,
  configuration, architecture, deployment, and release status.

No commit, tag, release, deployment, push, or remote metadata mutation was
performed as part of this preparation.
