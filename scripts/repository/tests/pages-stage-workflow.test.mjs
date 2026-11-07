/** Guard the disabled uploader's authority, first-run/code checks and complete-file handoff. */
import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { spawnSync } from "node:child_process";
const source = readFileSync(
  new URL("../../../.github/workflows/cms-pages-stage.yml", import.meta.url),
  "utf8",
).replace(/\r\n/g, "\n");
test("Pages staging is manual, off by default and has no deployment permission", () => {
  assert.match(source, /workflow_dispatch:/);
  assert.match(source, /HOMEPAGE_CMS_STAGE_ENABLED == 'true'/);
  assert.match(source, /environment: cms-pages-staging/);
  assert.match(source, /contents: read/);
  assert.match(source, /actions: read/);
  assert.doesNotMatch(
    source,
    /(?:pages|id-token|actions|contents): write|deploy-pages|configure-pages|upload-pages-artifact|secrets\.|workflow_run:|pull_request:|push:/,
  );
  assert.match(source, /cancel-in-progress: false/);
  for (const action of source.matchAll(/uses: (.+)/g))
    assert.match(action[1], /@[a-f0-9]{40}\b/);
});
test("Pages upload preserves exactly one prepared TAR and finite retention", () => {
  assert.match(source, /artifact\.tar/);
  assert.match(source, /compression-level: 6/);
  assert.match(source, /retention-days: 1\b/);
  assert.match(source, /timeout-minutes: 5\b/);
  assert.match(source, /if-no-files-found: error/);
  assert.match(source, /persist-credentials: false/);
  assert.match(source, /github\.run_attempt == 1/);
  assert.match(source, /github\.ref == 'refs\/heads\/main'/);
  for (const line of source.split("\n")) {
    if (/\brun:/.test(line)) assert.doesNotMatch(line, /\$\{\{/);
    if (line.includes("${{ inputs."))
      assert.match(
        line,
        /^\s*(run-name: CMS Astro staging|CMS_STAGE_ID:|CMS_STAGE_REQUEST:|name: astro-cms-pages-)/,
      );
  }
  assert.ok(
    source.indexOf("Check correlation") <
      source.indexOf("uses: actions/checkout@"),
  );
});
test("Before-checkout Python guard accepts only exact first-run ID and reviewed code", () => {
  const match = source.match(/python3 - <<'PY'\n([\s\S]*?)\n\s+PY/);
  assert.ok(match);
  const lines = match[1].split("\n"),
    indent = Math.min(
      ...lines.filter((l) => l.trim()).map((l) => l.match(/^ */)[0].length),
    ),
    script = lines.map((l) => l.slice(indent)).join("\n");
  const id = "33333333-3333-4333-8333-333333333333",
    sha = "a".repeat(40),
    request = { request_id: id, head_sha: sha };
  for (const change of [
    {},
    { CMS_STAGE_ID: "bad" },
    { GITHUB_SHA: "b".repeat(40) },
    { GITHUB_RUN_ATTEMPT: "2" },
    { GITHUB_REF: "refs/pull/1/merge" },
    {
      CMS_STAGE_REQUEST: '{"request_id":"duplicate","request_id":"duplicate"}',
    },
    { CMS_STAGE_REQUEST: " ".repeat(8193) },
  ]) {
    const result = spawnSync(
      process.env.HOMEPAGE_PYTHON_BINARY ?? "python",
      ["-c", script],
      {
        env: {
          ...process.env,
          CMS_STAGE_REQUEST: JSON.stringify(request),
          CMS_STAGE_ID: id,
          GITHUB_SHA: sha,
          GITHUB_RUN_ATTEMPT: "1",
          GITHUB_REF: "refs/heads/main",
          ...change,
        },
        encoding: "utf8",
        windowsHide: true,
      },
    );
    assert.equal(result.status, Object.keys(change).length ? 1 : 0);
    if (Object.keys(change).length)
      assert.equal(result.stderr.trim(), "PAGES_STAGE_UNVERIFIED");
  }
});
