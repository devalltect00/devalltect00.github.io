/** Validate explicit pinned input before Astro can replace dist; only the locked Astro CLI is executed. */
import { spawnSync } from "node:child_process";
import path from "node:path";
import fs from "node:fs";
import { readPrepared, selectRelease } from "../../src/server/cms/cache.ts";
try {
  const id = selectRelease(process.env.HOMEPAGE_CMS_RELEASE);
  if (id) readPrepared(process.cwd(), id);
  const root = path.resolve("node_modules/astro"),
    pkg = JSON.parse(fs.readFileSync(path.join(root, "package.json"), "utf8")),
    cli = path.resolve(root, pkg.bin.astro);
  if (!cli.startsWith(root + path.sep)) throw Error();
  const child = spawnSync(process.execPath, [cli, "build"], {
    stdio: "inherit",
    env: process.env,
  });
  process.exitCode = child.status ?? 1;
} catch {
  console.error(
    "CMS_BUILD_INPUT_INVALID: prepared input is unavailable or changed. Existing dist is retained; no authored fallback was used.",
  );
  process.exitCode = 1;
}
