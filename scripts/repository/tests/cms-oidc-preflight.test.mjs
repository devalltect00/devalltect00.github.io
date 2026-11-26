/** Run the standard-library Python security cases from the normal repository test command. No hosted requests. */
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import test from "node:test";

test("identity preflight rejects stale context and never reports runtime credentials", () => {
  const result = spawnSync(
    process.platform === "win32" ? "python" : "python3",
    [
      "-B",
      "-m",
      "unittest",
      "discover",
      "-s",
      "scripts/testing",
      "-p",
      "test_oidc_preflight.py",
    ],
    {
      cwd: fileURLToPath(new URL("../../../", import.meta.url)),
      encoding: "utf8",
      timeout: 30000,
      windowsHide: true,
    },
  );
  assert.equal(
    result.status,
    0,
    "Identity preflight regressions failed: " + result.stderr,
  );
});
