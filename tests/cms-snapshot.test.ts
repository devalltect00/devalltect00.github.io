// @vitest-environment node
/** Adversarial public projection tests use real Go wire fixtures and preserve exact lossless release identity. */
import { describe, it, expect } from "vitest";
import fs from "node:fs";
import {
  inspectSnapshot,
  hash,
  canonical,
  excerpt,
  packageLimit,
} from "../src/server/cms/snapshot.ts";
import { astroPackage } from "./fixtures/cms/astro-package.mjs";
const base = JSON.parse(
  fs.readFileSync(
    new URL("./fixtures/cms/approved-package.json", import.meta.url),
    "utf8",
  ),
);
const profile = JSON.parse(
  fs.readFileSync(
    new URL("./fixtures/cms/public-profile.json", import.meta.url),
    "utf8",
  ),
);
/** Rehash only for malformed producer tests; checksums do not authenticate approval. */
function pack(value: unknown): Buffer {
  return Buffer.from(JSON.stringify(value));
}
/** Mutate unknown fixture payloads only after an object check; adversarial test data never weakens domain DTO types. */
function setField(
  value: unknown,
  keys: (string | number)[],
  replacement: unknown,
): void {
  let parent = value;
  for (const key of keys.slice(0, -1)) {
    if (!parent || typeof parent !== "object") throw Error("fixture shape");
    parent = Reflect.get(parent, key) as unknown;
  }
  if (!parent || typeof parent !== "object") throw Error("fixture target");
  Reflect.set(parent, keys.at(-1)!, replacement);
}
describe("public gateway projection", () => {
  it("selects only enabled Astro projects, preserves explicit order and localized optional facts", () => {
    const s = inspectSnapshot(pack(astroPackage(base, profile)));
    expect(s.projects).toHaveLength(2);
    expect(s.projects[0].localizedName?.id).toBe("Proyek tanpa tautan");
    expect(s.projects[0].repository).toBe("");
    expect(s.projects[1].version).toBe("v1.2.3");
    expect(s.projects[1].technologies).toEqual(["React JS"]);
    expect(s.contacts[0].href).toBe("mailto:profile%40example.com");
    expect(s.profile?.displayName).toBe("Public Profile Fixture");
    expect(canonical("<&>")).toBe('"\\u003c\\u0026\\u003e"');
  });
  it("keeps full int64 sequence text, empty authoritative projects/profile and deterministic projection", () => {
    const p = astroPackage(base, profile);
    p.sequence = "9223372036854775807";
    expect(inspectSnapshot(pack(p)).sequence).toBe(p.sequence);
    const empty = inspectSnapshot(pack(astroPackage(base, profile, "empty")));
    expect(empty.projects).toEqual([]);
    expect(empty.profile).toBeNull();
    expect(empty.contacts).toEqual([]);
    expect(inspectSnapshot(pack(p)).projectionDigest).toBe(
      inspectSnapshot(pack(p)).projectionDigest,
    );
  });
  it("rejects private/malformed selected data and conflicting global identities despite a recomputed digest", () => {
    const patches: Array<[(string | number)[], unknown]> = [
      [["content", "records", 0, "data", "project", "row_version"], 1],
      [
        [
          "content",
          "records",
          0,
          "data",
          "project",
          "metadata",
          "links",
          0,
          "url",
        ],
        "javascript:alert(1)",
      ],
      [
        [
          "content",
          "records",
          0,
          "data",
          "project",
          "metadata",
          "links",
          0,
          "url",
        ],
        "https://example.com/?token=x",
      ],
      [
        ["content", "records", 0, "data", "project", "metadata", "channels"],
        [
          {
            channel: "astro",
            enabled: true,
            position: 0,
            featured_position: null,
          },
          {
            channel: "astro",
            enabled: true,
            position: 0,
            featured_position: null,
          },
        ],
      ],
      [
        ["content", "records", 1, "data", "technologies", 0, "display_name"],
        "Conflicting technology",
      ],
      [
        ["content", "records", 0, "data", "project", "media", 0, "media_id"],
        "77777777-7777-4777-8777-777777777777",
      ],
      [
        [
          "content",
          "records",
          4,
          "data",
          "profile",
          "translations",
          "id",
          "professional_summary",
        ],
        "",
      ],
      [["content", "records", 4, "data", "profile", "owner_id"], "private"],
    ];
    for (const [keys, value] of patches) {
      const p = astroPackage(base, profile);
      setField(p, keys, value);
      p.snapshot_digest = hash(canonical(p.content));
      expect(() => inspectSnapshot(pack(p))).toThrow("CMS_CONTENT_INVALID");
    }
    const duplicate = astroPackage(base, profile);
    duplicate.content.records.push(duplicate.content.records[0]);
    duplicate.snapshot_digest = hash(canonical(duplicate.content));
    expect(() => inspectSnapshot(pack(duplicate))).toThrow(
      "CMS_CONTENT_INVALID",
    );
  });
  it("rejects invalid identity/hash, duplicate/prototype JSON members, deep structures, malformed UTF8 and oversized bodies", () => {
    const p = astroPackage(base, profile);
    for (const v of [
      { ...p, sequence: 1 },
      { ...p, sequence: "9223372036854775808" },
      { ...p, release_id: "../release" },
      { ...p, snapshot_digest: "0".repeat(64) },
    ])
      expect(() => inspectSnapshot(pack(v))).toThrow("CMS_CONTENT_INVALID");
    for (const v of [
      Buffer.from('{"a":1,"a":2}'),
      Buffer.from('{"constructor":{}}'),
      Buffer.from("[".repeat(65) + "0" + "]".repeat(65)),
      Buffer.from([0xff]),
      Buffer.alloc(packageLimit + 1),
    ])
      expect(() => inspectSnapshot(v)).toThrow("CMS_CONTENT_INVALID");
  });
  it("bounds Unicode-safe excerpts without fabricating ratings, links or full profile content", () => {
    const value = "🧭 ".repeat(300);
    const result = excerpt(value);
    expect([...result].length).toBeLessThanOrEqual(320);
    expect(result.endsWith("…")).toBe(true);
    expect(excerpt("Line one\nLine two")).toBe("Line one Line two");
    expect(excerpt("x".repeat(500))).toHaveLength(320);
  });
});
