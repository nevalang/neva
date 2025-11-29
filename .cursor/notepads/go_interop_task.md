# Go → Neva Package Mode Design

## quick example first (nevalang)

before details, here is a small nevalang package with a multi-port export. this is what the go package will target.

```
neva.yml
userfmt/
  userfmt.neva
```

```neva
// userfmt.neva
// exported component with exactly one input struct and one output struct
type FormatUserInput struct {
    name string
    age int
    active bool
}
type FormatUserOutput struct {
    greeting string
    summary string
}
pub def FormatUser(input FormatUserInput) (output FormatUserOutput) {
  // implementation can be any neva graph; ports define the external api
}
```

the go backend maps the single input port to the start struct message and the single output port to the stop struct message. go-facing types stay idiomatic (plain structs with string/int/bool fields) and conversion happens in generated wiring.

## Goals & Constraints

* Generate Go packages that feel idiomatic to Go developers so Neva programs can be invoked from Go with the same ease as calling ordinary functions.
* Do not require changes to the Neva runtime or stdlib. All custom interop code must live in generated files or user packages.
* Reuse the existing Go backend so we continue to emit ready-to-run Go code without runtime unmarshalling or reflection.
* Keep performance overhead low by reusing `runtime.Program` execution semantics and avoiding extra goroutines.
* Leave room for future hooks (e.g. swapping the runtime implementation) without forcing that surface into the initial MVP.

This iteration deliberately focuses on **Go → Neva**. Calling user Go code from Neva will be revisited once the package flow proves out.

## Current Backend Behaviour

Running `neva build --target=go` produces an executable:

1. The compiler builds IR for a single entry-point program.
2. The `Middleend` orchestrates the pipeline (Analysis -> Desugaring -> IR Generation).
3. The `Backend` receives the IR program and emits `main.go`, wiring calls into one `runtime.Program`.
4. The backend copies `internal/runtime/**` and related support files next to the entry point so the emitted Go compiles in isolation.

Package mode extends this pipeline: instead of generating one main program, it discovers all exported components, generates an IR program for each export, and emits library-friendly files (`api.go`, `programs.go`) in a reusable Go package.

## CLI Surface

Package mode extends the existing command:

```bash
neva build --target=go --target-go-mode=pkg --output=./internal/gen ./pkg/userfmt
```

* `--target=go` continues to select the Go backend.
* `--target-go-mode` accepts `executable` (default) or `pkg`.
* `--output` points to the directory where generated Go code is written.
* The final positional argument stays the path to the Neva **package** being compiled.

## Exported Component Handling

Executable builds assume one entry point. Library emission surfaces every `pub def`.

**Stage 1 (Current):** Only single-port exports are fully supported.
- If a component has 1 input port `a`, the Go API takes a struct with field `A`, and the runtime sends the message directly to `a`.
- If a component has multiple input ports (`a`, `b`), the generated Go code *looks* correct (accepts a struct with `A` and `B`), and sends a `StructMsg{A, B}`. **However**, the runtime wiring connects this message to only *one* of the component's ports (e.g. `a`). The component will likely fail because:
  1. It receives a `StructMsg` instead of the expected primitive on port `a`.
  2. It receives nothing on port `b`.

**Stage 2 (Future):** Multi-port exports will be supported via bundling.
- A "bundling adapter" component will be synthesized by the Desugarer.
- This adapter will accept the `StructMsg`, split it into separate messages for `a` and `b`, and forward them to the user's component.
- The `Backend` emission logic will remain largely the same, but it will target the adapter instead of the raw component.

The compiler performs these steps when `--target-go-mode=pkg` is selected:

1. **Analysis & Desugaring:** The `Middleend` runs analysis in library mode (empty main package) and desugaring.
2. **Export Discovery:** `pkg.GetInteropableComponents()` identifies eligible exports.
3. **IR Generation:** The `Middleend` calls `Irgen.GenerateForComponent` for each export, producing multiple IR programs.
4. **Backend Emission:** The `Backend` receives the list of exports (each containing AST component info and IR program) and emits `exports.go` using a library template.

This design presents exported components as *call/return* style APIs to Go code.

## Generated Package Layout

Given `--output=./internal/gen`, the backend emits a single Go package rooted exactly at `internal/gen` with package name `gen` (or derived from go.mod), plus the runtime copy:

```
internal/gen/
  runtime/…   # local copy of runtime
  exports.go  # per-export: typed structs, factory, and free function
```

## Go → Neva Workflow

### Step-by-step

1. **Initialize Module:** `go mod init example.com/myapp/gen` in the output directory (or ensure parent dir has go.mod).
2. **Generate:** `neva build --target=go --target-go-mode=pkg --output=./gen ./src`.
   - The compiler detects the module path from `go.mod`.
   - It generates code importing the local `runtime` package (e.g. `example.com/myapp/gen/runtime`).
3. **Import & Call:**
   ```go
   import "example.com/myapp/gen"
   // ...
   out, err := gen.MyFunc(ctx, gen.MyFuncInput{...})
   ```

No `go mod replace` or `go get` of the Neva repository is required because the runtime is vendored and imports are rewritten to be local.

## Architecture: Proper Compiler Abstraction

The compiler supports two distinct compilation modes through the `Middleend` and `Backend` interfaces:

### Middleend

The `Middleend` struct (in `compiler.go`) orchestrates the pipeline:

* `ProcessExecutable(feResult) (*ir.Program, error)`: Analyzes (with main), desugars, and generates single IR program.
* `ProcessLibrary(feResult) ([]LibraryExport, error)`: Analyzes (library mode), desugars, discovers exports, and generates IR program for each export.

### Backend

The `Backend` interface (in `contract.go`) is pure and does NOT depend on analyzer/desugarer/irgen:

* `EmitExecutable(dst string, prog *ir.Program, trace bool) error`: Emits standalone program.
* `EmitLibrary(dst string, exports []LibraryExport, trace bool) error`: Emits library code for multiple exports.

Each backend (`golang`, `ir`, `native`, `wasm`) implements this interface. `native` and `wasm` wrap `golang` backend for executable emission but panic for library mode currently.

## Implementation Plan (Stage 1 - Completed)

1. **Backend Interface:** Updated to `EmitExecutable` and `EmitLibrary`.
2. **Middleend:** Refactored to handle `ProcessExecutable` and `ProcessLibrary`.
3. **Golang Backend:** Implemented `EmitLibrary` with `exports.go` generation and runtime copying with import rewriting.
4. **CLI:** Updated `build` and `run` commands to use new `compiler.Compile` flow with `Mode` flag and consistent IR format validation.
5. **E2E Test:** Verifies `go mod init` -> `neva build` -> `go run` flow without external dependencies.
6. **Self-Contained Modules:** Compiler now supports building modules where the main package is at the root (`.`), mirroring Go's `go.mod` + `main.go` behavior.
7. **Complex Type Support:** Backend `getGoFromMsg` and `getMsgFromGo` now support:
    - Primitives (`int`, `string`, `bool`, `float`) -> Native Go types
    - Struct Literals -> `runtime.StructMsg`
    - Union Literals -> `runtime.UnionMsg`
    - `list<T>` -> `[]runtime.Msg`
    - `dict<T>` -> `map[string]runtime.Msg`
    - Named Complex Types (e.g. `EntityRef` struct) -> `runtime.Msg` interface (fallback)

## Stage 2 (Future)

* **Multi-port Exports:** Implement bundling in Desugarer (synthesize wrapper component with single struct in/out ports).
* **Client Type:** Generate a `Client` struct to allow shared runtime instance and configuration.
* **Streaming:** Explore support for streaming APIs (channel-based interaction) if Go semantics allow.
* **Type Resolution for Named Structs:** Improve backend `mapFields` to resolve named types (like `EntityRef`) to their underlying struct definitions, avoiding generic `runtime.Msg` fallback.
