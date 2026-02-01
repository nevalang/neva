# AGENTS.md

Follow these instructions.

## 1. 🤖 Operating Protocol

1. Use `context7` MCP server (when available) to fetch libraries API documentation.
2. For GitHub context, try `gh` (GitHub CLI) first; if unavailable or insufficient, fall back to `curl` (e.g., issues/PR descriptions/comments).
3. Run `golangci-lint` and `go test`. Fix warnings.
4. If uncertainty > 10%, ask user.
5. Update this file if changes to process, architecture, or rules.
6. Examples and parser for `.neva` changes. `go.mod` for Go imports. `docs/style_guide.md` for naming/formatting rules (check when writing `*.neva` code).
7. Plan -> Review -> Execute -> Review.
8. Refactor: Actively identify and resolve unnecessary complexity or duplication. Prioritize code clarity and long-term maintainability over chasing theoretical perfection.
9. Use targeted tests and cap long-running commands to ~5 minutes unless explicitly requested otherwise.

## 2. 📈 Self-Improvement Protocol

**After each session** (bug fix, feature, brainstorm), update this file (`AGENTS.md`) with:

- **Language semantics** learned (e.g., how connections work, port behavior).
- **Common patterns** discovered (e.g., typical error causes, debugging approaches).
- **Architecture insights** gained (e.g., how compiler phases interact).
- **Gotchas** encountered (e.g., edge cases, non-obvious behaviors).

**Balance**: Keep it concise. Every line must earn its place. Remove outdated info. Split into workflows/rules if sections grow large.

**Goal**: Build perfect context so future sessions start smarter.

### Session Notes (2026-01-23)

- **Language semantics**: Chained connections nest; receiver type constraints must use the sender from the same chain link (not the outer sender).
- **Language semantics**: Network senders now parse any `constLit`; only primitives/union literals get a type expr without additional analyzer inference.
- **Language semantics**: Union tag/type compatibility relies on structural union checks via normal sender/receiver subtype validation.
- **Common patterns**: Overload/generic resolution relies on `deriveNodeConstraintsFromNetwork`; nested chains need explicit sender-to-receiver pairing to avoid union/any leakage.
- **Architecture insights**: Switch case output type inference is used in analyzer (before desugaring) for both overload resolution and network type checks.
- **Architecture insights**: Union node tag/data validation is driven by `buildUnionTagInfos` + `validateUnionDataReceivers` to keep Union<T> wiring consistent.
- **Gotchas**: Portless outports require explicit port names when multiple outports exist; generic std nodes like `fmt.Println` need explicit type args if type inference is unavailable.

### Session Notes (2026-01-26)

- **Language semantics**: Union:tag currently assumes a single sender; fan-in is explicitly rejected to keep wiring semantics unambiguous.
- **Common patterns**: Avoid extra union-type equivalence checks at tag-binding time; rely on standard subtype validation for compatibility.
- **Gotchas**: In sender position, `[...]` is parsed as fan-in (multiple senders), so list literals should not be used as senders without explicit disambiguation.
- **Gotchas**: Literal senders are limited to primitives/union literals; list/dict/struct literals are rejected for now (use const refs).

### Session Notes (2026-01-27)

- **Language semantics**: Literal senders must be in a chain triggered by a signal (e.g. `:start -> true -> ...`).
- **Language semantics**: Fan-out (`-> [...]`) can contain chained connections, but the chain head is the sender before the fan-out (e.g. `:start -> [foo, bar -> baz]`).
- **Language semantics**: `std/core.New<T>` now requires a `sig` inport; explicit uses must be triggered (e.g. `:start -> new:sig`).
- **Common patterns**: Stdlib Switch setups should seed case values via `:start` to avoid autonomous constant senders.
- **Architecture insights**: `analyzeSender` rejects const senders when `prevChainLink` is empty (analyzer enforces the rule).
- **Gotchas**: Stdlib modules can fail compilation after analyzer rule changes; update `std/*` constants accordingly.

### Session Notes (2026-01-28)

- **Common patterns**: Enabling `govet`'s `fieldalignment` in golangci-lint surfaces many existing struct-order warnings; plan for a staged rollout (only-new-issues) or bulk field reordering.
- **Common patterns**: `betteralign` is best run as a separate make target with `-fix` since it cannot integrate into golangci-lint.
- **Gotchas**: `Select` waits for all `then` array inputs on each `if`; if any `then` slot is only sent conditionally it can deadlock.
- **Common patterns**: golangci-lint v2.5.0 no longer exposes `gosimple`, `stylecheck`, or `tenv` as standalone linters; rely on `staticcheck` instead.
- **Gotchas**: `nolintlint` flags stale `//nolint:golint`; replace with a real doc comment or an active linter name.
- **Common patterns**: `gosec` G115 still flags checked conversions; use a guard plus `// #nosec G115` on the cast line.
- **Common patterns**: `gosec` G306 expects 0600 perms; for readable artifacts, keep 0644 and add `// #nosec G306` with rationale.
- **Architecture insights**: Runtime cancellation now flows via a typed context key helper to avoid raw string keys.
- **Common patterns**: `gocritic`'s commentedOutCode flags comments that resemble code; remove them instead of paraphrasing.
- **Common patterns**: `golangci-lint` thelper prefers `t.Helper()` inside helper constructors used in tests.
- **Common patterns**: `govet` fieldalignment warnings can be suppressed with `//nolint:govet` on struct or inline struct literal lines when reordering is too invasive.
- **Common patterns**: `gocyclo` can be temporarily suppressed with `//nolint:gocyclo` and a short rationale, but consider later refactor.

### Session Notes (2026-01-29)

- **Common patterns**: `copyloopvar` flags redundant loop-var shadowing (`name := name`) in Go 1.22; remove the copy unless needed for older versions.
- **Gotchas**: Runtime now panics on invalid array/case indices (negative or uint8 overflow) to surface compiler bugs early.
- **Common patterns**: Resolve union elements in deterministic key order to avoid nondeterministic error paths in tests.
- **Common patterns**: `golangci-lint` config uses a strict allowlist with targeted test/go:generate exclusions; tune `tparallel`/`gocyclo` per package if flakiness or complexity hotspots appear.
- **Common patterns**: Consider opt-in linters like `depguard`, `gofumpt`, `goimports`, `gochecknoglobals`, `gochecknoinits`, `wrapcheck`, and `forbidigo` based on team appetite; add them gradually to avoid noisy rollouts.
- **Common patterns**: `govet` fieldalignment can be satisfied by moving pointer-heavy fields before strings to reduce pointer bytes and GC scan range.
- **Common patterns**: `depguard` can enforce stdlib-only boundaries with separate rules for `internal/runtime/*.go` and `internal/runtime/funcs/**` when funcs may import `internal/runtime`.
- **Gotchas**: `gochecknoglobals` will flag template strings, expected outputs in tests, and shared helpers; plan exclusions or refactors before enabling widely.
- **Common patterns**: To satisfy `gochecknoglobals`, switch static template vars to `const`, and move shared values into helper funcs; use per-program counters instead of package globals.
- **Common patterns**: If performance trumps `gochecknoglobals`, allow a scoped global counter with `//nolint:gochecknoglobals` rather than threading it through every outport.

### Session Notes (2026-01-31)

- **Language semantics**: Import blocks now allow inline comments; struct literals accept trailing commas.
- **Language semantics**: Const/literal senders are only valid when chained from a triggering signal.
- **Language semantics**: `Main` must have exactly one `start` inport and one `stop` outport, both `any` and non-array.
- **Common patterns**: Use explicit stdlib components (`operators`, `routers.Switch`, `streams.Range`, `Union`) instead of composite sender syntax.
- **Architecture insights**: Debug runtime validation code is generated by the Go backend only when the CLI flag is enabled.
- **Gotchas**: Dotenv loaders are signal-style and do not return dictionaries; use `os.Environ` to read environment values.

### Session Notes (2026-02-01)

- **Language semantics**: Array-bypass now uses `[*]` on both sides of a normal connection; index `255` is reserved for the wildcard.
- **Architecture insights**: Array-bypass is handled as a normal connection with a sentinel slot index (`ArrayBypassIdx`) instead of a dedicated AST connection type.
- **Gotchas**: Using `[255]` directly is now a parser error; always use `[*]` for array-bypass.
- **Common patterns**: Stream/list conversions use `streams.FromList<T>` and `lists.FromStream<T>`; `StreamToList`/`ListToStream` were removed.
- **Architecture insights**: Stream helpers live under `std/streams` (map/filter/reduce in dedicated files; `ForEach` in `std/streams/for_each.neva`).
- **Gotchas**: `std/time` only provides delays/durations; there is no wall-clock date/time source, so pass timestamps/durations into components that need "today."
- **Common patterns**: `lists.FromArray<T>` builds a list by chaining `streams.FromArray<T>` into `lists.FromStream<T>`.

## 3. ⚡ Core Concepts

- **Dataflow**: Programs are graphs. Nodes process data; edges transport it.
- **Implicit Parallelism**: Every node runs in parallel.
- **Type System**: Static, structural, with generics and tagged-unions.
- **Visibility**: Entities are private by default. Export with `pub`.
- **Entities**:
  - **Components**: Logic containers, node blueprints.
  - **Interfaces**: Port definitions, abstract components.
  - **Types**: Type definitions.
  - **Constants**: Fixed values.
- **Program Hierarchy**:
  - **Module**: Root unit (has `neva.yml`).
  - **Package**: Directory with `*.neva` files.
  - **Component**: The building block.

## 4. 🧠 Architecture

### Compiler (`internal/compiler/`)

1. **Parser**: ANTLR4 -> Raw AST.
2. **Analyzer**: Semantics & Type checking.
3. **Desugaring**: Syntactic sugar -> Canonical AST.
4. **IR Gen**: Canonical AST -> IR.
5. **Backend**: IR -> Target Code.

### Runtime (`internal/runtime/`)

The runtime is a library embedded into every compiled program.

- **FuncRunner**: Executes "native components" (runtime functions).
- **Connector**: Manages message passing.
- **Extensibility**:
  - **Native Functions**: `internal/runtime/funcs`.
  - **Func Registry**: `internal/runtime/funcs/registry.go`.
  - **Func Interface**: `runtime.Func` & `FuncCreator`.

### Standard Library (`std/`)

The standard library provides components for all programs. Some are implemented in Neva, some use runtime functions written in Go via `#extern`.

## 5. 🛠️ Debugging Tips

**Debug Compiler Output**:

- **IR**: `neva run --target ir <pkg>`
- **Trace**: `neva run --emit-trace <pkg>`
- **Runtime validation**: `neva run --debug-runtime-validation <pkg>` or `neva build --debug-runtime-validation <pkg>` (compiler-only check that prints unconnected senders/receivers to validate runtime wiring; intended for language developers inspecting compiler output)

**Debug the CLI/Compiler**:

- **Logs**: Use `fmt.Printf`, remove before finishing.
- **Tests**: `go test ./...`

## 6. 🗺️ Key Locations

- `cmd/neva/`: CLI Entry point.
- `internal/cli/`: CLI implementation.
- `internal/compiler/parser/neva.g4`: Grammar Definition.
- `e2e/`: End-to-End Tests.
- `examples/`: Example programs.
- `pkg/`: Shared utilities.

## 7. 🎨 Coding Standards

- **Go Idioms**:
  - Comments: Every function should have a short doc comment. If it relates to Neva semantics, include a tiny Neva example when helpful.
  - Use `any` instead of `interface{}`.
  - TD tests: `tests := []struct{ name string ... }`
  - Test case names: `lower_snake_case`
  - KISS: simpler code > complex abstractions
  - Utils: `pkg/` for shared utils (EXCEPT `runtime`)
    - If duplicated in 3+ places, move it to `pkg/` (except `runtime`).

**Workflow**:

1. `make build` (Verify compilation).
2. `golangci-lint run ./...` then `go test ./...` (Verify lint + tests).
3. `make antlr` (Regenerate parser if `.g4` changed).
