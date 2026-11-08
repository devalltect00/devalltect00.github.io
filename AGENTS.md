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

## CMS build boundary

Read docs/cms-builds.md before changing CMS mode. src/core/cms owns public wire
validation; src/server/cms owns Node-only preparation/integrity; src/data/homepage.ts
owns one explicit process-environment build pin. Keep the cache out of dist/client
code. Empty CMS content remains empty; only an absent selector enables authored
mode. No public marker authorizes approval, deployment or job completion. Run
pnpm test:cms after changes to projection/build output in addition to pnpm validate.

Exact automation must pass HOMEPAGE_CMS_REQUEST along with the immutable release
pin. Read docs/cms-build-handoff.md before changing semantic data attributes or
candidate permissions. Do not replace visible fields with hidden proof payloads.
The separate candidate workflow cannot deploy or finish a CMS job; keep its
repository authority/input policy tests in scripts/repository/tests.

M5n documents private bounded recovery in docs/cms-candidate-recovery.md. Do not turn missing output or a public marker into deployment/completion authority. Keep candidate workflow read-only and dist-only; no budget reset, rerun or cancellation is implied.

## Exact Pages payload boundary

M5o adds the local verifier in scripts/content/pages-payload.mjs; read docs/cms-pages-payload.md. No extraction or provider access is allowed there. Keep complete asset inventory and strict gzip/USTAR bounds, and never promote an unsigned public manifest into deployment authority.

## Exact Pages staging boundary

M5p adds the optional standard-library Go stager in scripts/pages-tool and cms-pages-stage.yml. Read docs/cms-pages-staging.md. Keep the uploader disabled by default, main/first-attempt/code/environment gated, with no Pages/OIDC-write permission. Run standalone Go tests/vet plus public validation and joint private compiled proofs when either writer/contract changes.

M5q documents the private Owner coordinator in docs/cms-owner-staging.md. Public request JSON remains unsigned; direct operator dispatch does not gain private SQL budget/current-Owner guarantees. Keep hosted uploader off by default and acceptance separate from artifact/deployment proof.

M5s documents private successful-upload reconciliation in docs/cms-stage-reconciliation.md. Keep public workflows/tools unchanged and disabled by default. Historical uploaded-byte evidence does not authorize publication, retries or deployment.

M5t documents conservative private retirement in docs/cms-stage-recovery.md. Public workflows/tools remain unchanged and disabled by default. Metadata retirement never cancels/deletes/deploys, and unknown dispatches do not gain retry authority.

## Deployment receiver integrity prerequisite

Read docs/cms-deployment-receiver.md before changing the independent local intent
or wrapper checks. No CLI/runtime/workflow exposes them. Public unsigned intent
is not current CMS permission; do not attach Pages writes before authenticated
fresh receiver authority and ordering are designed. Preserve complete canonical
byte checks and rerun joint compiled initial/update/empty receiver proof.
