import { describe, expect, it } from "vitest";
import { projects } from "../src/data/projects";
import {
  DEFAULT_DOCUMENTATION_URL,
  DEFAULT_PORTFOLIO_URL,
  resolvePublicUrl,
  site,
} from "../src/data/site";
import { translations } from "../src/i18n";

describe("homepage content", () => {
  it("keeps project identifiers unique and links secure", () => {
    expect(new Set(projects.map((project) => project.id)).size).toBe(
      projects.length,
    );
    for (const project of projects) {
      expect(project.repository).toMatch(/^https:\/\//);
      expect(project.documentation).toMatch(/^https:\/\//);
      expect(project.description.en).not.toHaveLength(0);
      expect(project.description.id).not.toHaveLength(0);
    }
  });

  it("provides both supported languages", () => {
    expect(Object.keys(translations)).toEqual(["en", "id"]);
    expect(translations.en.hero.primaryAction).toBeTruthy();
    expect(translations.id.hero.primaryAction).toBeTruthy();
  });

  it("uses valid public URL overrides and safe defaults", () => {
    expect(
      resolvePublicUrl("https://example.com/portfolio", DEFAULT_PORTFOLIO_URL),
    ).toBe("https://example.com/portfolio");
    expect(resolvePublicUrl("javascript:alert(1)", DEFAULT_PORTFOLIO_URL)).toBe(
      DEFAULT_PORTFOLIO_URL,
    );
    expect(resolvePublicUrl("  ", DEFAULT_DOCUMENTATION_URL)).toBe(
      DEFAULT_DOCUMENTATION_URL,
    );
    expect(site.portfolioUrl).toMatch(/^https:\/\//);
    expect(site.documentationUrl).toMatch(/^https:\/\//);
  });
});
