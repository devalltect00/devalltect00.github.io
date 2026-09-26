import assert from "node:assert/strict";
import test from "node:test";

import {
  GITHUB,
  GITLAB,
  buildRequests,
  discoverRepository,
  getDescription,
  getTopics,
  normalizeRepositoryUrl,
  parseArguments,
  resolveCommand,
  runApiRequest,
} from "../src/sync_metadata.mjs";

test("loads and normalizes provider-specific metadata", () => {
  const packageData = {
    description: "  Documentation portal  ",
    devalltect: {
      github: { topics: ["Developer Tools", "docs", "docs", ""] },
    },
  };

  assert.equal(getDescription(packageData), "Documentation portal");
  assert.deepEqual(getTopics(packageData, "github"), [
    "developer-tools",
    "docs",
  ]);
  assert.deepEqual(getTopics(packageData, "gitlab"), []);
});

test("rejects invalid metadata instead of coercing unsafe values", () => {
  assert.throws(() => getDescription({}), /non-empty string/u);
  assert.throws(
    () => getTopics({ devalltect: { github: { topics: "docs" } } }, "github"),
    /must be an array/u,
  );
  assert.throws(
    () => getTopics({ devalltect: { github: { topics: [{}] } } }, "github"),
    /only strings/u,
  );
});

test("normalizes provider clone URLs and validates their host", () => {
  assert.equal(
    normalizeRepositoryUrl(
      "https://github.com/devalltect00/devalltect-docs.git",
      GITHUB,
    ),
    "devalltect00/devalltect-docs",
  );
  assert.equal(
    normalizeRepositoryUrl("git@gitlab.com:group/subgroup/docs.git", GITLAB),
    "group/subgroup/docs",
  );
  assert.throws(
    () =>
      normalizeRepositoryUrl(
        "https://example.com/devalltect00/devalltect-docs.git",
        GITHUB,
      ),
    /does not point to GitHub/u,
  );
  assert.throws(
    () => normalizeRepositoryUrl("git@gitlab.com:group/../docs.git", GITLAB),
    /invalid repository path/u,
  );
});

test("discovers a provider target without exposing rejected URLs", () => {
  const responses = new Map([
    [["remote"].join("\0"), "origin\nbackup\ngitlab"],
    [["remote", "get-url", "--", "github"].join("\0"), ""],
    [
      ["remote", "get-url", "--", "origin"].join("\0"),
      "https://example.com/private.git",
    ],
    [
      ["remote", "get-url", "--", "backup"].join("\0"),
      "git@github.com:devalltect00/devalltect-docs.git",
    ],
  ]);
  const target = discoverRepository(GITHUB, (args) => {
    const value = responses.get(args.join("\0"));
    if (value === undefined) {
      throw new Error("missing fixture");
    }
    return value;
  });

  assert.deepEqual(target, {
    remote: "backup",
    repository: "devalltect00/devalltect-docs",
  });
});

test("builds provider API requests with encoded, bounded payloads", () => {
  const githubRequests = buildRequests(
    GITHUB,
    "owner/docs",
    "Docs",
    Array.from({ length: 25 }, (_, index) => `topic-${index}`),
  );
  assert.equal(githubRequests.length, 2);
  assert.equal(githubRequests[1].payload.names.length, 20);
  assert.equal(githubRequests[1].args.at(-3), "repos/owner/docs/topics");

  const gitlabRequests = buildRequests(GITLAB, "group/subgroup/docs", "Docs", [
    "documentation",
  ]);
  assert.equal(gitlabRequests.length, 1);
  assert.ok(
    gitlabRequests[0].args.includes("projects/group%2Fsubgroup%2Fdocs"),
  );
  assert.deepEqual(gitlabRequests[0].payload, {
    description: "Docs",
    topics: ["documentation"],
  });
});

test("dry-run prints a request without invoking the provider CLI", () => {
  let invoked = false;
  const output = [];
  runApiRequest(
    { executable: "gh", args: ["api", "repos/owner/docs"], payload: {} },
    {
      dryRun: true,
      runner: () => {
        invoked = true;
      },
      output: (line) => output.push(line),
    },
  );

  assert.equal(invoked, false);
  assert.equal(output.length, 2);
});

test("resolves the provider CLI to the first platform-reported path", () => {
  const invocations = [];
  const executable = resolveCommand("gh", {
    platform: "win32",
    runner: (command, args) => {
      invocations.push({ command, args });
      return {
        error: null,
        status: 0,
        stdout:
          "C:\\Program Files\\GitHub CLI\\gh.exe\r\nD:\\tools\\gh.exe\r\n",
      };
    },
  });

  assert.deepEqual(invocations, [{ command: "where.exe", args: ["gh"] }]);
  assert.equal(executable, "C:\\Program Files\\GitHub CLI\\gh.exe");
});

test("runs a provider request through its resolved executable path", () => {
  const invocations = [];
  const output = [];
  runApiRequest(
    { executable: "gh", args: ["api", "repos/owner/docs"], payload: {} },
    {
      executable: "C:\\Program Files\\GitHub CLI\\gh.exe",
      runner: (command, args) => {
        invocations.push({ command, args });
        return { error: null, status: 0 };
      },
      output: (line) => output.push(line),
    },
  );

  assert.equal(invocations[0].command, "C:\\Program Files\\GitHub CLI\\gh.exe");
  assert.deepEqual(invocations[0].args, ["api", "repos/owner/docs"]);
  assert.match(output[0], /^\+ gh api/u);
});

test("parses dry-run and provider selection", () => {
  assert.deepEqual(parseArguments(["--dry-run", "--provider", "gitlab"]), {
    dryRun: true,
    provider: "gitlab",
    help: false,
  });
  assert.throws(() => parseArguments(["--provider", "other"]), /requires/u);
  assert.throws(() => parseArguments(["--unknown"]), /unknown argument/u);
});
