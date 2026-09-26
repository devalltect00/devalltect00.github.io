import type { Locale } from "../types/site";

export const DEFAULT_PORTFOLIO_URL = "https://devalltect-portfolio.vercel.app/";
export const DEFAULT_DOCUMENTATION_URL =
  "https://devalltect00.github.io/devalltect-docs/";

/**
 * Resolve a public URL override without allowing an invalid build variable to
 * break navigation. Blank, malformed, and non-HTTP values use the default.
 */
export function resolvePublicUrl(
  value: string | undefined,
  fallback: string,
): string {
  const candidate = value?.trim();
  if (!candidate) return fallback;

  try {
    const url = new URL(candidate);
    return ["http:", "https:"].includes(url.protocol) ? url.href : fallback;
  } catch {
    return fallback;
  }
}

export const site = {
  name: "Devalltect",
  author: "Rizky Fernandes",
  email: "devalltect00@gmail.com",
  url: "https://devalltect00.github.io/",
  portfolioUrl: resolvePublicUrl(
    import.meta.env.PUBLIC_PORTFOLIO_URL,
    DEFAULT_PORTFOLIO_URL,
  ),
  documentationUrl: resolvePublicUrl(
    import.meta.env.PUBLIC_DOCUMENTATION_URL,
    DEFAULT_DOCUMENTATION_URL,
  ),
  githubUrl: "https://github.com/devalltect00",
  gitlabUrl: "https://gitlab.com/devalltect00",
  linkedInUrl: "https://www.linkedin.com/in/rizky-purwanto-fernandes/",
} as const;

export const localePath = (locale: Locale) => (locale === "en" ? "/" : "/id/");
