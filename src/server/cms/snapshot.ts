/** Build-only projection validates complete public identity, then retains only selected gateway facts. */
import { createHash } from "node:crypto";
import {
  cmsEnvelopeSchema,
  cmsProjectSchema,
  cmsTechnologiesSchema,
  cmsUUID,
} from "../../core/cms/schema.ts";
import { cmsProfileSchema } from "../../core/cms/profile.ts";
import { parseCmsJSON } from "../../core/cms/strict-json.ts";
import type { Project, ContactLink, Locale } from "../../types/site.ts";
export const packageLimit = (2 << 20) + 8192;
/** Safe constant diagnostics never disclose upstream data, URLs, credentials or local paths. */
export class CmsError extends Error {
  constructor() {
    super("CMS_CONTENT_INVALID");
    this.name = "CmsError";
  }
}
/** Lossless public identity is a build receipt only, without activation or CMS job authority. */
export interface GatewaySnapshot {
  releaseID: string;
  sequence: string;
  snapshotDigest: string;
  createdAt: string;
  projectionDigest: string;
  projects: Project[];
  profile: { displayName: string; summary: Record<Locale, string> } | null;
  contacts: ContactLink[];
}
/** Go's canonical HTML escaping is required before hashing public snapshot content. */
export function hash(value: Uint8Array | string): string {
  return createHash("sha256").update(value).digest("hex");
}
/** Preserve approved member order and match Go JSON HTML escaping before digest comparison. */
export function canonical(value: unknown): string {
  return JSON.stringify(value).replace(
    /[<>&\u2028\u2029]/g,
    (c) => "\\u" + c.charCodeAt(0).toString(16).padStart(4, "0"),
  );
}
/** The gateway deliberately previews approved prose; code-point bounds avoid splitting a Unicode character. */
export function excerpt(value: string, limit = 320): string {
  const text = value.replace(/\s+/gu, " ").trim(),
    chars = [...text];
  if (chars.length <= limit) return text;
  const head = chars.slice(0, limit - 1).join(""),
    space = head.lastIndexOf(" ");
  return (space > limit / 2 ? head.slice(0, space) : head) + "…";
}
/** Public link kinds select only local icons; user text is never evaluated as an icon or component name. */
function icon(kind: string): ContactLink["icon"] {
  return kind === "github" || kind === "gitlab" || kind === "linkedin"
    ? kind
    : "arrow";
}
/** Parse and project one complete approved public package; malformed opt-in never falls back to authored data. */
export function inspectSnapshot(raw: Uint8Array): GatewaySnapshot {
  try {
    if (raw.byteLength < 1 || raw.byteLength > packageLimit)
      throw new CmsError();
    const input = parseCmsJSON(
      new TextDecoder("utf-8", { fatal: true }).decode(raw),
    );
    const envelope = cmsEnvelopeSchema.parse(input);
    if (
      !input ||
      typeof input !== "object" ||
      !("content" in input) ||
      hash(canonical(input.content)) !== envelope.snapshot_digest
    )
      throw new CmsError();
    const ids = new Set<string>(),
      slugs = new Set<string>(),
      assets = new Set<string>(),
      needed = new Set<string>(),
      vocabulary = new Map<string, string>();
    for (const asset of envelope.content.assets) {
      if (assets.has(asset.id)) throw new CmsError();
      assets.add(asset.id);
    }
    const selected: Array<{ position: number; project: Project }> = [];
    let profile: GatewaySnapshot["profile"] = null;
    const contacts: ContactLink[] = [];
    for (const record of envelope.content.records) {
      if (
        (record.kind === "profile"
          ? record.id !== "1"
          : !cmsUUID.safeParse(record.id).success) ||
        ids.has(record.kind + ":" + record.id)
      )
        throw new CmsError();
      ids.add(record.kind + ":" + record.id);
      if (
        Object.keys(record.data).length !== 2 ||
        !Object.hasOwn(record.data, record.kind) ||
        !Object.hasOwn(record.data, "technologies") ||
        !record.data[record.kind] ||
        typeof record.data[record.kind] !== "object" ||
        Array.isArray(record.data[record.kind])
      )
        throw new CmsError();
      const tech = cmsTechnologiesSchema.parse(record.data.technologies);
      if (new Set(tech.map((t) => t.id)).size !== tech.length)
        throw new CmsError();
      for (const t of tech) {
        const identity = JSON.stringify([t.slug, t.display_name]);
        if (vocabulary.has(t.id) && vocabulary.get(t.id) !== identity)
          throw new CmsError();
        vocabulary.set(t.id, identity);
      }
      if (record.kind === "profile") {
        if (tech.length) throw new CmsError();
        const p = cmsProfileSchema.parse(record.data.profile);
        profile = {
          displayName: p.display_name,
          summary: {
            en: excerpt(
              p.translations.en.current_focus ||
                p.translations.en.professional_summary,
            ),
            id: excerpt(
              p.translations.id.current_focus ||
                p.translations.id.professional_summary,
            ),
          },
        };
        if (p.public_email)
          contacts.push({
            kind: "email",
            href: "mailto:" + encodeURIComponent(p.public_email),
            icon: "mail",
          });
        for (const link of p.public_links)
          contacts.push({
            kind: link.kind,
            href: link.url,
            icon: icon(link.kind),
          });
      }
      if (record.kind !== "project") continue;
      const p = cmsProjectSchema.parse(record.data.project);
      if (
        slugs.has(p.slug) ||
        new Set(p.metadata.channels.map((c) => c.channel)).size !==
          p.metadata.channels.length ||
        new Set(p.metadata.links.map((l) => l.kind)).size !==
          p.metadata.links.length ||
        new Set(p.metadata.types).size !== p.metadata.types.length ||
        new Set(p.media.map((m) => m.media_id)).size !== p.media.length
      )
        throw new CmsError();
      slugs.add(p.slug);
      for (const m of p.media) {
        if (!assets.has(m.media_id)) throw new CmsError();
        needed.add(m.media_id);
      }
      const channel = p.metadata.channels.find(
        (c) => c.channel === "astro" && c.enabled,
      );
      if (!channel) continue;
      const link = (...kinds: string[]) =>
        p.metadata.links.find((l) => kinds.includes(l.kind))?.url ?? "";
      selected.push({
        position: channel.position,
        project: {
          id: record.id,
          name: p.translations.en.title,
          localizedName: {
            en: p.translations.en.title,
            id: p.translations.id.title,
          },
          description: {
            en: excerpt(p.translations.en.summary),
            id: excerpt(p.translations.id.summary),
          },
          category: "",
          status: p.status,
          version: p.metadata.version_label,
          technologies: tech.map((t) => t.display_name),
          repository: link("github", "gitlab"),
          documentation: link("documentation"),
          website: link("website", "demo"),
          accent: "magenta",
        },
      });
    }
    if (needed.size !== assets.size) throw new CmsError();
    selected.sort(
      (a, b) =>
        a.position - b.position || a.project.id.localeCompare(b.project.id),
    );
    const accents = ["magenta", "cyan", "teal", "navy"] as const;
    const projects = selected.map((v, i) => ({
      ...v.project,
      accent: accents[i % accents.length],
    }));
    return {
      releaseID: envelope.release_id,
      sequence: envelope.sequence,
      snapshotDigest: envelope.snapshot_digest,
      createdAt: envelope.created_at,
      projectionDigest: hash(canonical({ projects, profile, contacts })),
      projects,
      profile,
      contacts,
    };
  } catch {
    throw new CmsError();
  }
}
