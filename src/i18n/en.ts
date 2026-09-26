import type { SiteCopy } from "../types/site";

export const en: SiteCopy = {
  languageName: "English",
  navigation: {
    home: "Home",
    explore: "Explore",
    projects: "Projects",
    principles: "Principles",
    contact: "Connect",
    languageLabel: "View in Indonesian",
    themeLabel: "Change color theme",
  },
  hero: {
    eyebrow: "Software engineering · creative systems",
    titleLead: "Ideas become",
    titleAccent: "useful software.",
    description:
      "A focused gateway to Devalltect projects, documentation, experiments, and the thinking behind them.",
    primaryAction: "Explore the portfolio",
    secondaryAction: "Read the documentation",
    signal: "Pause · Think · Act",
  },
  explore: {
    eyebrow: "01 / Choose a path",
    title: "One identity. Different ways to explore.",
    description:
      "Start with the work, dive into the technical details, or follow development where it happens.",
    labels: {
      portfolio: {
        title: "Interactive portfolio",
        description:
          "Projects, experience, technology, and the story behind the work.",
        label: "Enter portfolio",
      },
      docs: {
        title: "Documentation portal",
        description:
          "Guides, commands, architecture, and release-aware project references.",
        label: "Open documentation",
      },
      github: {
        title: "GitHub workspace",
        description:
          "Source code, releases, issues, and public project activity.",
        label: "Visit GitHub",
      },
      gitlab: {
        title: "GitLab workspace",
        description:
          "Mirrored projects, pipelines, packages, and delivery workflows.",
        label: "Visit GitLab",
      },
    },
  },
  projects: {
    eyebrow: "02 / The ecosystem",
    title: "Small tools, connected by a larger workflow.",
    description:
      "Each project solves a focused problem. Together they support building, maintaining, documenting, and releasing software.",
    documentation: "Documentation",
    repository: "Repository",
  },
  principles: {
    eyebrow: "03 / The approach",
    title: "Built with curiosity. Refined with discipline.",
    description:
      "The work moves between playful exploration and deliberate engineering without treating either as an afterthought.",
    items: [
      {
        number: "01",
        title: "Pause",
        description:
          "Understand the problem, its constraints, and who the result is meant to help.",
      },
      {
        number: "02",
        title: "Think",
        description:
          "Shape the system, test assumptions, and choose the smallest useful abstraction.",
      },
      {
        number: "03",
        title: "Act",
        description:
          "Build, validate, document, learn, and improve the next iteration.",
      },
    ],
  },
  contact: {
    eyebrow: "04 / Continue the journey",
    title: "There is more behind every project.",
    description:
      "Explore the portfolio for the visual story, or open the documentation when you want the technical one.",
    action: "Start with the portfolio",
  },
  footer: "Designed as the front door to the Devalltect ecosystem.",
};
