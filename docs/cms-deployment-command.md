# Gated local deployment command

Updated 9 November 2026 (Asia/Jakarta).

The Astro Go tool now routes `deploy-local` to the guarded runtime. It remains
disabled by default. No hosted workflow, permission, live credentials, provider
setting, CMS route registration, database or production copy was changed.

## Exact inputs and trust

The command accepts exactly five flag/value pairs, in any order: `--policy`,
`--policy-sha256`, `--request`, `--candidate`, `--uploaded`. Unknown, repeated,
missing or ambiguous flags are rejected. All four files are local regular files.
There is no download URL, extraction, shell interpolation or automatic retry.
Limits are policy 8192 bytes, request 12288 bytes, candidate 8 MiB and uploaded
wrapper the existing `MaxDeployWrapper` bound. Empty/oversize/non-file/leaf-link
inputs fail. Opened-file identity and size are checked before bounded reading.
Parent directories and concurrent writers are trusted operator workspace inputs;
this does not claim adversarial filesystem isolation or race-free directory traversal.

The raw policy file must match the separately supplied lowercase SHA-256. This is
an integrity pin, not a signature or proof of who reviewed the policy. Obtain the
pin independently from reviewed protected configuration. Never calculate both the
policy and trusted pin from untrusted dispatch data or a public manifest. A caller
who can replace both policy and pin controls this operator boundary. The command
does not elevate that caller or protect already-compromised privileged workflow code.

Policy JSON has exactly these fields (no secrets):

```json
{
  "schema_version": 1,
  "scope": "astro-protected-runtime-policy-v1",
  "repository": "example/homepage",
  "repository_id": "123",
  "workflow_ref": "example/homepage/.github/workflows/cms-pages-deploy.yml@refs/heads/main",
  "head_sha": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "receiver_origin": "https://cms.example.test",
  "audience": "https://github.com/example",
  "issuer_origin": "https://runner.actions.githubusercontent.com"
}
```

These are placeholders, not working production configuration. Duplicate/unknown
JSON fields, a wrong schema/scope and missing origin are rejected. Full workflow
settings are checked before archive files are loaded; complete archive verification
still occurs in the runtime before credential access. No policy URL or expected
content digest comes from the deployment request. The policy contains no token.

An operator invocation, only after configuration review, has this form:

```text
go run . deploy-local --policy reviewed-policy.json --policy-sha256 <independently-reviewed-sha256> --request deployment.json --candidate candidate.zip --uploaded uploaded.zip
```

Run from `scripts/pages-tool`. A normal terminal with deployment opt-in absent
returns a constant failure before reading files. Setting the flag alone is not
enough: exact workflow observations, code/run identity, reviewed origins and valid
live credentials/issuer identity and current CMS admission are still required.
Do not run this command against hosted services as a verification shortcut.

## Outcomes and recovery

Success prints only a constant statement that the exact operation was reported;
it does not claim public deployment completion. Failures print a constant receiver
diagnostic; raw provider bodies, credentials, file paths and policy are never echoed.
Output failure never repeats the operation. There is no acknowledgement output
file or persistent recovery store. If reporting fails, CLI memory is discarded;
do not rerun deployment to recover it. Keep the hold and use a separately designed
authenticated reconciliation/report-only path. A new process cannot reset permanent
CMS admission. Timeout cannot undo a provider request already accepted remotely.

Tests use bounded real local files and a package-local runtime fixture; they never
acquire live tokens or deploy. Existing joined runtime tests cover real loopback
issuer/receiver/provider clients separately. No hosted or joined RSA/SQL proof is
claimed. No public writer/archive contract change or migration (CMS schema 35).

Next: protected workflow review and exact artifact acquisition, plus authenticated
reconciliation/completion and ordering-hold release before hosted rollout.
