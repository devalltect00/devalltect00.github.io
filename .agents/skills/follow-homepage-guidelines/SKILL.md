---
name: follow-homepage-guidelines
description: Preserve the architecture, bilingual content, design system, accessibility, validation, and root GitHub Pages deployment contract of the Devalltect Astro homepage.
---

# Follow Homepage Guidelines

Use this skill for implementation and review work in the Devalltect homepage.

1. Read the root `AGENTS.md` before changing files.
2. Keep page routes thin and reusable UI in `src/components/`.
3. Store project records in `src/data/` and translated copy in `src/i18n/`.
4. Update English and Indonesian content together.
5. Preserve static output, the root Pages URL, accessible focus, reduced motion,
   responsive behavior, and light/dark theme contrast.
6. Prefer Astro and small native browser scripts over a client framework.
7. Run focused tests, then `pnpm validate` before completion.
8. Never perform Git or provider mutations without explicit approval.
