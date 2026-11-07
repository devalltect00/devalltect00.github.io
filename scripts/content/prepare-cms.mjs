/** Public-only explicit preparation never loads env files, CMS credentials or working data. */
import { parseArgs } from "node:util";
import { publicSource } from "../../src/server/cms/http.ts";
import { prepare } from "../../src/server/cms/cache.ts";
try {
  const { values } = parseArgs({
    options: {
      origin: { type: "string" },
      "allow-loopback": { type: "boolean", default: false },
    },
    strict: true,
    allowPositionals: false,
  });
  if (!values.origin) throw Error();
  const s = await prepare(
    process.cwd(),
    publicSource(values.origin, values["allow-loopback"]),
  );
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
