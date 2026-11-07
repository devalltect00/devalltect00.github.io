/** Versioned anonymous CMS wire contract. Working rows, credentials and URLs for asset retrieval are never accepted. */
import { z } from "zod";
import { isProjectDate } from "./dates.ts";

export const cmsUUID = z
  .string()
  .regex(/^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$/);
export const cmsDigest = z.string().regex(/^[a-f0-9]{64}$/);
const slug = z
  .string()
  .max(200)
  .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/);
const text = z.string().max(20000);
const nonempty = text.min(1);
const position = z.number().int().min(0).max(2147483647);
const link = z.strictObject({
  kind: z
    .string()
    .max(60)
    .regex(/^[a-z][a-z0-9_-]*$/),
  url: z
    .url()
    .max(2048)
    .refine((value) => {
      const u = new URL(value);
      return (
        u.protocol === "https:" &&
        !u.username &&
        !u.password &&
        !u.search &&
        !value.split("#")[0].includes("?") &&
        (!u.port || u.port === "443") &&
        !/\s|\\/.test(value)
      );
    }),
});
const translation = z.strictObject({
  title: nonempty,
  short_title: text,
  summary: nonempty,
  role: text,
});
const mediaText = z.strictObject({ alt_text: nonempty, caption: text });
const paired = (a: string, b: string) => (a === "") === (b === "");
/** Full projected project metadata remains strict even for projects excluded from this destination. */
export const cmsProjectSchema = z.strictObject({
  slug,
  status: z.enum([
    "maintained",
    "in-progress",
    "completed",
    "archived",
    "planned",
  ]),
  translations: z
    .strictObject({ en: translation, id: translation })
    .refine(
      (v) =>
        paired(v.en.short_title, v.id.short_title) &&
        paired(v.en.role, v.id.role),
    ),
  metadata: z.strictObject({
    date: z
      .strictObject({
        value: z.string().refine(isProjectDate),
        precision: z.enum(["year", "month", "day"]),
      })
      .refine(
        (d) => d.value.length === { year: 4, month: 7, day: 10 }[d.precision],
      )
      .nullable(),
    version_label: text,
    license: text,
    types: z.array(z.string().max(100)).max(30),
    links: z.array(link).max(30),
    contributors: z
      .array(
        z.strictObject({
          name: nonempty,
          role: z
            .strictObject({ en: text, id: text })
            .refine((v) => paired(v.en, v.id)),
          links: z.array(link).max(30),
        }),
      )
      .max(50),
    distribution: z.strictObject({
      python: z
        .strictObject({
          package: slug,
          version: nonempty,
          access: z.literal("public"),
          registry_page: link.shape.url,
          index_url: link.shape.url,
        })
        .nullable(),
      containers: z
        .array(
          z.strictObject({
            provider: z.enum(["ghcr", "gitlab"]),
            image: z
              .string()
              .max(512)
              .regex(/^(?:ghcr\.io|registry\.gitlab\.com)\/[a-z0-9._/-]+$/),
            access: z.literal("public"),
          }),
        )
        .max(20),
    }),
    channels: z
      .array(
        z.strictObject({
          channel: z.enum(["portfolio", "astro", "github", "gitlab", "docs"]),
          enabled: z.boolean(),
          position,
          featured_position: position.nullable(),
        }),
      )
      .min(1)
      .max(5),
  }),
  media: z
    .array(
      z.strictObject({
        media_id: cmsUUID,
        translations: z
          .strictObject({ en: mediaText, id: mediaText })
          .refine((v) => paired(v.en.caption, v.id.caption)),
      }),
    )
    .max(100),
});
/** References retain distinct technology labels; no aliases or proficiency are inferred. */
export const cmsTechnologiesSchema = z
  .array(z.strictObject({ id: cmsUUID, slug, display_name: nonempty }))
  .max(30);
/** Public normalized rasters only. SVG, arbitrary paths and private storage URLs are deliberately absent. */
export const cmsAssetSchema = z.strictObject({
  id: cmsUUID,
  mime: z.enum(["image/png", "image/jpeg", "image/webp"]),
  width: z.number().int().min(1).max(8192),
  height: z.number().int().min(1).max(8192),
  digest: cmsDigest,
});
/** The complete approved envelope is hash-checked; only project/Profile bodies are parsed and projected here. */
export const cmsEnvelopeSchema = z.strictObject({
  schema_version: z.literal(1),
  release_id: cmsUUID,
  sequence: z
    .string()
    .regex(/^[1-9][0-9]{0,18}$/)
    .refine((v) => BigInt(v) <= 9223372036854775807n),
  created_at: z
    .string()
    .regex(/^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d{1,9})?Z$/)
    .refine((v) => Number.isFinite(Date.parse(v))),
  snapshot_digest: cmsDigest,
  locale: z.literal("en").optional(),
  content: z.strictObject({
    schema_version: z.literal(1),
    records: z
      .array(
        z.strictObject({
          kind: z.enum([
            "project",
            "profile",
            "skill",
            "experience",
            "education",
            "publication",
          ]),
          id: z.string(),
          data: z.record(z.string(), z.unknown()),
        }),
      )
      .max(2001),
    assets: z.array(cmsAssetSchema).max(256),
  }),
});
export type CmsProject = z.infer<typeof cmsProjectSchema>;
export type CmsAsset = z.infer<typeof cmsAssetSchema>;
export type CmsTechnology = z.infer<typeof cmsTechnologiesSchema>[number];
