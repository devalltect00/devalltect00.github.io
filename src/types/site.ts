export type Locale = "en" | "id";

export interface Destination {
  id: string;
  title: string;
  description: string;
  href: string;
  label: string;
  icon: "portfolio" | "docs" | "github" | "gitlab";
  external: boolean;
}

/** Public contacts use only local icon kinds and explicitly approved destinations. */
export interface ContactLink {
  kind:
    | "email"
    | "website"
    | "github"
    | "gitlab"
    | "linkedin"
    | "documentation"
    | "telegram"
    | "youtube"
    | "devto";
  href: string;
  icon: "mail" | "github" | "gitlab" | "linkedin" | "arrow";
}
export interface Project {
  id: string;
  name: string;
  localizedName?: Record<Locale, string>;
  status?: "maintained" | "in-progress" | "completed" | "archived" | "planned";
  website?: string;
  description: Record<Locale, string>;
  category: string;
  version: string;
  technologies: string[];
  repository: string;
  documentation: string;
  accent: "magenta" | "cyan" | "navy" | "teal";
}

export interface SiteCopy {
  languageName: string;
  navigation: {
    home: string;
    explore: string;
    projects: string;
    principles: string;
    contact: string;
    languageLabel: string;
    themeLabel: string;
  };
  hero: {
    eyebrow: string;
    titleLead: string;
    titleAccent: string;
    description: string;
    primaryAction: string;
    secondaryAction: string;
    signal: string;
  };
  explore: {
    eyebrow: string;
    title: string;
    description: string;
    labels: Record<
      string,
      { title: string; description: string; label: string }
    >;
  };
  projects: {
    eyebrow: string;
    title: string;
    description: string;
    documentation: string;
    repository: string;
    website: string;
    empty: string;
    technologies: string;
    statuses: Record<NonNullable<Project["status"]>, string>;
  };
  principles: {
    eyebrow: string;
    title: string;
    description: string;
    items: Array<{ number: string; title: string; description: string }>;
  };
  contact: {
    eyebrow: string;
    title: string;
    description: string;
    action: string;
    noContacts: string;
    labels: Record<ContactLink["kind"], string>;
  };
  footer: string;
}
