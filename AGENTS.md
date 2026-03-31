# AGENTS.md

This file is a compact operating guide for coding agents in this repository.
It is intentionally short and stable. Use linked docs for deep details.

## 1) Operating Protocol

1. Plan -> Review -> Execute -> Review.
2. Use `context7` MCP (when available) for library API docs.
3. For GitHub context, use `gh` first; fall back to `curl` only if needed.
4. If uncertainty is >10% and local context cannot resolve it safely, ask user.
5. Refactor proactively when it clearly improves clarity/maintainability.
6. Run targeted checks with ~5m cap unless user requests longer runs.
7. Run `golangci-lint` and `go test` for touched scope; fix issues you introduce.
8. For PR comment tasks: apply changes first, then reply to each addressed review comment via `gh`; do not resolve threads unless user asks.
9. For generated tests, include short intent comments.
10. Keep this file updated when process/architecture/rules change.
11. For repository-local skills, prefer concise English `SKILL.md` guidance (tool list + workflow); avoid bundled scripts unless explicitly requested.
12. Error semantics policy: return `*compiler.Error` for invalid user programs (input/domain failures), but `panic` on internal invariant violations or impossible cross-stage states.

## 2) High-Level Project Context

- Neva is a statically typed, compiled **dataflow** language.
- Programs are explicit node/edge graphs; implicit parallelism is default.
- Compiler pipeline: parser -> analyzer -> desugarer -> IR gen -> backend.
- Runtime executes generated Go with message passing primitives.
- Standard library mixes pure Neva components and `#extern` runtime-backed components.
- LSP source is externalized (`nevalang/neva-lsp`), not under `cmd/lsp` in this repo.

## 3) Documentation Router

- Project overview and quick start: [README.md](./README.md)
- Architecture map: [ARCHITECTURE.md](./ARCHITECTURE.md)
- Contributor workflow and release basics: [CONTRIBUTING.md](./CONTRIBUTING.md)
- Documentation index: [docs/README.md](./docs/README.md)
- Language rationale and design decisions: [docs/qna.md](./docs/qna.md)
- Neva style rules: [docs/style_guide.md](./docs/style_guide.md)
- Language deep dive: [docs/book/README.md](./docs/book/README.md)
- Build/test/dev commands: [Makefile](./Makefile)
- Go lint policy (source of truth): [.golangci.yml](./.golangci.yml)

## 4) Stable Language/Architecture Invariants

These are distilled from long-term session notes and should be treated as
high-signal constraints:

1. Preserve 1:1 mapping between text and graph. Prefer explicit stdlib nodes
   over hidden control-flow sugar.
2. `Main` must have exactly one `start` inport and one `stop` outport
   (both `any`, non-array).
3. Literal/const senders are constrained:
   - valid only in signal-triggered chains (e.g. `:start -> ...`)
   - primitive/union literals are allowed; bytes literal is not
4. Array-bypass uses `[*]` on both sides of a single connection.
   - index `255` is reserved for wildcard (`[*]`)
   - do not mix bypass and indexed usage on same port
5. Error flow convention is stable: `res` for primary result, `err error`
   for failures, `?` for propagation when custom handling is unnecessary.
6. Conversions are explicit stdlib components (no implicit cast syntax).
   - policy split: builtin (total casts) vs `std/strconv` (parse/format)
7. `bytes` is transport-focused and currently has no literal syntax.
   Use explicit converters (`bytes.FromString`, `strings.FromBytes`).
8. Keep compiler-contract stdlib semantics explicit (`builtin`, `Union`,
   `Struct`, `#autoports`, desugaring-sensitive behavior).

## 5) Code Style

### Go (do not invent style; follow repository policy)

- Source of truth: `.golangci.yml`, `go.mod`, and existing code patterns.
- Treat `go.mod` as the version ceiling for Go features and stdlib APIs.
- Prefer modern Go idioms available in the target version over legacy
  patterns. By default, favor `any`, `errors.AsType[T]`, `new(val)`,
  `wg.Go`, `b.Loop`, `t.Context()`, `slices`/`maps`/`cmp`, `min`/`max`,
  `SplitSeq`/`FieldsSeq` for iteration-only use, and `omitzero` when
  zero-value JSON semantics matter.
- Always run `gofmt` on changed Go files.
- Respect active lints (notable strict checks include `depguard`,
  `fieldalignment`, `nil*`, `gosec`, `gochecknoglobals`, `gochecknoinits`).
- Runtime boundary rule is strict:
  - `internal/runtime/*.go`: stdlib imports only
  - `internal/runtime/funcs/**`: stdlib + `internal/runtime`
- Prefer the repo-pinned linter command when local binary is stale:
  `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 run ./...`
- Add doc comments for new Go functions/types (including unexported helpers).
- For newly generated Go code blocks longer than 3 lines, add a short one-line intent comment before the block when the purpose is not immediately obvious.

### Neva (`*.neva`)

- Source of truth: `docs/style_guide.md`.
- Mandatory basics:
  - tabs for indentation
  - node aliases in `lower_snake_case`
  - import ordering/grouping per style guide
  - omit explicit port names when unambiguous

## 6) Validation Workflow

Prefer smallest meaningful scope first, then widen only if needed.

1. Lint changed scope (`golangci-lint`).
2. Run targeted tests for touched packages/components.
3. For parser/grammar work, run parser smoke tests and regenerate parser.
4. Run broader `go test ./...` only when appropriate; it can be long due to
   `examples/` + `e2e/`.

Useful commands:

```bash
golangci-lint run ./...
go test ./...
make antlr
make gofix
```

## 7) Change-Specific Checklists

### Grammar / parser / analyzer changes

- Update grammar/parser artifacts when needed (`make antlr`).
- Re-run parser smoke tests and touched analyzer tests.
- Preserve documented core semantics (array-bypass, sender rules, `Main` rules).

### Stdlib / runtime extern changes

- Keep `std/* #extern(...)` names synchronized with
  `internal/runtime/funcs/registry.go`.
- If stdlib naming/conventions change, update docs in same PR:
  `docs/style_guide.md` and `docs/qna.md`.
- If compiler bootstrap utils are affected, regenerate from repo CLI
  (avoid stale global binaries).

### Tests and e2e

- `examples/` is one module (all examples must compile together).
- `e2e/` contains separate modules with Go harness tests.
- Keep e2e runs time-bounded; avoid orphaned long-running subprocess chains.

## 8) Keep AGENTS.md Lean

- Do not append per-session journals here.
- Keep only durable process rules, architecture deltas, and high-value gotchas.
- If a note is transient, put it in issue/PR discussion instead.
