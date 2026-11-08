# Deployment configuration security evidence

9 November 2026 (Asia/Jakarta).

Tests in Astro `deployment_settings_test.go` check exact context and deny substituted observations before credentials are read. No signed identity is inferred from environment values. Current CMS authentication and durable admission remain required. No API/SQL/public-output contract changed.

Verification: standalone `go test ./... -count=1 -cover` passed (3.980 seconds,
90.6% tool statement coverage); `go vet ./...` passed. Full `pnpm validate`
passed formatting, ESLint, 53-file type check with zero diagnostics, 25 unit tests,
23 repository tests and static build. Native diagram previews exported and reviewed.
No race-detector or hosted proof is claimed. No database implementation changed.
Both checkpoint messages passed Custy format validation; retained-history and
changed-source Gitleaks scans and whitespace checks passed. No commits, tags,
staging, deployment or credential acquisition was performed.
