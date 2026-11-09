// Protected deployment workflow tests execute only its pre-checkout Python guard, never provider actions.
import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import {
  readFileSync,
  mkdtempSync,
  realpathSync,
  existsSync,
  rmSync,
} from "node:fs";
import { tmpdir } from "node:os";
import { join, sep, basename } from "node:path";
import { spawnSync } from "node:child_process";
import test from "node:test";

const workflow = readFileSync(
  new URL("../../../.github/workflows/cms-pages-deploy.yml", import.meta.url),
  "utf8",
);
const authored = readFileSync(
  new URL("../../../.github/workflows/pages.yml", import.meta.url),
  "utf8",
);
const sha = "a".repeat(40),
  id = "66666666-6666-4666-8666-666666666666";

test("CMS deployment is opt-in, first-main-run, protected and credential-separated", () => {
  assert.match(workflow, /workflow_dispatch:/);
  assert.doesNotMatch(
    workflow,
    /pull_request:|pull_request_target:|schedule:|push:/,
  );
  assert.match(workflow, /HOMEPAGE_CMS_DEPLOY_ENABLED == 'true'/);
  assert.match(workflow, /github.run_attempt == 1/);
  assert.match(workflow, /environment: github-pages/);
  assert.match(workflow, /group: github-pages/);
  assert.match(workflow, /cancel-in-progress: false/);
  assert.match(workflow, /actions: read/);
  assert.match(workflow, /pages: write/);
  assert.match(workflow, /id-token: write/);
  assert.match(workflow, /persist-credentials: false/);
  assert.match(workflow, /ref: \$\{\{ github.sha \}\}/);
  assert.match(workflow, /run \. deploy-workflow/);
  assert.doesNotMatch(
    workflow,
    /upload-artifact@|curl |wget |continue-on-error|printenv|set -x/,
  );
  assert.match(authored, /HOMEPAGE_CMS_DEPLOY_ENABLED != 'true'/);
  assert.match(authored, /github-pages-authored-disabled/);
  for (const line of workflow
    .split("\n")
    .filter((line) => /uses: actions\//.test(line))) {
    assert.match(line, /@[a-f0-9]{40}/);
  }
  assert.ok(
    workflow.indexOf("Validate reviewed code") <
      workflow.indexOf("actions/checkout@"),
  );
  assert.ok(
    workflow.indexOf("GITHUB_TOKEN:") >
      workflow.indexOf("Acquire exact artifacts"),
  );
});

test("pre-checkout guard pins reviewed policy and correlation and refuses substitutions before writing files", () => {
  const captured = workflow
    .replaceAll("\r\n", "\n")
    .match(/python3 - <<'PY'\n([\s\S]*?)\n\s+PY/);
  assert.ok(captured);
  const guard = captured[1]
    .split("\n")
    .map((line) => line.slice(10))
    .join("\n");
  for (const kind of [
    "valid",
    "pin",
    "identity",
    "code",
    "workflow-code",
    "repository",
    "workflow",
    "attempt",
    "branch",
    "duplicate-policy",
    "large-request",
  ]) {
    const root = mkdtempSync(join(tmpdir(), "cms-deploy-guard-"));
    try {
      let policy = JSON.stringify({
        head_sha: sha,
        repository: "example/homepage",
        repository_id: "123",
        workflow_ref:
          "example/homepage/.github/workflows/cms-pages-deploy.yml@refs/heads/main",
      });
      if (kind === "duplicate-policy")
        policy = policy.slice(0, -1) + ',"head_sha":"other"}';
      const env = {
        ...process.env,
        RUNNER_TEMP: root,
        CMS_REVIEWED_POLICY: policy,
        CMS_POLICY_SHA256: createHash("sha256").update(policy).digest("hex"),
        CMS_DEPLOYMENT_REQUEST: JSON.stringify({
          request_id: id,
          stage: { head_sha: sha },
        }),
        CMS_DEPLOYMENT_ID: id,
        GITHUB_SHA: sha,
        GITHUB_WORKFLOW_SHA: sha,
        GITHUB_REPOSITORY: "example/homepage",
        GITHUB_REPOSITORY_ID: "123",
        GITHUB_WORKFLOW_REF:
          "example/homepage/.github/workflows/cms-pages-deploy.yml@refs/heads/main",
        GITHUB_REF: "refs/heads/main",
        GITHUB_RUN_ATTEMPT: "1",
      };
      const patch = {
        pin: "CMS_POLICY_SHA256",
        identity: "CMS_DEPLOYMENT_ID",
        code: "GITHUB_SHA",
        "workflow-code": "GITHUB_WORKFLOW_SHA",
        repository: "GITHUB_REPOSITORY_ID",
        workflow: "GITHUB_WORKFLOW_REF",
        attempt: "GITHUB_RUN_ATTEMPT",
        branch: "GITHUB_REF",
      }[kind];
      if (patch) env[patch] = "substituted";
      if (kind === "large-request")
        env.CMS_DEPLOYMENT_REQUEST = "x".repeat(12289);
      const result = spawnSync("python", ["-c", guard], {
        env,
        encoding: "utf8",
        windowsHide: true,
        timeout: 15000,
      });
      if (kind === "valid") {
        assert.equal(result.status, 0, result.stderr);
        assert.equal(
          readFileSync(join(root, "cms-pages-deploy", "policy.json"), "utf8"),
          policy,
        );
      } else {
        assert.notEqual(result.status, 0, kind);
        assert.match(result.stderr, /PAGES_RECEIVER_UNVERIFIED/);
        assert.equal(existsSync(join(root, "cms-pages-deploy")), false);
      }
    } finally {
      const target = realpathSync(root),
        base = realpathSync(tmpdir());
      if (
        !target.startsWith(base + sep) ||
        !basename(target).startsWith("cms-deploy-guard-")
      ) {
        throw new Error("Unsafe test cleanup target");
      }
      rmSync(target, { recursive: true, force: true });
    }
  }
});
