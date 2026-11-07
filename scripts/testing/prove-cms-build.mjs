/** Real static Astro proof. All fixture authority is ephemeral; no deployment, private settings or destination receipt is touched. */
import assert from "node:assert/strict";
import fs from "node:fs/promises";
import path from "node:path";
import { spawn } from "node:child_process";
import { createServer } from "node:http";
import { once } from "node:events";
import { createInterface } from "node:readline";
import { astroPackage } from "../../tests/fixtures/cms/astro-package.mjs";
import { inspectSnapshot } from "../../src/server/cms/snapshot.ts";
import { chromium, expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";
const root = process.cwd(),
  evidence = path.resolve(".cms-cache/proof"),
  workerFlag = process.argv.indexOf("--worker"),
  goFlag = process.argv.indexOf("--go");
await fs.mkdir(evidence, { recursive: true });
const requestRoot = await fs.mkdtemp(path.join(evidence, "requests-"));
/** Child output is retained locally; safe failures print only the stage label, never credentials or upstream bodies. */
async function command(label, executable, args, env = {}, cwd = root) {
  const filename = path.join(evidence, label + ".log"),
    log = await fs.open(filename, "w");
  const child = spawn(executable, args, {
    cwd,
    env: { ...process.env, ...env },
    stdio: ["ignore", log.fd, log.fd],
    windowsHide: true,
  });
  try {
    const code = await new Promise((resolve, reject) => {
      child.once("error", reject);
      child.once("close", resolve);
    });
    assert.equal(code, 0, label + " failed; inspect its local log");
  } finally {
    await log.close();
  }
}
let worker, lines, iterator, sourceServer, origin;
/** Only public origin/identity is read from the existing private test fixture protocol. */
async function line() {
  let timer;
  try {
    return await Promise.race([
      iterator.next().then((v) => {
        assert.ok(!v.done);
        return JSON.parse(v.value);
      }),
      new Promise((_, reject) => {
        timer = setTimeout(() => reject(Error("FIXTURE_TIMEOUT")), 30000);
      }),
    ]);
  } finally {
    clearTimeout(timer);
  }
}
const base = JSON.parse(
    await fs.readFile("tests/fixtures/cms/approved-package.json", "utf8"),
  ),
  profile = JSON.parse(
    await fs.readFile("tests/fixtures/cms/public-profile.json", "utf8"),
  );
let current = astroPackage(base, profile),
  requests = [];
const packages = [];
try {
  if (workerFlag >= 0) {
    const fixture = process.argv[workerFlag + 1];
    assert.ok(fixture);
    worker = spawn(process.execPath, [fixture, "--astro"], {
      stdio: ["pipe", "pipe", "ignore"],
      windowsHide: true,
    });
    lines = createInterface({ input: worker.stdout });
    iterator = lines[Symbol.asyncIterator]();
    origin = (await line()).origin;
  } else {
    sourceServer = createServer((req, res) => {
      requests.push(req.url);
      assert.equal(req.headers.authorization, undefined);
      assert.equal(req.headers.cookie, undefined);
      if (req.url !== "/public/v1/content/en") {
        res.writeHead(404);
        res.end();
        return;
      }
      res.writeHead(200, { "content-type": "application/json" });
      res.end(JSON.stringify(current));
    });
    sourceServer.listen(0, "127.0.0.1");
    await once(sourceServer, "listening");
    origin = "http://127.0.0.1:" + sourceServer.address().port;
  }
  for (const stage of ["initial", "updated", "empty"]) {
    if (stage !== "initial") {
      if (worker) {
        worker.stdin.write((stage === "updated" ? "advance" : "empty") + "\n");
        await line();
      } else current = astroPackage(base, profile, stage);
    }
    const raw = Buffer.from(
        await (
          await fetch(origin + "/public/v1/content/en", { redirect: "error" })
        ).arrayBuffer(),
      ),
      snapshot = inspectSnapshot(raw);
    const requestFile = path.join(requestRoot, stage + ".json");
    const packageFile = path.resolve(
      ".cms-cache/releases",
      snapshot.releaseID,
      "package.json",
    );
    packages.push({ stage, snapshot, requestFile, packageFile });
    await command("prepare-" + stage, process.execPath, [
      "scripts/content/prepare-cms.mjs",
      "--origin",
      origin,
      "--allow-loopback",
    ]);
    if (goFlag >= 0) {
      await command(
        "go-" + stage,
        process.env.HOMEPAGE_GO_BINARY ?? "go",
        [
          "test",
          "-count=1",
          "-run",
          "^TestCompiledAstroSnapshot$",
          "./internal/publishing",
        ],
        {
          CMS_ASTRO_TEST_PACKAGE_FILE: path.resolve(
            ".cms-cache/releases",
            snapshot.releaseID,
            "package.json",
          ),
        },
        process.argv[goFlag + 1],
      );
    }
    if (goFlag >= 0) {
      await command(
        "request-" + stage,
        process.env.HOMEPAGE_GO_BINARY ?? "go",
        [
          "run",
          "./cmd/cms",
          "astro-build",
          "request",
          "--package",
          packageFile,
          "--out",
          requestFile,
        ],
        {},
        process.argv[goFlag + 1],
      );
    } else {
      await fs.writeFile(
        requestFile,
        JSON.stringify({
          schema_version: 1,
          scope: "astro-project-profile-build-v1",
          release_id: snapshot.releaseID,
          sequence: snapshot.sequence,
          snapshot_digest: snapshot.snapshotDigest,
          projection_digest: snapshot.projectionDigest,
        }),
        { flag: "wx" },
      );
    }
    // Preparation itself must enforce the intent before changing a retained cache.
    await command("prepare-intent-" + stage, process.execPath, [
      "scripts/content/prepare-cms.mjs",
      "--origin",
      origin,
      "--allow-loopback",
      "--request",
      requestFile,
    ]);
  }
  if (worker) {
    worker.stdin.write("stop\n");
    await once(worker, "close");
    lines.close();
    worker = undefined;
  } else {
    await new Promise((resolve) => sourceServer.close(resolve));
    sourceServer = undefined;
  }
  // All builds and browser reads below happen after the source has stopped.
  for (const { stage, snapshot, requestFile, packageFile } of packages) {
    await command(
      "build-" + stage,
      process.execPath,
      ["scripts/content/build.mjs"],
      {
        HOMEPAGE_CMS_RELEASE: snapshot.releaseID,
        HOMEPAGE_CMS_REQUEST: requestFile,
      },
    );
    if (goFlag >= 0)
      await command(
        "verify-" + stage,
        process.env.HOMEPAGE_GO_BINARY ?? "go",
        [
          "run",
          "./cmd/cms",
          "astro-build",
          "verify",
          "--package",
          packageFile,
          "--dist",
          path.resolve("dist"),
          "--request",
          requestFile,
        ],
        {},
        process.argv[goFlag + 1],
      );
    if (goFlag >= 0)
      await command(
        "artifact-" + stage,
        process.env.HOMEPAGE_GO_BINARY ?? "go",
        [
          "test",
          "-count=1",
          "-run",
          "^TestCompiledAstroArtifact$",
          "./internal/platform/astrogithub",
        ],
        {
          CMS_ASTRO_ARTIFACT_PACKAGE: packageFile,
          CMS_ASTRO_ARTIFACT_DIST: path.resolve("dist"),
          CMS_ASTRO_PAGES_OUTPUT: path.join(requestRoot, "pages-" + stage),
        },
        process.argv[goFlag + 1],
      );
    if (goFlag >= 0)
      await command("pages-" + stage, process.execPath, [
        "scripts/content/verify-pages-payload.mjs",
        path.join(requestRoot, "pages-" + stage, "candidate.tar.gz"),
        path.join(requestRoot, "pages-" + stage, "pages-evidence.json"),
      ]);
    const marker = JSON.parse(
      await fs.readFile("dist/.well-known/homepage-content.json", "utf8"),
    );
    assert.equal(marker.source, "cms-build-preview");
    assert.equal(marker.scope, "astro-project-profile-preview-v1");
    assert.equal(marker.release_id, snapshot.releaseID);
    assert.equal(marker.sequence, snapshot.sequence);
    assert.equal(marker.snapshot_digest, snapshot.snapshotDigest);
    assert.equal(marker.projection_digest, snapshot.projectionDigest);
    assert.deepEqual(marker.consumed_kinds, ["project", "profile"]);
    for (const [locale, file] of [
      ["en", "dist/index.html"],
      ["id", "dist/id/index.html"],
    ]) {
      const html = await fs.readFile(file, "utf8");
      assert.ok(
        html.includes(
          snapshot.releaseID +
            ":" +
            snapshot.sequence +
            ":" +
            snapshot.snapshotDigest +
            ":" +
            snapshot.projectionDigest,
        ),
      );
      assert.ok(!html.includes("EXCLUDED_PROJECT_CANARY"));
      assert.ok(!html.includes("mailto:devalltect00@gmail.com"));
      assert.ok(!html.includes("Rizky Fernandes"));
      assert.ok(!html.includes("public-cv.pdf"));
      assert.ok(!html.includes("Custy"));
      assert.equal(
        (html.match(/data-cms-project=/g) || []).length,
        stage === "initial" ? 2 : stage === "updated" ? 1 : 0,
      );
      if (stage === "empty") {
        assert.ok(
          html.includes(
            locale === "en"
              ? "No projects have been published for this homepage yet."
              : "Belum ada proyek yang dipublikasikan untuk beranda ini.",
          ),
        );
        assert.ok(
          html.includes(
            locale === "en"
              ? "Public contact links have not been added."
              : "Tautan kontak publik belum ditambahkan.",
          ),
        );
        assert.ok(!html.includes("profile%40example.com"));
      } else {
        assert.ok(
          html.includes(
            stage === "updated"
              ? "Updated Public Profile"
              : "Public Profile Fixture",
          ),
        );
      }
    }
    if (stage !== "updated") continue;
    /** Serve only files discovered in this build; request text can never become an arbitrary filesystem path. */
    const files = new Map();
    async function collect(dir, relative = "") {
      for (const item of await fs.readdir(dir, { withFileTypes: true })) {
        assert.ok(!item.isSymbolicLink());
        const name = path.posix.join(relative, item.name);
        if (item.isDirectory()) await collect(path.join(dir, item.name), name);
        else
          files.set("/" + name, await fs.readFile(path.join(dir, item.name)));
      }
    }
    await collect(path.resolve("dist"));
    assert.ok(
      ![...files.keys()].some(
        (name) => name.includes(".cms-cache") || name.endsWith("package.json"),
      ),
    );
    const mime = {
      ".html": "text/html",
      ".css": "text/css",
      ".js": "text/javascript",
      ".json": "application/json",
      ".svg": "image/svg+xml",
      ".woff2": "font/woff2",
      ".woff": "font/woff",
    };
    const server = createServer((req, res) => {
      const pathname = new URL(req.url, "http://127.0.0.1").pathname;
      const name = pathname.endsWith("/") ? pathname + "index.html" : pathname;
      const bytes = files.get(name);
      res.writeHead(bytes ? 200 : 404, {
        "content-type": mime[path.extname(name)] ?? "application/octet-stream",
      });
      res.end(bytes ?? "");
    });
    server.listen(0, "127.0.0.1");
    await once(server, "listening");
    const url = "http://127.0.0.1:" + server.address().port;
    const browser = await chromium.launch(
      process.env.HOMEPAGE_BROWSER_CHANNEL === "msedge"
        ? { channel: "msedge" }
        : {},
    );
    try {
      for (const [size, viewport] of [
        ["desktop", { width: 1280, height: 900 }],
        ["mobile", { width: 390, height: 844 }],
      ])
        for (const locale of ["en", "id"]) {
          const context = await browser.newContext({
            viewport,
            reducedMotion: "reduce",
          });
          const page = await context.newPage();
          await page.route("**/*", (route) =>
            route.request().url().startsWith(url)
              ? route.continue()
              : route.abort(),
          );
          await page.goto(url + (locale === "en" ? "/" : "/id/"));
          await expect(page.locator("h1")).toBeVisible();
          await expect(page.locator(".hero__description")).toHaveText(
            locale === "en"
              ? "Updated plain <gateway> focus."
              : "Fokus <beranda> yang diperbarui.",
          );
          await expect(page.locator("article[data-cms-project]")).toHaveCount(
            1,
          );
          await expect(page.locator("article[data-cms-project] h3")).toHaveText(
            locale === "en"
              ? "Updated gateway project"
              : "Proyek beranda terbaru",
          );
          await expect(
            page.locator(
              '.contact__links a[href="mailto:updated-profile%40example.com"]',
            ),
          ).toHaveCount(1);
          for (const item of await page.locator(".reveal").all()) {
            await item.scrollIntoViewIfNeeded();
            await expect(item).toHaveCSS("opacity", "1");
          }
          assert.ok(
            await page.evaluate(
              () => document.documentElement.scrollWidth <= innerWidth,
            ),
          );
          const accessibility = await new AxeBuilder({ page })
            .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
            .analyze();
          assert.deepEqual(accessibility.violations, []);
          await page.screenshot({
            path: path.join(evidence, `${locale}-${size}-light.png`),
            fullPage: true,
          });
          const toggle = page.locator("button[data-theme-toggle]");
          await toggle.click();
          await expect(page.locator("html")).toHaveAttribute(
            "data-theme",
            "dark",
          );
          await page.reload();
          await expect(page.locator("html")).toHaveAttribute(
            "data-theme",
            "dark",
          );
          const dark = await new AxeBuilder({ page })
            .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
            .analyze();
          assert.deepEqual(dark.violations, []);
          assert.ok(
            await page.evaluate(
              () => document.documentElement.scrollWidth <= innerWidth,
            ),
          );
          await page.screenshot({
            path: path.join(evidence, `${locale}-${size}-dark.png`),
            fullPage: true,
          });
          await context.close();
        }
    } finally {
      await browser.close();
      await new Promise((resolve) => server.close(resolve));
    }
  }
  // An invalid explicit selector must fail before Astro can empty the last successful local dist.
  const before = await fs.readFile("dist/index.html");
  const failed = spawn(process.execPath, ["scripts/content/build.mjs"], {
    env: { ...process.env, HOMEPAGE_CMS_RELEASE: "../invalid" },
    stdio: "ignore",
    windowsHide: true,
  });
  assert.notEqual(
    await new Promise((resolve) => failed.once("close", resolve)),
    0,
  );
  assert.deepEqual(await fs.readFile("dist/index.html"), before);
  const staleRequest = path.join(requestRoot, "stale.json");
  const stale = JSON.parse(await fs.readFile(packages[0].requestFile, "utf8"));
  stale.projection_digest = "0".repeat(64);
  await fs.writeFile(staleRequest, JSON.stringify(stale));
  const rejected = spawn(process.execPath, ["scripts/content/build.mjs"], {
    env: {
      ...process.env,
      HOMEPAGE_CMS_RELEASE: packages[0].snapshot.releaseID,
      HOMEPAGE_CMS_REQUEST: staleRequest,
    },
    stdio: "ignore",
    windowsHide: true,
  });
  assert.notEqual(
    await new Promise((resolve) => rejected.once("close", resolve)),
    0,
  );
  assert.deepEqual(await fs.readFile("dist/index.html"), before);
  assert.ok(requests.every((name) => name === "/public/v1/content/en"));
  await fs.writeFile(
    path.join(evidence, "result.json"),
    JSON.stringify(
      {
        source:
          workerFlag >= 0 ? "production-signed-worker" : "public-golden-http",
        checks: [
          "go-canonical-public-package",
          "astro-channel-only",
          "localized-projects-profile",
          "optional-links",
          "empty-no-authored-fallback",
          "offline-build-and-browsing",
          "exact-build-marker",
          "exact-release-request",
          ...(goFlag >= 0
            ? [
                "independent-go-semantic-html",
                "go-pages-payload",
                "independent-node-pages-inventory",
              ]
            : []),
          "invalid-pin-retains-dist",
          "en-id-desktop-mobile-light-dark-accessibility",
        ],
        releases: packages.map((p) => ({
          stage: p.stage,
          release_id: p.snapshot.releaseID,
          sequence: p.snapshot.sequence,
        })),
      },
      null,
      2,
    ),
  );
  console.log(
    "Compiled Astro proof passed: pinned source-independent EN/ID builds, selected projects/profile, removals, optional links, themes/accessibility and safe invalid-pin failure. No CMS destination receipt was recorded.",
  );
} finally {
  if (worker) {
    worker.stdin.write("stop\n");
    lines?.close();
  }
  if (sourceServer) {
    sourceServer.closeAllConnections();
    await new Promise((resolve) => sourceServer.close(resolve));
  }
}
