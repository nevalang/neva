# AGENTS.md

Follow these instructions.

## 1. 🤖 Operating Protocol

1. Use `context7` MCP server (when available) to fetch libraries API documentation.
2. Run `golangci-lint` and `go test`. Fix warnings.
3. If uncertainty > 10%, ask user.
4. Update this file if changes to process, architecture, or rules.
5. Examples and parser for `.neva` changes. `go.mod` for Go imports. `docs/style_guide.md` for naming/formatting rules (check when writing `*.neva` code).
6. Plan -> Review -> Execute -> Review.
7. Refactor: Actively identify and resolve unnecessary complexity or duplication. Prioritize code clarity and long-term maintainability over chasing theoretical perfection.
8. Use targeted tests and cap long-running commands to ~5 minutes unless explicitly requested otherwise.

## 2. 📈 Self-Improvement Protocol

**After each session** (bug fix, feature, brainstorm), update this file (`AGENTS.md`) with:

- **Language semantics** learned (e.g., how connections work, port behavior).
- **Common patterns** discovered (e.g., typical error causes, debugging approaches).
- **Architecture insights** gained (e.g., how compiler phases interact).
- **Gotchas** encountered (e.g., edge cases, non-obvious behaviors).

**Balance**: Keep it concise. Every line must earn its place. Remove outdated info. Split into workflows/rules if sections grow large.

**Goal**: Build perfect context so future sessions start smarter.

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
  - Use `any` instead of `interface{}`.
  - TD tests: `tests := []struct{ name string ... }`
  - Test case names: lower_snake_case
  - KISS: simpler code > complex abstractions
  - Utils: `pkg/` for shared utils (EXCEPT `runtime`)
    - If duplicated in 3+ places, move it to `pkg/` (except `runtime`).

**Workflow**:

1. `make build` (Verify compilation).
2. `golangci-lint run ./...` then `go test ./...` (Verify lint + tests).
3. `make antlr` (Regenerate parser if `.g4` changed).
