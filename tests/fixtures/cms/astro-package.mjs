/** Synthetic public fixture factory; original Go member order and explicit Astro channel rules are retained. */
import { createHash } from "node:crypto";
/** Match Go's canonical JSON bytes, not an independently guessed project serializer. */
function digest(value) {
  return createHash("sha256")
    .update(
      JSON.stringify(value).replace(
        /[<>&\u2028\u2029]/g,
        (c) => "\\u" + c.charCodeAt(0).toString(16).padStart(4, "0"),
      ),
    )
    .digest("hex");
}
/** Content-addressed fixture UUIDs avoid overwriting retained test releases; these identifiers confer no approval authority. */
function seal(pkg) {
  pkg.snapshot_digest = digest(pkg.content);
  const h = digest({ sequence: pkg.sequence, content: pkg.snapshot_digest });
  pkg.release_id = `${h.slice(0, 8)}-${h.slice(8, 12)}-4${h.slice(13, 16)}-8${h.slice(17, 20)}-${h.slice(20, 32)}`;
  return pkg;
}
/** Stage changes exercise selected/unselected projects, optional links, locales, contacts and removal. */
export function astroPackage(input, profile, stage = "initial") {
  const pkg = structuredClone(input);
  pkg.release_id =
    stage === "initial"
      ? "11111111-1111-4111-8111-111111111111"
      : stage === "updated"
        ? "77777777-7777-4777-8777-777777777777"
        : "88888888-8888-4888-8888-888888888888";
  pkg.sequence = stage === "initial" ? "1" : stage === "updated" ? "2" : "3";
  if (stage === "empty") {
    pkg.content.records = [];
    pkg.content.assets = [];
    return seal(pkg);
  }
  const record = structuredClone(
    pkg.content.records.find(
      (r) =>
        r.kind === "project" && r.id === "22222222-2222-4222-8222-222222222222",
    ),
  );
  const project = record.data.project;
  project.metadata.channels = [
    {
      channel: "portfolio",
      enabled: true,
      position: 0,
      featured_position: null,
    },
    { channel: "astro", enabled: true, position: 20, featured_position: null },
  ];
  project.metadata.links = [
    { kind: "github", url: "https://github.com/example/public-gateway" },
    { kind: "documentation", url: "https://example.com/public-guide" },
  ];
  project.metadata.version_label = "v1.2.3";
  if (record.data.technologies.length === 0)
    record.data.technologies = [
      {
        id: "77777771-1111-4111-8111-111111111111",
        slug: "react-js",
        display_name: "React JS",
      },
    ];
  const excluded = structuredClone(record);
  excluded.id = "55555555-5555-4555-8555-555555555555";
  excluded.data.project.slug = "portfolio-only-fixture";
  excluded.data.project.metadata.channels = [
    {
      channel: "portfolio",
      enabled: true,
      position: 0,
      featured_position: null,
    },
  ];
  excluded.data.project.translations.en.title = "EXCLUDED_PROJECT_CANARY";
  excluded.data.project.translations.id.title = "EXCLUDED_PROJECT_CANARY";
  const disabled = structuredClone(excluded);
  disabled.id = "66666666-6666-4666-8666-666666666666";
  disabled.data.project.slug = "disabled-astro-fixture";
  disabled.data.project.metadata.channels = [
    {
      channel: "portfolio",
      enabled: true,
      position: 1,
      featured_position: null,
    },
  ];
  const minimal = structuredClone(record);
  minimal.id = "99999999-9999-4999-8999-999999999999";
  minimal.data.project.slug = "minimal-gateway";
  minimal.data.project.metadata.channels = [
    { channel: "astro", enabled: true, position: 0, featured_position: null },
  ];
  minimal.data.project.metadata.links = [];
  minimal.data.project.metadata.version_label = "";
  minimal.data.project.media = [];
  minimal.data.technologies = [];
  minimal.data.project.translations.en.title = "A project without links";
  minimal.data.project.translations.id.title = "Proyek tanpa tautan";
  minimal.data.project.status = "in-progress";
  const p = structuredClone(profile);
  if (stage === "updated") {
    p.display_name = "Updated Public Profile";
    p.public_email = "updated-profile@example.com";
    p.cv = { enabled: false, url: "" };
    p.translations.en.current_focus = "Updated plain <gateway> focus.";
    p.translations.id.current_focus = "Fokus <beranda> yang diperbarui.";
    project.slug = "runtime-new";
    project.translations.en.title = "Updated gateway project";
    project.translations.id.title = "Proyek beranda terbaru";
    project.translations.en.summary = "Updated plain <gateway> text.";
    project.translations.id.summary = "Teks <beranda> yang diperbarui.";
  }
  pkg.content.records = [
    record,
    excluded,
    disabled,
    ...(stage === "initial" ? [minimal] : []),
    ...pkg.content.records.filter(
      (r) => r.kind !== "project" && r.kind !== "profile",
    ),
    { kind: "profile", id: "1", data: { profile: p, technologies: [] } },
  ];
  return seal(pkg);
}
