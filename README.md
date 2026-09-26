# Devalltect Homepage

> The front door to the Devalltect project ecosystem.

[![Astro](https://img.shields.io/badge/Astro-7-BC52EE?logo=astro&logoColor=white)](https://astro.build/)
[![TypeScript](https://img.shields.io/badge/TypeScript-strict-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Node.js](https://img.shields.io/badge/Node.js-24-339933?logo=nodedotjs&logoColor=white)](https://nodejs.org/)
[![GitHub Pages](https://img.shields.io/badge/GitHub%20Pages-ready-222?logo=github)](https://pages.github.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-3FCAFF.svg)](LICENSE)

Devalltect Homepage is a bilingual Astro gateway that introduces the Devalltect
identity and directs visitors to the interactive portfolio, documentation
portal, source repositories, and maintained developer tools.

It deliberately stays focused: the homepage helps visitors choose a path, the
portfolio tells the broader visual story, and the documentation portal provides
the technical detail.

## Project metadata

| Property        | Value                             |
| --------------- | --------------------------------- |
| Release         | `v1.0.0`                          |
| Runtime         | Node.js 24                        |
| Package manager | pnpm 12.5.1                       |
| Locales         | English and Bahasa Indonesia      |
| Deployment      | Root GitHub Pages through Actions |

## Highlights

- English at `/` and Indonesian at `/id/`
- Responsive light and dark themes
- Typed, centralized project and translation data
- Astro-native components with minimal client-side JavaScript
- Accessible navigation and reduced-motion support
- Static GitHub Pages deployment from the repository root
- Repository metadata synchronization for GitHub and GitLab

## Requirements

- Node.js 24
- pnpm 12.5.1 through Corepack

## Development

```bash
corepack enable
pnpm install --frozen-lockfile
pnpm dev
```

Open `http://localhost:4321/`. For a network preview, pass an explicit host:

```bash
pnpm dev --host 0.0.0.0
```

### Optional destination overrides

Copy `.env.example` to `.env` when a local build should point somewhere other
than the production destinations:

```dotenv
PUBLIC_PORTFOLIO_URL=https://devalltect-portfolio.vercel.app/
PUBLIC_DOCUMENTATION_URL=https://devalltect00.github.io/devalltect-docs/
```

Both variables are optional. Blank, malformed, and non-HTTP values fall back to
the URLs shown above. For GitHub Actions, create repository variables with the
same names; the CI and Pages workflows pass them to Astro at build time.

## Validation

```bash
pnpm validate
```

This checks formatting, lint rules, Astro and TypeScript diagnostics, unit and
repository-maintenance tests, and the production build.

## Production build

```bash
pnpm build
pnpm preview
```

Generated files are written to `dist/`. GitHub Pages publishes that directory
through `.github/workflows/pages.yml`; select **GitHub Actions** as the Pages
source in repository settings.

## Content maintenance

- Project cards: `src/data/projects.ts`
- Shared destinations and profile links: `src/data/site.ts`
- English copy: `src/i18n/en.ts`
- Indonesian copy: `src/i18n/id.ts`
- Design tokens and responsive behavior: `src/styles/global.css`

Review metadata updates safely before applying them:

```bash
pnpm repository:metadata:dry-run --provider github
pnpm repository:metadata:sync --provider github
```

## Related destinations

- [Interactive portfolio](https://devalltect-portfolio.vercel.app/)
- [Documentation portal](https://devalltect00.github.io/devalltect-docs/)
- [GitHub profile](https://github.com/devalltect00)
- [GitLab profile](https://gitlab.com/devalltect00)

## Documentation

- [Development journey](docs/DEVELOPMENT_JOURNEY.md)
- [Badge references](docs/badges.md)
- [Public documentation](https://devalltect00.github.io/devalltect-docs/docs/devalltect-homepage)
- [Security policy](SECURITY.md)

## License

Released under the [MIT License](LICENSE).

_Pause · Think · Act_
