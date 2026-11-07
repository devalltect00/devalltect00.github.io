/** Strict approved public profile contract; disabled CV locators and private fields never enter visitor props. */
import { z } from "zod";

/** Mirror Go rune bounds and reject control/ambiguous URL input without silently repairing approved facts. */
const text = (max: number, multiline = false) =>
  z
    .string()
    .refine(
      (v) =>
        [...v].length <= max &&
        v === v.trim() &&
        !/\p{Cc}/u.test(multiline ? v.replace(/[\r\n\t]/g, "") : v),
    );
const https = text(2048).refine((v) => {
  try {
    const u = new URL(v);
    return (
      u.protocol === "https:" &&
      !u.username &&
      !u.password &&
      !u.search &&
      !v.split("#")[0].includes("?") &&
      (!u.port || u.port === "443") &&
      !/\s|\\/.test(v)
    );
  } catch {
    return false;
  }
});
const translated = z.strictObject({
  professional_summary: text(4000, true).refine((v) => v.length > 0),
  current_focus: text(1000, true),
});
export const cmsProfileSchema = z.strictObject({
  display_name: text(120).refine((v) => v.length > 0),
  handle: text(65).refine(
    (v) => v === "" || /^@?[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$/.test(v),
  ),
  public_email: text(254).refine(
    (v) => v === "" || /^[^\s<>@]+@[^\s<>@]+\.[^\s<>@]+$/.test(v),
  ),
  public_links: z
    .array(
      z.strictObject({
        kind: z.enum([
          "website",
          "github",
          "gitlab",
          "linkedin",
          "documentation",
          "telegram",
          "youtube",
          "devto",
        ]),
        url: https,
      }),
    )
    .max(8)
    .refine((v) => new Set(v.map((a) => a.kind)).size === v.length),
  programming_since: z.number().int().min(1900).max(2200).nullable(),
  cv: z
    .strictObject({
      enabled: z.boolean(),
      url: z.union([z.literal(""), https]),
    })
    .refine((v) => (v.enabled ? v.url !== "" : v.url === "")),
  translations: z
    .strictObject({ en: translated, id: translated })
    .refine((v) => (v.en.current_focus === "") === (v.id.current_focus === "")),
});
export type CmsProfile = z.infer<typeof cmsProfileSchema>;
