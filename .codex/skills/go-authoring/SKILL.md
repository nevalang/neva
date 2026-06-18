---
name: "go-authoring"
description: "Use when editing, reviewing, or generating Go files in this repository. Covers Neva-specific Go style, runtime boundaries, lint expectations, and validation."
---

# Go Authoring

Use this skill for changes touching `*.go`.

## Rules

- Follow the nearest `AGENTS.md` as scoped source of truth.
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
- Never ignore errors. If an error cannot be returned or sent through the
  component API, handle it explicitly and document or panic only when it
  represents an internal invariant violation.

## Runtime Boundaries

- `internal/runtime/*.go`: stdlib imports only.
- `internal/runtime/funcs/**`: stdlib + `internal/runtime`.
- Binary runtime funcs in `internal/runtime/funcs/**` must receive `left` and
  `right` inputs concurrently; do not add sequential two-input receives in
  operator helpers or creators.

## Comments

- Add doc comments for new Go functions and types.
- For newly generated Go code blocks longer than 3 lines, add a short one-line
  intent comment when the purpose is not immediately obvious.
