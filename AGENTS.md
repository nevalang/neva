# Neva Engineering Guide

This file is the repository's high-level engineering map for humans and robots.
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

```text
.
|- cmd/                    CLI entrypoint.
|- internal/
|  |- compiler/            Frontend, lowering, IR, and code-generation backends.
|  `- runtime/             Messages, ports, program execution, and native functions.
|- std/                    Public standard-library packages and component contracts.
|- e2e/                    Isolated regression modules with Go test harnesses.
|- examples/               User-facing executable examples in one shared module.
|- benchmarks/             Explicit runtime and language performance scenarios.
|- docs/
|  |- user/                Language behavior, APIs, style, and learning material.
|  `- developer/           Compiler, runtime, test, and contributor guidance.
`- .codex/                 Repository-local skills, plans, and agent automation.
```

The language server lives in `nevalang/neva-lsp`, not in this repository.

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

## Validation

Prefer the smallest meaningful scope first, then widen only when needed.

1. Run focused lint and tests for the changed subsystem.
2. For parser work, update grammar artifacts and run parser smoke tests.
3. For runtime and standard-library behavior, run focused unit tests and e2e
   coverage; benchmarks do not replace behavior tests.
4. Run broader `go test ./...` only when the blast radius warrants it.

Each e2e directory is independently testable. Before pushing a behavior
change, select and run the relevant `go test ./e2e/<case>` packages locally;
use the full e2e suite only when the changed contract has broad reach.

The repository's Codex stop hook runs `pre-commit` after a turn that leaves
working-tree changes, so formatting and autofix output is available before the
next response or push. Review and trust the hook explicitly when Codex asks.

For PR review comments, apply code or documentation changes first, reply to
every addressed comment through GitHub, and do not resolve threads unless the
user requests it. Verify the latest unresolved threads after replying.
