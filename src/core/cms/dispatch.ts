/** Public run correlation is data, never Owner authentication, deployment authority or a completion receipt. */
import { z } from "zod";
import { cmsUUID } from "./schema.ts";
const contextSchema = z.strictObject({
  attempt_id: cmsUUID,
  head_sha: z.string().regex(/^[a-f0-9]{40}$/),
  actual_sha: z.string().regex(/^[a-f0-9]{40}$/),
});
/** Legacy local previews have no run context; explicitly supplied workflow context must be complete and match GitHub's source commit. */
export function assertDispatchContext(
  env: Record<string, string | undefined>,
): void {
  const attempt = env.HOMEPAGE_CMS_DISPATCH_ID,
    head = env.HOMEPAGE_EXPECTED_HEAD_SHA;
  if (attempt === undefined && head === undefined) return;
  const value = contextSchema.parse({
    attempt_id: attempt,
    head_sha: head,
    actual_sha: env.GITHUB_SHA,
  });
  if (value.head_sha !== value.actual_sha)
    throw Error("CMS_DISPATCH_CONTEXT_INVALID");
}
