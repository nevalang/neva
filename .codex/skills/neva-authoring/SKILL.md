---
name: "neva-authoring"
description: "Use when editing, reviewing, or generating .neva source files or Neva snippets in docs/tests. Covers syntax, style, and validation expectations."
---

# Neva Authoring

Use this skill for changes touching `*.neva` or Neva code snippets.

## Rules

- Follow the nearest `AGENTS.md` as scoped source of truth.
- Source of truth for style: `docs/style_guide.md`.
- Keep edited `.neva` files syntactically valid.
- If syntax is uncertain, consult `internal/compiler/parser/neva.g4` and nearby
  `.neva` examples before editing.
- When feasible, validate changed Neva programs by compiling or running them.
- When editing documentation or snippets, prefer a quick local compile/run check
  of an equivalent minimal example.
