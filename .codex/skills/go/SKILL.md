---
name: "go"
description: "Use for Go changes in Neva: authoring, refactoring, debugging, or review."
---

# Go

Use this skill for changes touching `*.go`. Read the matching developer
documentation named by `AGENTS.md` before editing.

## Rules

- Source of truth for style and APIs: `.golangci.yml`, `go.mod`, and existing
  code patterns.
- Treat `go.mod` as the version ceiling for Go features and stdlib APIs.
- Prefer modern Go idioms available in the target version over legacy patterns.
- Always run `gofmt` on changed Go files.
- Respect active lints.
- Never add file-wide or package-wide `nolint` directives. If suppression is
  unavoidable, keep it narrowly scoped to the exact line/rule and add a short
  justification.
- Prefer `Makefile` targets for standard checks when applicable.
- Do not add helper functions or methods unless they abstract at least two
  meaningful operations or remove real complexity.
- Never ignore errors, including assignment to `_`. Handle the error locally or
  pass it to the caller; choose one outcome such as return (wrapping when useful),
  log, or panic, but do not both report and continue with the same error.

## Runtime

- For `internal/runtime/funcs/**`, follow
  `docs/developer/runtime-functions.md`.
- Add doc comments for new Go functions and types.
