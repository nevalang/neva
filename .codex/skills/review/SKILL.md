---
name: "review"
description: "Use for focused review of a Neva pull request, branch, diff, or patch."
---

# Review

Use this skill for high-signal review of a code change. The change may be a
pull request, a branch, a local diff, or a patch. The goal is to publish only
required changes or blocking questions.

## Workflow

1. Read the diff and required repository context:
   - `AGENTS.md`
   - canonical documentation named by its change-routing table
2. Select the focused perspectives that match the change and its risk:
   - correctness: `references/review-correctness.md`
   - style and readability: `docs/user/style_guide.md` for Neva code, plus
     `references/review-readability.md`
   - performance: `references/review-performance.md`
   - security: `references/review-security.md`
3. Choose whether to use subagents based on the change's size and risk. For a
   small, focused diff, review directly; for a broad or high-risk diff, split
   the selected perspectives into focused passes with task-specific prompts.
4. Deduplicate overlapping findings.
5. Output only actionable findings and blocking questions. If there are no
   meaningful findings, say that clearly.

## Output Contract

- Findings first, ordered by severity.
- Include file/line references whenever possible.
- Keep the output proportionate to the findings; do not impose a template when
  there is no meaningful feedback.
- Report clear violations of `docs/user/style_guide.md` for Neva code. Do not
  report subjective style preferences.
