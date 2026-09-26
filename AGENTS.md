# AGENTS.md

## Project purpose

This repository is the Astro-based homepage published at
`https://devalltect00.github.io/`. It is a concise bilingual gateway to the
Devalltect portfolio, documentation portal, projects, GitHub, and GitLab.

Do not turn this site into a duplicate of the interactive portfolio or the
Docusaurus documentation portal. Keep navigation direct, content concise, and
the visual system consistent with the wider Devalltect ecosystem.

## Technology

- Astro with static output
- Strict TypeScript
- pnpm and Node.js 24
- Scoped Astro styles plus shared CSS design tokens
- Vitest and Node's test runner
- GitHub Actions and GitHub Pages

## Architecture

- `src/pages/` owns routes only.
- `src/layouts/` owns document structure and shared metadata.
- `src/components/` owns reusable presentation.
- `src/data/` owns destinations, project records, versions, and external URLs.
- `src/i18n/` owns all user-facing English and Indonesian copy.
- `src/styles/` owns tokens, global layout, responsive behavior, and motion.
- `scripts/repository/` owns repository-maintenance helpers and tests.

Do not hardcode project content inside page routes. Do not introduce React or
another client framework unless an interaction cannot be expressed cleanly with
Astro and small browser scripts.

## Content and localization

English is served at `/`; Indonesian is served at `/id/`. Every public content
change must be considered for both languages. Keep project names, product
versions, repository URLs, and technology names centralized.

## Design and accessibility

Preserve the Devalltect palette, typography, light/dark themes, responsive
layouts, keyboard focus, semantic headings, and reduced-motion behavior.
Animations should clarify hierarchy or state rather than delay navigation.

## Validation

Run targeted checks while editing, then complete:

```text
pnpm validate
```

Do not commit generated `dist/`, `.astro/`, or local environment files.

## Deployment

The repository name must remain `devalltect00.github.io` for root GitHub Pages
deployment. `astro.config.ts` therefore uses the root site URL and no project
base path. Pages must use the **GitHub Actions** source.

## Release safeguards

Do not create commits, tags, releases, pushes, or provider metadata mutations
without explicit approval. Run metadata synchronization in dry-run mode first.
