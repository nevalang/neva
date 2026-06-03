# AGENTS.md

This file defines local workflow for `examples/` for both humans and machines.

## Purpose

- `examples/` are living documentation.
- Example packages are also executed as end-to-end tests.

## Structure

- `examples/` is one Neva module.
- Each example is a separate package inside that module.
- Running one package still requires all example packages in the module to compile.

## Relation to `e2e/`

- `e2e/` contains focused test scenarios.
- `examples/` should stay readable and educational first.

## Local `AGENTS.md` Style

- Put scenario intent and topology in local example-level `AGENTS.md`.
- Prefer concise Mermaid diagrams for non-trivial flows.
- Avoid duplicating exact expected output values when those assertions are already covered by tests.
