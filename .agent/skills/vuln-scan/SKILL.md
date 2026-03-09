---
name: "vuln-scan"
description: "Run the repository security baseline: golangci-lint, go test, and govulncheck with a summarized reachable-vulnerability report. Use when asked to perform or repeat a vulnerability audit."
---

# Vulnerability Scan Skill

## Overview

Use this skill to run a repeatable repository security audit that combines:
- static checks (`golangci-lint`)
- runtime validation (`go test`)
- dependency and standard-library CVE analysis (`govulncheck`)

The bundled script writes raw logs and a compact summary to `.agent/reports/`.

## Inputs

- `repo`: repository path (default `.`)
- `out_dir`: output directory for logs/reports (default `<repo>/.agent/reports`)

## Quick Start

- `bash .agent/skills/vuln-scan/scripts/run_vuln_scan.sh .`
- `bash .agent/skills/vuln-scan/scripts/run_vuln_scan.sh . .agent/reports`

## Workflow

1. Ensure scanner binaries are available.
   - Installs `govulncheck` and `golangci-lint v2` if missing.
2. Run static checks.
   - Executes `golangci-lint run --timeout=4m ./...`.
3. Run tests.
   - Executes `go test -count=1 -p 1 ./...`.
4. Run CVE scan.
   - Executes `govulncheck -json ./...`.
5. Summarize findings.
   - Classifies vulnerabilities into `called`, `imported`, `required`.
   - Emits a concise text summary with reachable (`called`) vulnerabilities.

## Outputs

The script emits timestamped files:
- `golangci-lint-<ts>.log`
- `go-test-<ts>.log`
- `govulncheck-<ts>.json`
- `govuln-summary-<ts>.json`
- `vuln-summary-<ts>.txt`

## Exit Behavior

The script exits non-zero when any of the following is true:
- lint failed
- tests failed
- `govulncheck` command failed
- at least one `called` vulnerability was found
