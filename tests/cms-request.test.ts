// @vitest-environment node
/** An exact intent is public input, but stale or malformed requests must fail before replacing prepared input or dist. */
import { expect, it, afterEach } from "vitest";
import fs from "node:fs";
import { createServer } from "node:http";
import { once } from "node:events";
import { spawn } from "node:child_process";
import { fileURLToPath } from "node:url";
import os from "node:os";
import path from "node:path";
import {
  readBuildRequest,
  assertBuildRequest,
  buildRequestSchema,
} from "../src/server/cms/request.ts";
import { inspectSnapshot } from "../src/server/cms/snapshot.ts";
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
const snapshot = inspectSnapshot(
  Buffer.from(JSON.stringify(astroPackage(base, profile))),
);
const request = {
  schema_version: 1,
  scope: "astro-project-profile-build-v1",
  release_id: snapshot.releaseID,
  sequence: snapshot.sequence,
  snapshot_digest: snapshot.snapshotDigest,
  projection_digest: snapshot.projectionDigest,
};
const roots: string[] = [];
afterEach(() => {
  for (const root of roots.splice(0))
    fs.rmSync(root, { recursive: true, force: true });
});
/** Only disposable test-owned directories are used for malformed local intent files. */
function file(raw: string | Buffer): string {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "astro-request-"));
  roots.push(root);
  const name = path.join(root, "request.json");
  fs.writeFileSync(name, raw);
  return name;
}
it("accepts the exact public intent and retains the entire int64 sequence as text", () => {
  const parsed = readBuildRequest(file(JSON.stringify(request)));
  expect(() => assertBuildRequest(parsed, snapshot)).not.toThrow();
  expect(
    buildRequestSchema.parse({ ...request, sequence: "9223372036854775807" })
      .sequence,
  ).toBe("9223372036854775807");
});
it("rejects each independently mismatched identity and malformed envelope without coercing it", () => {
  for (const key of [
    "release_id",
    "sequence",
    "snapshot_digest",
    "projection_digest",
  ] as const) {
    const wrong = { ...request, [key]: "different" };
    expect(() =>
      assertBuildRequest(
        wrong as ReturnType<typeof readBuildRequest>,
        snapshot,
      ),
    ).toThrow("CMS_CONTENT_INVALID");
  }
  for (const raw of [
    "",
    "{",
    JSON.stringify({ ...request, unknown: true }),
    JSON.stringify({ ...request, sequence: "01" }),
    JSON.stringify({ ...request, sequence: "9223372036854775808" }),
    JSON.stringify({ ...request, sequence: 1 }),
    JSON.stringify(request).replace(
      '"schema_version":1',
      '"schema_version":1,"schema_version":1',
    ),
    "x".repeat(8193),
    Buffer.from([0xff]),
  ])
    expect(() => readBuildRequest(file(raw))).toThrow("CMS_CONTENT_INVALID");
});
it("rejects missing, blank, directory and linked request files", () => {
  const name = file(JSON.stringify(request));
  expect(() => readBuildRequest("")).toThrow();
  expect(() => readBuildRequest(name + "-missing")).toThrow();
  expect(() => readBuildRequest(path.dirname(name))).toThrow();
  const alias = path.join(path.dirname(name), "alias");
  // Directory junctions do not require Windows symlink privileges and still exercise non-file rejection.
  fs.symlinkSync(
    path.dirname(name),
    alias,
    process.platform === "win32" ? "junction" : "dir",
  );
  expect(() => readBuildRequest(alias)).toThrow();
});

/** A real CLI child keeps URL/body/path construction in the production entrypoint while the test controls only public fixture input. */
function command(
  script: string,
  args: string[],
  cwd: string,
  env: NodeJS.ProcessEnv = process.env,
): Promise<number | null> {
  const child = spawn(
    process.execPath,
    [fileURLToPath(new URL(script, import.meta.url)), ...args],
    { cwd, env, stdio: "ignore", windowsHide: true },
  );
  return new Promise((resolve, reject) => {
    child.once("error", reject);
    child.once("close", resolve);
  });
}
it("rejects a moving/stale public source before preparing any files and retains the earlier exact release", async () => {
  const filename = file(
    JSON.stringify({ ...request, projection_digest: "0".repeat(64) }),
  );
  const cwd = path.dirname(filename);
  let current = astroPackage(base, profile);
  const server = createServer((_req, res) => {
    res.writeHead(200, { "content-type": "application/json" });
    res.end(JSON.stringify(current));
  });
  server.listen(0, "127.0.0.1");
  await once(server, "listening");
  const address = server.address();
  if (!address || typeof address === "string") throw Error("fixture address");
  const args = [
    "--origin",
    `http://127.0.0.1:${address.port}`,
    "--allow-loopback",
    "--request",
    filename,
  ];
  try {
    expect(
      await command("../scripts/content/prepare-cms.mjs", args, cwd),
    ).not.toBe(0);
    expect(fs.existsSync(path.join(cwd, ".cms-cache"))).toBe(false);
    fs.writeFileSync(filename, JSON.stringify(request));
    expect(await command("../scripts/content/prepare-cms.mjs", args, cwd)).toBe(
      0,
    );
    const releases = path.join(cwd, ".cms-cache/releases");
    const retained = fs.readFileSync(
      path.join(releases, snapshot.releaseID, "package.json"),
    );
    current = astroPackage(base, profile, "updated");
    expect(
      await command("../scripts/content/prepare-cms.mjs", args, cwd),
    ).not.toBe(0);
    expect(fs.readdirSync(releases)).toEqual([snapshot.releaseID]);
    expect(
      fs.readFileSync(path.join(releases, snapshot.releaseID, "package.json")),
    ).toEqual(retained);
  } finally {
    server.closeAllConnections();
    await new Promise<void>((resolve) => server.close(() => resolve()));
  }
}, 15000);
it("treats manual workflow inputs as data and creates its public request exclusively", async () => {
  const name = file("public fixture");
  const cwd = path.dirname(name);
  const env = {
    ...process.env,
    HOMEPAGE_CMS_RELEASE: request.release_id,
    HOMEPAGE_EXPECTED_SEQUENCE: request.sequence,
    HOMEPAGE_EXPECTED_SNAPSHOT_DIGEST: request.snapshot_digest,
    HOMEPAGE_EXPECTED_PROJECTION_DIGEST: request.projection_digest,
  };
  expect(
    await command("../scripts/content/create-build-request.mjs", [], cwd, {
      ...env,
      HOMEPAGE_CMS_RELEASE: "invalid; echo untrusted",
    }),
  ).not.toBe(0);
  expect(fs.existsSync(path.join(cwd, ".cms-cache"))).toBe(false);
  expect(
    await command("../scripts/content/create-build-request.mjs", [], cwd, env),
  ).toBe(0);
  const output = path.join(cwd, ".cms-cache/build-request.json");
  expect(readBuildRequest(output)).toEqual(request);
  expect(
    await command("../scripts/content/create-build-request.mjs", [], cwd, env),
  ).not.toBe(0);
  expect(readBuildRequest(output)).toEqual(request);
}, 15000);
