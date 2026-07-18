# Compiler

The compiler turns a resolved multi-module source build into an executable Go
program or an interop library. Its implementation is under
`internal/compiler/`:

```text
internal/compiler/
|- parser/                ANTLR grammar, generated parser, and AST construction.
|- analyzer/              Name, graph, port, and type validation.
|- typesystem/            Type resolution, compatibility, and subtyping.
|- desugarer/             Lowers validated source-level conveniences.
|- ir/                    Backend-neutral program representation.
|- irgen/                 Expands component graphs into IR connections and calls.
`- backend/
   |- golang/             Emits generated Go and the embedded runtime.
   `- ir/                 Emits inspectable JSON, YAML, DOT, Mermaid, or ThreeJS IR.
```

The [architecture page](./architecture.md) shows the surrounding build flow.
User documentation owns public language semantics; this page records the
implementation boundaries that preserve them.

## Pipeline

```text
builder -> parser -> analyzer -> desugarer -> IR generator -> backend
```

`Compiler.Compile` first asks the builder for a resolved module graph, then
parses every module into an AST build. The analyzer validates that source build;
the desugarer then lowers validated source constructs. Executable compilation
generates IR rooted at `Main`. Library compilation generates a separate IR
program for each interopable exported component.

Analyze before desugaring. This keeps diagnostics in the user's source terms
and prevents lowering unchecked constructs. IR generation consumes an analyzed,
desugared build; it must not parse source syntax or recreate source-level type
rules.

## Parser and Analysis

The parser produces structured AST data. Do not leave a syntax-level
mini-language encoded in strings for later stages to parse again. When syntax
changes, update the grammar and generated parser artifacts together and run the
parser smoke tests.

The analyzer owns graph correctness: entities resolve in their lexical/module
scope, connections refer to valid compatible ports, and type expressions are
resolved through `typesystem/`. Some standard-library directives are compiler
contracts rather than ordinary library calls. For example, struct selectors and
directives such as `#extern`, `#bind`, and `#autoports` require analysis that
cannot be expressed solely in Neva source. Keep that knowledge explicit and
narrowly scoped; do not disguise it as a component-specific language exception.

## IR Generation

`irgen.Generator` recursively expands an exported root component. For each
native `#extern` node it emits an IR function call with its addressed input and
output ports and configuration message. For an ordinary Neva component it
walks the validated network, tracks the ports actually used by every subnode,
and recursively emits the component graph. The result is `ir.Program`, a list
of function calls plus a connection map between concrete port addresses.

IR generation may rely on analyzer invariants. An impossible cross-stage state
is an internal error, not a user diagnostic. Conversely, a malformed user
program must fail in parser or analyzer with `*compiler.Error` before it reaches
IR generation.

## Backends

The Go backend graph-reduces intermediate connections, assigns channel
variables to the remaining port addresses, generates runtime-function calls,
and writes a standalone Go module with the required runtime files. It supports
both executable output and exported library wrappers.

The IR backend is an inspection backend, not a separate lowering path: it
serializes the same `ir.Program` as JSON, YAML, DOT, Mermaid, or ThreeJS. Use it
to inspect graph lowering or support visualization tooling without conflating
that representation with generated Go.

## Type-System Changes

Before changing grammar, AST, analyzer, or type-system code, create a minimal
reproducer and establish which existing language invariant it violates. Check
whether the requested behavior is already expressible by an existing type,
union, component, or standard-library contract. Do not introduce a special
case for a component when a general language rule is the actual model.
