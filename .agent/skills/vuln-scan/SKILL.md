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

This skill requires `golangci-lint`, `govulncheck`, and `jq` to be available in PATH.

## Quick Start

```bash
# Install required tools
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run scan
golangci-lint run --timeout=4m ./...
go test -count=1 -p 1 ./...
govulncheck -json ./... | tee govulncheck.json

# Summarize findings with jq
jq -s '
  def finding_class($f):
    if ([$f.trace[]?.module] | index("github.com/nevalang/neva")) != null then "called"
    elif ([$f.trace[]? | select(has("package"))] | length) > 0 then "imported"
    else "required"
    end;
  (map(select(has("osv")) | .osv) | INDEX(.id)) as $osvmap |
  (map(select(has("finding")) | .finding) | group_by(.osv) |
    map({
      id: .[0].osv,
      class: (if (map(finding_class(.)) | index("called")) != null then "called" elif (map(finding_class(.)) | index("imported")) != null then "imported" else "required" end),
      fixed_version: .[0].fixed_version,
      found_modules: ([.[].trace[0].module] | unique),
      summary: ($osvmap[.[0].osv].summary // "")
    })
  )
' govulncheck.json > govuln-summary.json

# Report reachable vulnerabilities
jq -r '.[] | select(.class == "called") | "- \(.id) | fixed: \(.fixed_version) | module(s): \(.found_modules | join(",")) | \(.summary)"' govuln-summary.json
```

## Required Tools

- `golangci-lint` v2+
- `govulncheck` (x/vuln)
- `jq` (for JSON processing)

## Workflow

1. Run static checks:
   ```bash
   golangci-lint run --timeout=4m ./...
   ```

2. Run tests:
   ```bash
   go test -count=1 -p 1 ./...
   ```

3. Run CVE scan:
   ```bash
   govulncheck -json ./...
   ```

4. Process results:
   - Classify vulnerabilities into `called`, `imported`, `required`
   - Extract reachable (`called`) vulnerabilities
   - Generate human-readable summary

## Exit Behavior

Exit non-zero when:
- lint fails
- tests fail
- `govulncheck` command fails
- at least one `called` vulnerability is found
