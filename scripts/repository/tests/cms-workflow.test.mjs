/** Guard the candidate workflow's authority and input boundary; a build artifact is never deployment/job-completion authority. */
import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
const source = readFileSync(
  new URL("../../../.github/workflows/cms-build-preview.yml", import.meta.url),
  "utf8",
);
test("CMS candidate workflow has only read repository permissions and no deployment authority", () => {
  assert.match(source, /permissions:\s*\n\s+contents: read/);
  assert.doesNotMatch(
    source,
    /(?:pages|id-token|actions|contents): write|deploy-pages|configure-pages|upload-pages-artifact|secrets\./,
  );
  assert.match(source, /path: dist\//);
  assert.match(source, /retention-days: 1\b/);
  assert.match(source, /include-hidden-files: true/);
  assert.match(source, /timeout-minutes: 10\b/);
  for (const action of source.matchAll(/uses: (.+)/g))
    assert.match(action[1], /@[a-f0-9]{40}\b/);
});
test("CMS candidate inputs stay in environment data and cannot choose the source origin or shell text", () => {
  assert.match(
    source,
    /HOMEPAGE_CMS_PUBLIC_ORIGIN: \$\{\{ vars\.HOMEPAGE_CMS_PUBLIC_ORIGIN \}\}/,
  );
  for (const line of source.split("\n"))
    if (/\brun:/.test(line)) assert.doesNotMatch(line, /\$\{\{/);
  assert.match(source, /create-build-request\.mjs/);
  assert.match(source, /--request "\$HOMEPAGE_CMS_REQUEST"/);
  assert.doesNotMatch(source, /--allow-loopback|curl|wget|eval\b/);
});
