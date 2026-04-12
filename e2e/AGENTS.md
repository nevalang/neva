# AGENTS.md

This file defines local workflow for `e2e/` tests for both humans and machines.

## Purpose

- `e2e/` contains end-to-end test modules.
- These packages are not user-facing examples; see [`examples/`](../examples/).

## Structure

- Each e2e test is a separate Neva module.
- Every test package includes Neva source plus a Go test harness.
- One broken module should not prevent other e2e modules from compiling.

## Execution

- Run targeted e2e packages while iterating.
- Use broader e2e runs only when needed, because they are expensive.
