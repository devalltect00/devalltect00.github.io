/** CI inputs are data, never shell source. This creates only a new ignored public intent, never credentials or approval. */
import fs from "node:fs";
import { assertDispatchContext } from "../../src/core/cms/dispatch.ts";
import { buildRequestSchema } from "../../src/server/cms/request.ts";
try {
  if (process.argv.length !== 2) throw Error();
  assertDispatchContext(process.env);
  const request = buildRequestSchema.parse({
    schema_version: 1,
    scope: "astro-project-profile-build-v1",
    release_id: process.env.HOMEPAGE_CMS_RELEASE,
    sequence: process.env.HOMEPAGE_EXPECTED_SEQUENCE,
    snapshot_digest: process.env.HOMEPAGE_EXPECTED_SNAPSHOT_DIGEST,
    projection_digest: process.env.HOMEPAGE_EXPECTED_PROJECTION_DIGEST,
  });
  if (!fs.existsSync(".cms-cache")) fs.mkdirSync(".cms-cache");
  const stat = fs.lstatSync(".cms-cache");
  if (!stat.isDirectory() || stat.isSymbolicLink()) throw Error();
  fs.writeFileSync(".cms-cache/build-request.json", JSON.stringify(request), {
    flag: "wx",
  });
  console.log(
    "Exact public Astro build request created. No deployment or CMS receipt recorded.",
  );
} catch {
  console.error(
    "CMS_BUILD_REQUEST_INVALID: verify the public identity inputs. Existing request files are retained.",
  );
  process.exitCode = 1;
}
