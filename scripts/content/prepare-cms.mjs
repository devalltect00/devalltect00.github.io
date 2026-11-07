/** Public-only explicit preparation never loads env files, CMS credentials or working data. */
import { parseArgs } from "node:util";
import { publicSource } from "../../src/server/cms/http.ts";
import {
  readBuildRequest,
  assertBuildRequest,
} from "../../src/server/cms/request.ts";
import { inspectSnapshot } from "../../src/server/cms/snapshot.ts";
import { prepare } from "../../src/server/cms/cache.ts";
try {
  const { values } = parseArgs({
    options: {
      origin: { type: "string" },
      request: { type: "string" },
      "allow-loopback": { type: "boolean", default: false },
    },
    strict: true,
    allowPositionals: false,
  });
  if (!values.origin) throw Error();
  const request =
    values.request !== undefined ? readBuildRequest(values.request) : undefined;
  const source = publicSource(values.origin, values["allow-loopback"]);
  const s = await prepare(process.cwd(), async () => {
    const raw = await source();
    if (request) assertBuildRequest(request, inspectSnapshot(raw));
    return raw;
  });
  console.log(
    `Prepared public release ${s.releaseID}, sequence ${s.sequence}, ${s.projects.length} Astro projects.`,
  );
  console.log(
    `Set HOMEPAGE_CMS_RELEASE=${s.releaseID} in the build process environment. No deployment or CMS receipt was recorded.`,
  );
} catch {
  console.error(
    "CMS_PREPARE_FAILED: verify the public source/package. Existing prepared releases are retained.",
  );
  process.exitCode = 1;
}
