import type { Project } from "../types/site";

export const projects: Project[] = [
  {
    id: "custy",
    name: "Custy",
    description: {
      en: "Configurable Git workflow, versioning, changelog, backup, cleanup, and release automation.",
      id: "Otomasi alur kerja Git, versi, changelog, pencadangan, pembersihan, dan rilis yang dapat dikonfigurasi.",
    },
    category: "Release automation",
    version: "v2.1.2",
    technologies: ["Python", "Git", "Typer", "Jinja2"],
    repository: "https://github.com/devalltect00/Custy",
    documentation: "https://devalltect00.github.io/devalltect-docs/docs/custy/",
    accent: "magenta",
  },
  {
    id: "reflow",
    name: "Reflow",
    description: {
      en: "Repository-aware tag conversion, release recovery, and container image publishing.",
      id: "Konversi tag, pemulihan rilis, dan publikasi image container yang memahami konteks repositori.",
    },
    category: "Release maintenance",
    version: "v1.0.2",
    technologies: ["Python", "Git", "Docker", "Rich"],
    repository: "https://github.com/devalltect00/Git-Reflow",
    documentation:
      "https://devalltect00.github.io/devalltect-docs/docs/reflow/",
    accent: "cyan",
  },
  {
    id: "path-header-scanner",
    name: "Path Header Scanner",
    description: {
      en: "Preview, validate, and apply consistent path headers across source code and documentation.",
      id: "Pratinjau, validasi, dan terapkan path header yang konsisten pada kode sumber dan dokumentasi.",
    },
    category: "Code maintenance",
    version: "v1.0.1",
    technologies: ["Python", "Typer", "TOML", "Docker"],
    repository: "https://github.com/devalltect00/Path-Header-Scanner",
    documentation:
      "https://devalltect00.github.io/devalltect-docs/docs/path-header-scanner/",
    accent: "teal",
  },
  {
    id: "doc-gen",
    name: "Doc-Gen",
    description: {
      en: "Generate, print, and analyze repository structure documentation in Markdown.",
      id: "Buat, tampilkan, dan analisis dokumentasi struktur repositori dalam Markdown.",
    },
    category: "Documentation tooling",
    version: "v1.0.2",
    technologies: ["Python", "Markdown", "Typer", "Rich"],
    repository: "https://github.com/devalltect00/Doc-Gen",
    documentation:
      "https://devalltect00.github.io/devalltect-docs/docs/doc-gen/",
    accent: "navy",
  },
];
