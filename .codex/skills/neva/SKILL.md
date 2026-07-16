---
name: "neva"
description: "Use for Neva source or snippets: authoring, refactoring, debugging, or review."
---

# Neva

Use this skill for changes touching `*.neva` or Neva code snippets.

## Rules

- Read `docs/user/style_guide.md` before editing.
- Keep edited `.neva` files syntactically valid.
- If syntax is uncertain, consult `internal/compiler/parser/neva.g4` and nearby
  `.neva` examples before editing.
- When feasible, validate changed Neva programs by compiling or running them.
- When editing documentation or snippets, prefer a quick local compile/run check
  of an equivalent minimal example.
