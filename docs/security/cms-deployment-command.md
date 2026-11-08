# Deployment command security evidence

9 November 2026 (Asia/Jakarta).

Astro `deployment_command.go` / `deployment_command_test.go` test nineteen local
file/argument/policy/runtime cases: success, disabled mode, unknown/duplicate/missing
flags, pin mismatch, unknown/duplicate policy JSON, scope/settings rejection,
missing/directory/empty/oversize files, corrupted archives, runner failure, invalid
report and output failure. Default main routing refuses before any credential read.
Opened-file identity/size guards are implemented; leaf-link rejection is not claimed
as tested on a platform where creating symlinks requires special permission.
No hostile-parent-directory isolation, signed policy, hosted execution or penetration
test is claimed. Policy+pin are trusted operator input, never dispatch authority.
No secret is echoed or persisted. CLI failure does not authorize a create retry.

Final `go test ./... -count=1 -cover` passed in 4.470 seconds at 90.2% module
statement coverage; `go vet ./...` passed. Full `pnpm validate` passed formatting,
ESLint, 53-file Astro checks (zero diagnostics), 25 unit tests, 23 repository tests
and static build. An actual no-input `go run . deploy-local` returned the constant
receiver error with no stdout or network request. No positive command was run with
live credentials. Native diagram previews exported and reviewed. Both checkpoint
messages passed Custy validation. Gitleaks retained-history and changed-source
scans and whitespace checks passed. No race-detector result is claimed.
No database code changed; schema remains 35. No staging, commit, tag, hosted
deployment, credential change or editing database mutation was performed.
