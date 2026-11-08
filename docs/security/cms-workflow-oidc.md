# Workflow OIDC acquisition evidence

9 November 2026 (Asia/Jakarta).

Source/tests: Astro `workflow_oidc.go` and `workflow_oidc_test.go`. Synthetic local
HTTP fixtures exercise fixed bearer origin, exact audience/query preservation,
fourteen acquisition cases, two-attempt lifetime cap, endpoint/configuration
rejection and credential injection. No live token, Pages request or CMS mutation
is made. Returned token shape is not signature verification; private CMS RSA and
current-authority tests remain separate. Trusted configuration, TLS/system trust,
protected runner integrity and synchronized clocks remain assumptions. No hosted
proof, penetration test or vulnerability-free guarantee is claimed.

Standalone Go tests and vet passed; module statement coverage is 90.7% (not
full CMS coverage). Full `pnpm validate` passed formatting, lint, 53-file Astro
check with zero diagnostics, 25 unit tests, 23 repository tests and static build.
Native diagrams exported and rendered labels reviewed. Both checkpoint messages
passed the installed Custy validator. Gitleaks retained-history and changed-source
scans and whitespace checks passed. No race-detector or hosted proof is claimed.
No migration is needed (35); no database code changed. Production copies and
live credentials remain untouched. No staging, commit, tag or deployment occurred.
