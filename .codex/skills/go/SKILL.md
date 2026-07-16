---
name: "go"
description: "Use for Go changes in Neva: authoring, refactoring, debugging, or review."
---

# Go

Use this skill for changes touching `*.go`. Read the matching developer
documentation named by `AGENTS.md` before editing.

## Rules

- Follow idiomatic Go as described by [Effective Go](https://go.dev/doc/effective_go),
  the Go documentation, and established repository patterns.
- Source of truth for style and APIs: `.golangci.yml`, `go.mod`, and existing
  code patterns.
- Treat `go.mod` as the version ceiling for Go features and stdlib APIs.
- Prefer modern Go idioms available in the target version over legacy patterns.
- Return early to keep control flow flat. Prefer straightforward code and
  standard-library primitives over clever abstractions.
- Always run `gofmt` on changed Go files.
- Respect active lints.
- Never add file-wide or package-wide `nolint` directives. If suppression is
  unavoidable, keep it narrowly scoped to the exact line/rule and add a short
  justification.
- Prefer `Makefile` targets for standard checks when applicable.
- Do not add helper functions or methods unless they abstract at least two
  meaningful operations or remove real complexity.
- Never ignore errors, including assignment to `_`. Either return an error
  (wrapping when useful) or handle it locally; do not report an error and then
  continue as though it were handled.
- Add doc comments for new exported Go functions and types.

## Runtime

- For `internal/runtime/funcs/**`, follow
  `docs/developer/runtime-functions.md`.
