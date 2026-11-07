// @vitest-environment node
/** Filesystem/atomicity regressions guard explicit immutable preparation rather than a mutable latest directory. */
import { it, expect, afterEach } from "vitest";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import {
  prepare,
  readPrepared,
  selectRelease,
} from "../src/server/cms/cache.ts";
import { hash, canonical } from "../src/server/cms/snapshot.ts";
import { astroPackage } from "./fixtures/cms/astro-package.mjs";
const base = JSON.parse(
    fs.readFileSync(
      new URL("./fixtures/cms/approved-package.json", import.meta.url),
      "utf8",
    ),
  ),
  profile = JSON.parse(
    fs.readFileSync(
      new URL("./fixtures/cms/public-profile.json", import.meta.url),
      "utf8",
    ),
  );
const roots: string[] = [];
function root() {
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), "astro-cache-test-"));
  roots.push(dir);
  return dir;
}
afterEach(() => {
  for (const dir of roots.splice(0))
    fs.rmSync(dir, { recursive: true, force: true });
});
it("prepares a complete immutable release, validates reads and treats same identity retries as idempotent", async () => {
  const dir = root(),
    pkg = astroPackage(base, profile),
    bytes = Buffer.from(JSON.stringify(pkg));
  const first = await prepare(dir, async () => bytes);
  expect(readPrepared(dir, pkg.release_id)).toEqual(first);
  expect(await prepare(dir, async () => bytes)).toEqual(first);
  expect(fs.readdirSync(path.join(dir, ".cms-cache/releases"))).toEqual([
    pkg.release_id,
  ]);
});
it("preserves retained releases when source or conflicting immutable input fails", async () => {
  const dir = root(),
    pkg = astroPackage(base, profile);
  await prepare(dir, async () => Buffer.from(JSON.stringify(pkg)));
  await expect(
    prepare(dir, async () => {
      throw new Error("RAW_PRIVATE_DIAGNOSTIC");
    }),
  ).rejects.toThrow(/^CMS_CONTENT_INVALID$/);
  const bad = structuredClone(pkg);
  bad.sequence = "2";
  await expect(
    prepare(dir, async () => Buffer.from(JSON.stringify(bad))),
  ).rejects.toThrow("CMS_CONTENT_INVALID");
  expect(readPrepared(dir, pkg.release_id).sequence).toBe("1");
});
it("rejects changed bytes, incomplete receipts and traversal selectors without authored fallback", async () => {
  const dir = root(),
    pkg = astroPackage(base, profile);
  await prepare(dir, async () => Buffer.from(JSON.stringify(pkg)));
  const folder = path.join(dir, ".cms-cache/releases", pkg.release_id);
  fs.appendFileSync(path.join(folder, "package.json"), " ");
  expect(() => readPrepared(dir, pkg.release_id)).toThrow(
    "CMS_CONTENT_INVALID",
  );
  fs.unlinkSync(path.join(folder, "complete.json"));
  expect(() => readPrepared(dir, pkg.release_id)).toThrow(
    "CMS_CONTENT_INVALID",
  );
  expect(selectRelease(undefined)).toBeUndefined();
  for (const bad of ["", "../x", "/tmp/x", pkg.release_id + " "])
    expect(() => selectRelease(bad)).toThrow("CMS_CONTENT_INVALID");
});
it("rejects directory junctions/symlinks before importing a prepared release", async () => {
  const dir = root(),
    target = root();
  await prepare(target, async () =>
    Buffer.from(JSON.stringify(astroPackage(base, profile))),
  );
  fs.symlinkSync(
    path.join(target, ".cms-cache"),
    path.join(dir, ".cms-cache"),
    process.platform === "win32" ? "junction" : "dir",
  );
  expect(() => readPrepared(dir, base.release_id)).toThrow(
    "CMS_CONTENT_INVALID",
  );
});
it("does not accept copied completion metadata for a different canonical projection", async () => {
  const dir = root(),
    pkg = astroPackage(base, profile);
  await prepare(dir, async () => Buffer.from(JSON.stringify(pkg)));
  const folder = path.join(dir, ".cms-cache/releases", pkg.release_id);
  pkg.content.records[0].data.project.translations.en.title = "Changed";
  pkg.snapshot_digest = hash(canonical(pkg.content));
  fs.writeFileSync(path.join(folder, "package.json"), JSON.stringify(pkg));
  expect(() => readPrepared(dir, pkg.release_id)).toThrow(
    "CMS_CONTENT_INVALID",
  );
});
