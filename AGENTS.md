# Neva Engineering Guide

This file is the repository's high-level engineering map for humans and Codex.
Read the canonical document for the subsystem before changing it. Keep this
guide concise: it routes work; it does not duplicate documentation.

## Project Model

Neva is a statically typed, compiled dataflow language. Programs are explicit
node-and-connection graphs, and concurrency is the default execution model.

Compiler pipeline:

```text
parser -> analyzer -> desugarer -> IR generator -> Go backend
```

The standard library combines Neva components with native runtime-backed
components declared through `#extern`.

## Repository Map

- `cmd/`: CLI entrypoints.
- `internal/compiler/`: parser, analyzer, desugarer, IR generation, and Go backend.
- `internal/runtime/`: messages, ports, program execution, and native functions.
- `std/`: public standard-library packages and component contracts.
- `e2e/`: isolated regression modules with Go test harnesses.
- `examples/`: executable, user-facing examples in one shared module.
- `benchmarks/`: explicit runtime and language performance scenarios.
- `docs/user/`: language behavior, API usage, style, and learning materials.
- `docs/developer/`: implementation, testing, runtime, and contributor guidance.
- The language server lives in `nevalang/neva-lsp`, not in this repository.

## Documentation

Documentation is organized by primary reader, not by a one-way dependency
graph. Link between `docs/user/` and `docs/developer/` when useful, but keep a
topic canonical in one place.

- Public language behavior, API semantics, and Neva style belong in `docs/user/`.
- Compiler, runtime, standard-library implementation, and test strategy belong
  in `docs/developer/`.

Start from `docs/README.md` for documentation navigation.

## Change Routing

| Change | Read before editing | Update when behavior changes |
| --- | --- | --- |
| `*.neva` | `docs/user/style_guide.md` | User docs for public API or style |
| `std/**` | User component docs and `docs/developer/runtime-functions.md` when native | User API docs and developer implementation docs as applicable |
| `internal/runtime/funcs/**` | `docs/developer/runtime-functions.md` | Developer runtime docs; user docs only for public behavior |
| `internal/compiler/**` | `docs/developer/compiler.md` and relevant user semantics | Developer compiler docs; user docs for language-surface changes |
| `e2e/**`, `examples/**`, `benchmarks/**` | `docs/developer/testing.md` | Developer testing docs when the test contract changes |

## Stable Language Constraints

1. Preserve the 1:1 mapping between text and graph. Prefer explicit standard
   library nodes over hidden control-flow sugar.
2. `Main` has exactly one non-array `start any` inport and one non-array
   `stop any` outport.
3. Literal and const senders are valid only in signal-triggered chains.
   Primitive and union literals are allowed; bytes literals are not.
4. Array bypass uses `[*]` on both sides of one connection. Index `255` is
   reserved for wildcard bypass; do not mix bypass and indexed usage on a port.
5. Use `res` for the primary result and `err error` for failures. Propagate
   errors with `?` unless custom handling is necessary.
6. Conversions are explicit standard-library components. Do not introduce
   implicit casts.
7. `bytes` is transport-focused and has no literal syntax. Use explicit
   converters such as `bytes.FromString` and `strings.FromBytes`.
8. Keep compiler-contract standard-library behavior explicit: `builtin`,
   `Union`, `Struct`, `#autoports`, and desugaring-sensitive constructs are
   language boundaries, not ordinary helpers.
9. Return `*compiler.Error` for invalid user programs. Panic only for internal
   invariant violations or impossible cross-stage states.

## Working Protocol

1. Inspect the relevant code, tests, and canonical documentation before making
   an architectural decision.
2. Prefer existing repository patterns and public primitives over new APIs or
   special cases.
3. For a grammar, AST, analyzer, or type-system change, create a minimal
   reproducer before implementation and identify the violated language rule.
4. For a runtime function, establish why a Neva graph is insufficient and
   follow the runtime-function protocol.
5. Keep changes scoped. Record unresolved design work in an issue or a
   repository-local plan rather than shipping speculative infrastructure.

## Operating Standards

- Use `gh` for GitHub context before falling back to `curl`.
- Ask the user when material uncertainty remains after repository and
  documentation research.
- Run targeted checks with an approximately five-minute cap unless the user
  asks for a broader run.
- Start implementation branches from current `origin/main` unless the user
  explicitly asks to work on an existing branch.
- Never merge, enable auto-merge, or close a review PR without a direct user
  command in the current conversation.
- Keep this guide current when a recurring process or architecture rule changes.

## Skills

Use the matching repository skill for the task:

- `Go` for Go files.
- `Neva` for `.neva` files and Neva snippets.
- `Review` for pull requests, diffs, and branches.

## Validation

Prefer the smallest meaningful scope first, then widen only when needed.

1. Run focused lint and tests for the changed subsystem.
2. For parser work, update grammar artifacts and run parser smoke tests.
3. For runtime and standard-library behavior, run focused unit tests and e2e
   coverage; benchmarks do not replace behavior tests.
4. Run broader `go test ./...` only when the blast radius warrants it.

For PR review comments, apply code or documentation changes first, reply to
every addressed comment through GitHub, and do not resolve threads unless the
user requests it. Verify the latest unresolved threads after replying.
