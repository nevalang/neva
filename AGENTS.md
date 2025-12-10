# AGENTS.md

> [!NOTE]
> This file is your primary source of truth for working on the Neva repository. STRICTLY FOLLOW these instructions.

## 1. 🤖 Agent Operating Protocol

1.  **Context7 MCP**: ALWAYS use the `context7` MCP server to avoid hallucinations.
2.  **Linter**: Always run `make test` (includes `golangci-lint`). **FIX ALL WARNINGS**.
3.  **Uncertainty Check**: If uncertainty > 10% (0.1), **ASK** the user.
4.  **Self-Correction**: IF you change the build process, architecture, or strict rules, **YOU MUST UPDATE THIS FILE** (`AGENTS.md`) to reflect the new reality.
5.  **Cross-Reference**:
    - **Editing `.neva` code**: Read `examples/` and `internal/compiler/parser/neva.g4` (grammar) first.
    - **Editing Go code**: Read `go.mod` before importing.
6.  **Step-by-Step**: Plan -> Review -> Execute one step -> Review.

## 2. ⚡ The Neva Language (Internal View)

**Core Concepts**:

- **Dataflow**: Programs are graphs. Nodes process data; edges transport it.
- **Implicit Parallelism**: Every node runs in parallel (goroutines).
- **Type System**: Static, structural (like TypeScript), with generics and Rust-like tagged-unions.
- **Visibility**: Entities are private by default. Use `pub` keyword to export (e.g., `pub def Main`).
- **Entities**:
  1.  **Components**: Logic containers (Interface + Implementation), node blueprints.
  2.  **Interfaces**: Port definitions (Inports/Outports), abstract components.
  3.  **Types**: Type definitions (Structs, Unions, etc.).
  4.  **Constants**: Fixed values (`const msg string = 'hello'`).

**Program Hierarchy**:

- **Module**: Root unit (has `neva.yml`).
- **Package**: Directory with `*.neva` files.
- **Component**: The building block.

## 3. 🧠 Architecture & Runtime

### Compiler Pipeline (`internal/compiler/`)

1.  **Parser** (`parser/`): ANTLR4 (`neva.g4`) -> Raw AST.
2.  **Analyzer** (`analyzer/`): Semantic analysis & Type checking.
3.  **Desugarer** (`desugarer/`): Lowers syntactic sugar -> **Desugared AST** (Canonical).
4.  **IR Gen** (`irgen/`): Desugared AST -> **IR** (Intermediate Representation).
5.  **Backend** (`backend/`): IR -> Target Code (Go/Native/WASM).

### Runtime Model (`internal/runtime/`)

The runtime (`internal/runtime`) is a library embedded into every compiled program.

- **FuncRunner**: Executes "Native" flows (Go functions).
- **Connector**: Manages message passing.
- **Extensibility**:
  - **Native Functions**: Implemented in `internal/runtime/funcs`.
  - **Registry**: Register new Go funcs in `internal/runtime/funcs/registry.go`.
  - **Interface**: All runtime funcs implement `runtime.Func` (and `FuncCreator`).

## 4. 🛠️ Debugging & Observability

**Debug Compiler Output**:

- **Emit IR**: `neva build --target ir --target-ir-format yaml <pkg>` (View the compiler's intermediate state).
- **Trace Execution**: `neva build --emit-trace <pkg>` -> Generates `trace.json` (visualize message flow).

**Debug the CLI/Compiler**:

- **Logs**: Use `fmt.Printf` for debugging, but **MUST REMOVE** before finishing.
- **Test**: `go test -v ./...` for verbose output.

## 5. 🗺️ Key Locations

- `cmd/neva/`: CLI Entry point.
- `internal/compiler/parser/neva.g4`: **Grammar Definition**.
- `internal/runtime/funcs/`: **Standard Library Implementation** (Go side).
- `std/`: **Standard Library Definition** (Neva side, `#extern`).
- `e2e/`: **End-to-End Tests**.

## 6. 🎨 Coding Standards

**Go Idioms**:

- **Modern Go**: Use `any` instead of `interface{}`.
- **Table-Driven Tests**: `tests := []struct { name string ... }`
- **Naming**: Test cases MUST be `lower_snake_case`.
- **KISS**: Simple code > Complex abstractions.
- **Utils**: Use `pkg/` for shared utilities (EXCEPT in `internal/runtime` -> keep zero deps).
  - If functionality is duplicated in 3+ places, consider moving it to `pkg/` (except `runtime`).

**Workflow**:

1.  `make build` (Verify compilation).
2.  `make test` (Verify logic & linting).
3.  `make antlr` (Regenerate parser if `.g4` changed).
