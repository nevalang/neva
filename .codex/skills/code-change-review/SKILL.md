---
name: "code-change-review"
description: "Run a focused review of a Neva code change, branch, diff, or pull request by combining correctness, readability, performance, and security perspectives."
---

# Code Change Review

Use this skill for high-signal review of a code change. The change may be a
pull request, a branch, a local diff, or a patch. The goal is to publish only
required changes or blocking questions.

## Workflow

1. Read the diff and required repository context:
   - root `AGENTS.md`
   - nested `AGENTS.md` for touched paths
   - docs/style guides when the diff depends on them
2. Review from four focused perspectives:
   - correctness: `references/review-correctness.md`
   - readability: `references/review-readability.md`
   - performance: `references/review-performance.md`
   - security: `references/review-security.md`
3. Use subagents when available and useful. If subagent tooling is unavailable,
   run the four review passes yourself.
4. Deduplicate overlapping findings.
5. Output only actionable findings and blocking questions. If there are no
   meaningful findings, say that clearly.

## Output Contract

- Findings first, ordered by severity.
- Include file/line references whenever possible.
- Do not include praise, strengths, or broad summary prose before findings.
- Do not spend review budget on style-only remarks unless they materially affect
  correctness, maintainability, performance, or security.
