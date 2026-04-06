# .agent layout

Purpose: keep local agent artifacts organized and predictable.

## Canonical directories

- `skills/`: reusable skill definitions (`*/SKILL.md`).
- `plans/`: active planning documents.

## Legacy directories

- `workflows/`, `tasks/`, `notes/` are deprecated.
- New content should not be added there.
- Existing one-off task notes should be tracked as GitHub issues.
- Long-form vision/strategy notes should live in GitHub Wiki.

## Related roots

- `.opencode/` remains the runtime source for OpenCode CI agents/skills.
- `.claude/rules/` remains the source for file-type authoring rules.
- Cross-root unification is tracked separately in issue #1094.
