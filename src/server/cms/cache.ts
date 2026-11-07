/** Immutable local public build input; UUID/member paths are fixed and incomplete/symlinked input fails closed. */
import fs from "node:fs";
import path from "node:path";
import { z } from "zod";
import { cmsUUID, cmsDigest } from "../../core/cms/schema.ts";
import { parseCmsJSON } from "../../core/cms/strict-json.ts";
import {
  CmsError,
  hash,
  inspectSnapshot,
  packageLimit,
  type GatewaySnapshot,
} from "./snapshot.ts";
const receipt = z.strictObject({
  schema_version: z.literal(1),
  release_id: cmsUUID,
  sequence: z.string(),
  snapshot_digest: cmsDigest,
  projection_digest: cmsDigest,
  package_sha256: cmsDigest,
});
/** Missing opt-in means authored; malformed/blank explicit opt-in is a configuration error. */
export function selectRelease(value: string | undefined): string | undefined {
  if (value === undefined) return undefined;
  if (!cmsUUID.safeParse(value).success) throw new CmsError();
  return value;
}
/** Every ancestor under the trusted workspace is checked; no upstream filename or arbitrary URL is accepted. */
function folders(root: string, id?: string, create = false): string {
  if (id !== undefined && !cmsUUID.safeParse(id).success) throw new CmsError();
  let dir = path.resolve(root);
  for (const member of [".cms-cache", "releases", ...(id ? [id] : [])]) {
    dir = path.join(dir, member);
    if (create && !fs.existsSync(dir)) fs.mkdirSync(dir);
    const stat = fs.lstatSync(dir);
    if (!stat.isDirectory() || stat.isSymbolicLink()) throw new CmsError();
  }
  return dir;
}
/** Bound allocation and reject known filesystem links before reading complete bytes. */
function file(name: string, limit: number): Buffer {
  const s = fs.lstatSync(name);
  if (!s.isFile() || s.isSymbolicLink() || s.size < 1 || s.size > limit)
    throw new CmsError();
  const raw = fs.readFileSync(name);
  if (raw.length > limit) throw new CmsError();
  return raw;
}
/** Recheck both raw package and strict completion receipt at build time; a directory name is not an integrity proof. */
export function readPrepared(root: string, id: string): GatewaySnapshot {
  try {
    const dir = folders(root, id);
    const raw = file(path.join(dir, "package.json"), packageLimit),
      snapshot = inspectSnapshot(raw),
      r = receipt.parse(
        parseCmsJSON(
          file(path.join(dir, "complete.json"), 8192).toString("utf8"),
        ),
      );
    if (
      snapshot.releaseID !== id ||
      r.release_id !== id ||
      r.sequence !== snapshot.sequence ||
      r.snapshot_digest !== snapshot.snapshotDigest ||
      r.projection_digest !== snapshot.projectionDigest ||
      r.package_sha256 !== hash(raw)
    )
      throw new CmsError();
    return snapshot;
  } catch {
    throw new CmsError();
  }
}
/** Publish only a complete task-owned child directory; retained releases are validated, never overwritten or force-cleaned. */
export async function prepare(
  root: string,
  source: () => Promise<Uint8Array>,
): Promise<GatewaySnapshot> {
  let temporary: string | undefined, parent: string | undefined;
  try {
    const raw = await source(),
      snapshot = inspectSnapshot(raw);
    parent = folders(root, undefined, true);
    const target = path.join(parent, snapshot.releaseID);
    if (fs.existsSync(target)) {
      const old = readPrepared(root, snapshot.releaseID);
      if (
        old.sequence !== snapshot.sequence ||
        old.snapshotDigest !== snapshot.snapshotDigest ||
        old.createdAt !== snapshot.createdAt ||
        old.projectionDigest !== snapshot.projectionDigest
      )
        throw new CmsError();
      return old;
    }
    temporary = fs.mkdtempSync(path.join(parent, ".prepare-"));
    fs.writeFileSync(path.join(temporary, "package.json"), raw, { flag: "wx" });
    fs.writeFileSync(
      path.join(temporary, "complete.json"),
      JSON.stringify({
        schema_version: 1,
        release_id: snapshot.releaseID,
        sequence: snapshot.sequence,
        snapshot_digest: snapshot.snapshotDigest,
        projection_digest: snapshot.projectionDigest,
        package_sha256: hash(raw),
      }),
      { flag: "wx" },
    );
    try {
      fs.renameSync(temporary, target);
    } catch {
      if (!fs.existsSync(target)) throw new CmsError();
      const winner = readPrepared(root, snapshot.releaseID);
      if (
        winner.sequence !== snapshot.sequence ||
        winner.snapshotDigest !== snapshot.snapshotDigest ||
        winner.createdAt !== snapshot.createdAt ||
        winner.projectionDigest !== snapshot.projectionDigest
      )
        throw new CmsError();
    }
    return readPrepared(root, snapshot.releaseID);
  } catch {
    throw new CmsError();
  } finally {
    if (temporary && parent && fs.existsSync(temporary)) {
      if (
        path.dirname(temporary) !== parent ||
        !path.basename(temporary).startsWith(".prepare-") ||
        fs.lstatSync(temporary).isSymbolicLink()
      )
        throw new CmsError();
      fs.rmSync(temporary, { recursive: true });
    }
  }
}
