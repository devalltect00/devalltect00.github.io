/** Run correlation rejects partial/moving code references and never invents authentication from a public UUID. */
import { it, expect } from "vitest";
import { assertDispatchContext } from "../src/core/cms/dispatch.ts";
const valid = {
  HOMEPAGE_CMS_DISPATCH_ID: "11111111-1111-4111-8111-111111111111",
  HOMEPAGE_EXPECTED_HEAD_SHA: "a".repeat(40),
  GITHUB_SHA: "a".repeat(40),
};
it("allows explicit local preview compatibility without workflow context", () => {
  expect(() => assertDispatchContext({})).not.toThrow();
});
it("binds a public candidate UUID to the exact GitHub source commit", () => {
  expect(() => assertDispatchContext(valid)).not.toThrow();
  expect(() =>
    assertDispatchContext({ ...valid, GITHUB_SHA: "b".repeat(40) }),
  ).toThrow();
});
it("rejects partial, malformed and shell-looking context before any request file is written", () => {
  for (const key of Object.keys(valid) as (keyof typeof valid)[]) {
    expect(() =>
      assertDispatchContext({ ...valid, [key]: undefined }),
    ).toThrow();
    expect(() =>
      assertDispatchContext({ ...valid, [key]: "invalid; echo untrusted" }),
    ).toThrow();
  }
});
