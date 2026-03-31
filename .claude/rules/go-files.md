---
paths:
  - "**/*.go"
---

# Go Files

- Follow the nearest `AGENTS.md` as the scoped source of truth.
- Source of truth for style and APIs: `.golangci.yml`, `go.mod`, and existing code patterns.
- Treat `go.mod` as the version ceiling for Go features and stdlib APIs.
- Prefer modern Go idioms available in the target version over legacy patterns.
- Always run `gofmt` on changed Go files.
- Respect active lints.
- Prefer `Makefile` targets for standard checks when applicable.
- Runtime boundary rule is strict:
  - `internal/runtime/*.go`: stdlib imports only
  - `internal/runtime/funcs/**`: stdlib + `internal/runtime`
- Add doc comments for new Go functions and types.
- For newly generated Go code blocks longer than 3 lines, add a short one-line intent comment when the purpose is not immediately obvious.
