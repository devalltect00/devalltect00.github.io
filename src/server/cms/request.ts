/** Exact public build intent is process-local input, never an authentication token or a mutable current pointer. */
import fs from "node:fs";
import { z } from "zod";
import { cmsUUID, cmsDigest } from "../../core/cms/schema.ts";
import { parseCmsJSON } from "../../core/cms/strict-json.ts";
import { CmsError, type GatewaySnapshot } from "./snapshot.ts";

export const buildRequestSchema = z.strictObject({
  schema_version: z.literal(1),
  scope: z.literal("astro-project-profile-build-v1"),
  release_id: cmsUUID,
  sequence: z
    .string()
    .regex(/^[1-9][0-9]{0,18}$/)
    .refine((v) => BigInt(v) <= 9223372036854775807n),
  snapshot_digest: cmsDigest,
  projection_digest: cmsDigest,
});
export type BuildRequest = z.infer<typeof buildRequestSchema>;
/** Bounded strict local reads reject duplicate keys, known links and malformed explicit paths before a build can replace dist. */
export function readBuildRequest(filename: string): BuildRequest {
  try {
    if (!filename) throw new CmsError();
    const stat = fs.lstatSync(filename);
    if (
      !stat.isFile() ||
      stat.isSymbolicLink() ||
      stat.size < 1 ||
      stat.size > 8192
    )
      throw new CmsError();
    const raw = fs.readFileSync(filename);
    if (raw.length > 8192) throw new CmsError();
    return buildRequestSchema.parse(
      parseCmsJSON(new TextDecoder("utf-8", { fatal: true }).decode(raw)),
    );
  } catch {
    throw new CmsError();
  }
}
/** Repeat all expected identity checks; an explicit request never selects authored content or a different current release. */
export function assertBuildRequest(
  request: BuildRequest,
  snapshot: GatewaySnapshot,
): void {
  if (
    request.release_id !== snapshot.releaseID ||
    request.sequence !== snapshot.sequence ||
    request.snapshot_digest !== snapshot.snapshotDigest ||
    request.projection_digest !== snapshot.projectionDigest
  )
    throw new CmsError();
}
