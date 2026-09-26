# Project Structure

# Repository Overview

This repository was analyzed using composable project and framework metadata.

Detected technologies: `nodejs`, `javascript`, `typescript`, `pnpm`

Recognized top-level directories:

- `.config/` — Project configuration files.
- `docs/` — Project documentation and technical references.
- `public/` — Static files copied or served without compilation.
- `scripts/` — Utility scripts for development or automation.
- `src/` — JavaScript or TypeScript application source code.
- `tests/` — Automated tests.

---

## Repository Structure

(project type: ProjectType.NODEJS)

```text
.
├── .agents
│   └── skills
│       └── follow-homepage-guidelines
│           └── SKILL.md
├── .astro/ ... (collapsed)
├── .config
│   ├── custy
│   │   ├── templates
│   │   │   ├── backups/ ... (collapsed)
│   │   │   ├── changelog
│   │   │   │   └── changelog.j2
│   │   │   ├── examples
│   │   │   │   ├── commit_message
│   │   │   │   └── tag_message
│   │   │   ├── commit-message.txt
│   │   │   └── tag-message.txt
│   │   └── config.toml
│   └── doc_gen
│       └── config.toml
├── docs
│   ├── badges.md
│   ├── DEVELOPMENT_JOURNEY.md
│   └── project_structure.md
├── logs/ ... (collapsed)
├── other_venv/ ... (collapsed)
├── public
│   ├── favicon.ico
│   ├── favicon.svg
│   ├── robots.txt
│   └── site.webmanifest
├── scripts
│   └── repository
│       ├── src
│       │   └── sync_metadata.mjs
│       └── tests
│           └── sync_metadata.test.mjs
├── src
│   ├── components
│   │   ├── common
│   │   │   ├── Icon.astro
│   │   │   └── SiteFooter.astro
│   │   ├── home
│   │   │   ├── Contact.astro
│   │   │   ├── Explore.astro
│   │   │   ├── Hero.astro
│   │   │   ├── LandingPage.astro
│   │   │   ├── Principles.astro
│   │   │   └── Projects.astro
│   │   ├── navigation
│   │   │   └── Navigation.astro
│   │   └── theme
│   │       └── ThemeToggle.astro
│   ├── data
│   │   ├── projects.ts
│   │   └── site.ts
│   ├── i18n
│   │   ├── en.ts
│   │   ├── id.ts
│   │   └── index.ts
│   ├── layouts
│   │   └── BaseLayout.astro
│   ├── pages
│   │   ├── id
│   │   │   └── index.astro
│   │   ├── 404.astro
│   │   └── index.astro
│   ├── styles
│   │   └── global.css
│   └── types
│       └── site.ts
├── tests
│   └── site-data.test.ts
├── .editorconfig
├── .gitignore
├── .pre-commit-config.yaml
├── .prettierignore
├── .prettierrc.json
├── AGENTS.md
├── astro.config.ts
├── CHANGELOG.md
├── CLAUDE.md
├── eslint.config.js
├── LICENSE
├── package.json
├── pnpm-workspace.yaml
├── README.md
├── SECURITY.md
├── tsconfig.json
└── vitest.config.ts
```

---

## Recognized Files

- `README.md` — Project overview and introduction.
- `CHANGELOG.md` — History of notable changes between releases.
- `LICENSE` — Project license information.
- `SECURITY.md` — Security policy and vulnerability reporting instructions.
- `AGENTS.md` — Instructions and guidance for AI agents and automation tools.
- `.gitignore` — Specifies files and directories ignored by Git.
- `.editorconfig` — Editor configuration for consistent coding styles.
- `.prettierrc.json` — Prettier code formatting configuration.
- `.prettierignore` — Files ignored by Prettier.
- `.pre-commit-config.yaml` — Pre-commit hooks configuration.
- `.env.example` — Example environment variables configuration.
- `package.json` — Node.js package metadata, scripts, and dependencies.
- `pnpm-lock.yaml` — pnpm dependency lock file.
- `pnpm-workspace.yaml` — pnpm workspace and package-import configuration.
- `tsconfig.json` — TypeScript compiler configuration.

---

## Directory Details

### `.config/`

Project configuration files.

Stores reusable configuration files used by the project.
Helps keep the repository root clean and organized.

Common examples:

- .config/tool-config/
- .config/templates/
- .config/settings/

### `docs/`

Project documentation and technical references.

Contains user guides, technical references, and project records.

Possible sections include:

- architecture and design decisions
- developer and contributor guides
- user guides and command references
- generated repository documentation

Doc Gen writes `docs/project_structure.md` by default, but the
destination is configurable.

### `public/`

Static files copied or served without compilation.

May contain images, icons, manifests, and other assets.

### `scripts/`

Utility scripts for development or automation.

May include deployment scripts, maintenance tools, or helpers.

### `src/`

JavaScript or TypeScript application source code.

Contains runtime modules, entry points, and shared source code.

### `tests/`

Automated tests.

Contains unit tests and integration tests.
Ensures code reliability and correctness.

---

## Notes

- Temporary files, caches, and environment directories are excluded.
- Structure is generated automatically using DocGen.
