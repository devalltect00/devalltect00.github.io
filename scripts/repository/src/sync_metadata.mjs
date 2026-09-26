#!/usr/bin/env node

/**
 * Synchronize GitHub and GitLab repository metadata from package.json.
 *
 * Repository identities come from validated Git remotes rather than package
 * metadata, so a token-bearing URL is never required in package.json. Run a
 * dry run first to review the selected repositories and API payloads.
 */

import { readFileSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";
import { pathToFileURL } from "node:url";

export const PACKAGE_JSON = "package.json";
export const GITHUB_TOPIC_LIMIT = 20;

export const GITHUB = Object.freeze({
  key: "github",
  label: "GitHub",
  hostname: "github.com",
  hostAliases: ["www.github.com"],
  remoteCandidates: ["github", "origin", "backup", "upstream"],
  allowSubgroups: false,
  executable: "gh",
});

export const GITLAB = Object.freeze({
  key: "gitlab",
  label: "GitLab",
  hostname: "gitlab.com",
  hostAliases: ["www.gitlab.com"],
  remoteCandidates: ["gitlab", "backup", "origin", "upstream"],
  allowSubgroups: true,
  executable: "glab",
});

const PROVIDERS = Object.freeze({ github: GITHUB, gitlab: GITLAB });

/** Load and parse package.json from the current working directory. */
export function loadPackageJson(path = PACKAGE_JSON) {
  let source;
  try {
    source = readFileSync(path, "utf8");
  } catch (error) {
    throw new Error(`${path} not found; run from the repository root`, {
      cause: error,
    });
  }

  try {
    return JSON.parse(source);
  } catch (error) {
    throw new Error(`${path} contains invalid JSON`, { cause: error });
  }
}

/** Return the non-empty repository description stored in package.json. */
export function getDescription(packageData) {
  const description = packageData?.description;
  if (typeof description !== "string" || !description.trim()) {
    throw new Error('package.json "description" must be a non-empty string');
  }
  return description.trim();
}

/** Return normalized provider topics from package.json's devalltect block. */
export function getTopics(packageData, providerKey) {
  const configured = packageData?.devalltect?.[providerKey]?.topics ?? [];
  if (!Array.isArray(configured)) {
    throw new Error(
      `package.json "devalltect.${providerKey}.topics" must be an array`,
    );
  }

  const topics = [];
  const seen = new Set();
  for (const topic of configured) {
    if (typeof topic !== "string") {
      throw new Error(
        `package.json "devalltect.${providerKey}.topics" must contain only strings`,
      );
    }
    const normalized = topic.trim().toLowerCase().replaceAll(" ", "-");
    if (normalized && !seen.has(normalized)) {
      seen.add(normalized);
      topics.push(normalized);
    }
  }
  return topics;
}

/**
 * Normalize an HTTPS, SSH, Git, or SCP-style clone URL to a repository path.
 * Remote URLs are validated before their repository paths reach provider APIs.
 */
export function normalizeRepositoryUrl(remoteUrl, provider) {
  const value = remoteUrl.trim();
  if (!value || /\s/u.test(value)) {
    throw new Error("empty URL or unescaped whitespace in URL");
  }

  let hostname;
  let pathname;
  if (value.includes("://")) {
    let parsed;
    try {
      parsed = new URL(value);
    } catch {
      throw new Error("malformed remote URL");
    }
    if (!["https:", "http:", "ssh:", "git:"].includes(parsed.protocol)) {
      throw new Error("unsupported remote URL scheme");
    }
    if (parsed.search || parsed.hash) {
      throw new Error("expected a clone URL without a query or fragment");
    }
    hostname = parsed.hostname;
    pathname = parsed.pathname;
  } else {
    const match = value.match(/^(?:[^@/:\s]+@)?([^@/:\s]+):(.+)$/u);
    if (!match) {
      throw new Error("expected a clone URL, not a local path or identifier");
    }
    [, hostname, pathname] = match;
  }

  const acceptedHosts = new Set(
    [provider.hostname, ...provider.hostAliases].map((host) =>
      host.toLowerCase(),
    ),
  );
  if (!hostname || !acceptedHosts.has(hostname.toLowerCase())) {
    throw new Error(`remote does not point to ${provider.label}`);
  }

  let decodedPath;
  try {
    decodedPath = decodeURIComponent(pathname);
  } catch {
    throw new Error("remote URL contains invalid percent encoding");
  }
  decodedPath = decodedPath.replace(/^\/+|\/+$/gu, "");
  if (decodedPath.endsWith(".git")) {
    decodedPath = decodedPath.slice(0, -4);
  }

  const segments = decodedPath.split("/");
  const expectedSegments = provider.allowSubgroups
    ? segments.length >= 2
    : segments.length === 2;
  if (!expectedSegments) {
    const expected = provider.allowSubgroups
      ? "group[/subgroup]/project"
      : "owner/repository";
    throw new Error(`expected a ${expected} repository path`);
  }
  if (
    segments.some(
      (segment) =>
        ["", ".", "..", "-"].includes(segment) ||
        !/^[A-Za-z0-9_.-]+$/u.test(segment),
    )
  ) {
    throw new Error("invalid repository path segment");
  }
  return segments.join("/");
}

/** Run a read-only Git command and return trimmed standard output. */
function runGit(args) {
  const result = spawnSync("git", args, {
    encoding: "utf8",
    stdio: ["ignore", "pipe", "pipe"],
  });
  if (result.error) {
    throw new Error(`unable to run git: ${result.error.message}`);
  }
  if (result.status !== 0) {
    throw new Error(`git ${args[0]} failed with exit code ${result.status}`);
  }
  return result.stdout.trim();
}

/** Discover the first configured Git remote that belongs to a provider. */
export function discoverRepository(provider, gitRunner = runGit) {
  const availableRemotes = new Set(
    gitRunner(["remote"]).split(/\r?\n/u).filter(Boolean),
  );
  const failures = [];

  for (const remote of provider.remoteCandidates) {
    if (!availableRemotes.has(remote)) {
      failures.push(`${JSON.stringify(remote)}: not found`);
      continue;
    }
    let remoteUrl;
    try {
      remoteUrl = gitRunner(["remote", "get-url", "--", remote]);
    } catch {
      failures.push(`${JSON.stringify(remote)}: could not read URL`);
      continue;
    }
    try {
      const repository = normalizeRepositoryUrl(remoteUrl, provider);
      return { remote, repository };
    } catch (error) {
      // Never echo the URL because it may contain embedded credentials.
      failures.push(`${JSON.stringify(remote)}: ${error.message}`);
    }
  }

  throw new Error(
    `No valid ${provider.label} remote found. Tried: ${failures.join("; ")}`,
  );
}

/** Build the API requests needed to synchronize one provider. */
export function buildRequests(provider, repository, description, topics) {
  if (provider.key === GITHUB.key) {
    const endpoint = `repos/${repository.split("/").map(encodeURIComponent).join("/")}`;
    return [
      {
        executable: GITHUB.executable,
        args: [
          "api",
          "--hostname",
          GITHUB.hostname,
          "--method",
          "PATCH",
          endpoint,
          "--input",
          "-",
        ],
        payload: { description },
      },
      {
        executable: GITHUB.executable,
        args: [
          "api",
          "--hostname",
          GITHUB.hostname,
          "--method",
          "PUT",
          `${endpoint}/topics`,
          "--input",
          "-",
        ],
        payload: { names: topics.slice(0, GITHUB_TOPIC_LIMIT) },
      },
    ];
  }

  return [
    {
      executable: GITLAB.executable,
      args: [
        "api",
        "--hostname",
        GITLAB.hostname,
        "--method",
        "PUT",
        `projects/${encodeURIComponent(repository)}`,
        "--input",
        "-",
        "--header",
        "Content-Type: application/json",
      ],
      payload: { description, topics },
    },
  ];
}

/** Execute one provider API request, or print it without executing in dry-run. */
export function runApiRequest(
  request,
  {
    dryRun = false,
    runner = spawnSync,
    output = console.log,
    executable = request.executable,
  } = {},
) {
  output(`+ ${request.executable} ${request.args.join(" ")}`);
  output(`  JSON: ${JSON.stringify(request.payload)}`);
  if (dryRun) {
    return;
  }

  const result = runner(executable, request.args, {
    input: JSON.stringify(request.payload),
    encoding: "utf8",
    stdio: ["pipe", "inherit", "inherit"],
  });
  if (result.error) {
    throw new Error(
      `unable to run ${request.executable}: ${result.error.message}`,
    );
  }
  if (result.status !== 0) {
    throw new Error(
      `${request.executable} failed with exit code ${result.status ?? "unknown"}`,
    );
  }
}

/** Parse the small command-line interface without adding a runtime dependency. */
export function parseArguments(argv) {
  const options = { dryRun: false, provider: "all", help: false };
  for (let index = 0; index < argv.length; index += 1) {
    const argument = argv[index];
    if (argument === "--dry-run") {
      options.dryRun = true;
    } else if (argument === "--help" || argument === "-h") {
      options.help = true;
    } else if (argument === "--provider") {
      index += 1;
      if (index >= argv.length) {
        throw new Error("--provider requires github, gitlab, or all");
      }
      options.provider = argv[index];
    } else if (argument.startsWith("--provider=")) {
      options.provider = argument.slice("--provider=".length);
    } else {
      throw new Error(`unknown argument: ${argument}`);
    }
  }
  if (!["github", "gitlab", "all"].includes(options.provider)) {
    throw new Error("--provider requires github, gitlab, or all");
  }
  return options;
}

/** Resolve a command to the first executable path reported by the platform. */
export function resolveCommand(
  command,
  { runner = spawnSync, platform = process.platform } = {},
) {
  const lookup = platform === "win32" ? "where.exe" : "which";
  const result = runner(lookup, [command], {
    encoding: "utf8",
    stdio: ["ignore", "pipe", "ignore"],
  });
  if (result.error || result.status !== 0) {
    return null;
  }
  return (
    result.stdout
      ?.split(/\r?\n/u)
      .map((candidate) => candidate.trim())
      .find(Boolean) ?? null
  );
}

function showHelp() {
  console.log(`Usage: node scripts/repository/src/sync_metadata.mjs [options]

Synchronize package.json description and provider-specific topics.

Options:
  --dry-run                    Show targets and payloads without API requests
  --provider github|gitlab|all Limit synchronization (default: all)
  -h, --help                   Show this help`);
}

/** Coordinate validation, target discovery, and metadata synchronization. */
export function main(argv = process.argv.slice(2)) {
  const options = parseArguments(argv);
  if (options.help) {
    showHelp();
    return;
  }

  const packageData = loadPackageJson();
  const description = getDescription(packageData);
  const providers =
    options.provider === "all"
      ? [GITHUB, GITLAB]
      : [PROVIDERS[options.provider]];

  // Resolve every selected target before the first possible remote mutation.
  const plans = providers.map((provider) => {
    const target = discoverRepository(provider);
    return {
      provider,
      ...target,
      topics: getTopics(packageData, provider.key),
    };
  });

  const resolvedExecutables = new Map();
  if (!options.dryRun) {
    const missing = [];
    for (const { provider } of plans) {
      const executable = resolveCommand(provider.executable);
      if (executable) {
        resolvedExecutables.set(provider.key, executable);
      } else {
        missing.push(provider.executable);
      }
    }
    if (missing.length > 0) {
      throw new Error(`Required CLI tools not found: ${missing.join(", ")}`);
    }
  }

  let completedApiRequests = 0;
  for (const plan of plans) {
    console.log(`\n${plan.provider.label}`);
    console.log(`Remote:      ${plan.remote}`);
    console.log(`Repository:  ${plan.repository}`);
    console.log(`Description: ${description}`);
    console.log(`Topics:      ${plan.topics.join(", ") || "(none)"}`);
    if (
      plan.provider.key === GITHUB.key &&
      plan.topics.length > GITHUB_TOPIC_LIMIT
    ) {
      console.log(
        `${GITHUB.label}: using only the first ${GITHUB_TOPIC_LIMIT} topics.`,
      );
    }
    for (const request of buildRequests(
      plan.provider,
      plan.repository,
      description,
      plan.topics,
    )) {
      try {
        runApiRequest(request, {
          dryRun: options.dryRun,
          executable:
            resolvedExecutables.get(plan.provider.key) ?? request.executable,
        });
        if (!options.dryRun) {
          completedApiRequests += 1;
        }
      } catch (error) {
        if (completedApiRequests > 0) {
          throw new Error(
            `${error.message}. ${completedApiRequests} earlier API request(s) succeeded and remain applied.`,
            { cause: error },
          );
        }
        throw error;
      }
    }
  }

  console.log(
    options.dryRun
      ? "\nDry run complete. No repository metadata was changed."
      : "\nRepository metadata synchronized successfully.",
  );
}

const invokedPath = process.argv[1]
  ? pathToFileURL(resolve(process.argv[1])).href
  : "";
if (invokedPath === import.meta.url) {
  try {
    main();
  } catch (error) {
    console.error(`ERROR: ${error.message}`);
    process.exitCode = 1;
  }
}
