---
name: "vuln-scan"
description: "Run a repository vulnerability audit with lint, tests, and govulncheck. Use this when asked for security scanning or CVE triage in this repo."
---

# Vulnerability Scan Skill

## Overview

Use this skill to run a lightweight, repeatable security audit without hardcoded scripts.
Prefer concise command execution and summarize results with focus on reachable vulnerabilities.

## Core Tools

- `golangci-lint` (v2)
- `go test`
- `govulncheck`
- `jq` (optional, for JSON summarization)
- `gh` (when the scan is part of a PR/review workflow)

## Recommended Workflow

1. Verify baseline quality checks.
   - `golangci-lint run --timeout=4m`
   - `go test -count=1 -p 1 ./...`
2. Run vulnerability scan.
   - `govulncheck ./...`
3. If detailed triage is needed, switch to JSON mode.
   - `govulncheck -json ./... > /tmp/govulncheck.json`
4. Classify findings by impact.
   - prioritize reachable (`called`) findings
   - then `imported` and `required`
5. Apply minimal safe fixes.
   - prefer patch-level toolchain updates for stdlib CVEs
   - bump vulnerable modules to fixed versions
6. Re-run lint/tests/govulncheck and report before/after.

## Reporting Guidance

Keep the report short and actionable:
- findings count by class (`called/imported/required`)
- affected modules and fixed versions
- exact validation commands and outcomes
- residual risks (if any)
