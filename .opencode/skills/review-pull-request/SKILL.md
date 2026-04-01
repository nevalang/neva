---
name: review-pull-request
description: Review pull requests with focused subagents and prefer precise inline GitHub review comments.
---

# Review Pull Request

Use this skill when reviewing a pull request in GitHub.

## Review Goal
- Help the author make the pull request objectively better.
- Focus on real risk, real confusion, or real missed opportunity.
- Stay silent when there is nothing useful to say.

## Comment Policy
- Prefer inline file/line comments whenever there is a concrete location.
- When a concrete finding exists, publish it directly at that file/line instead of handing it off to another agent to relay.
- Comments must be either `actionable` or `questionable`.
- `actionable` means proposing a concrete change.
- `questionable` means asking a clear, unambiguous question the author can answer directly.
- Avoid vague observations without a clear next step.
- Avoid bulk summary comments when the feedback can be attached precisely to code.

## Nits
- `nit:` is optional-only feedback.
- Use `nit:` sparingly.
- If the comment is not worth the author's attention, do not write it.

## Workflow
1. Use the pull request context already available in the GitHub run.
2. Do not re-fetch pull request metadata or diff through ad-hoc shell commands unless the built-in context is missing.
3. Split review by focus area instead of mixing everything into one reviewer voice.
4. Prefer subagents for:
   - correctness
   - readability
   - performance
   - security
5. Let focused reviewers publish their own precise comments directly.
6. Only leave a top-level summary comment if there is important high-level feedback that cannot be attached to a specific file or line.

## Output Quality
- Be concrete.
- Be concise.
- Prefer high-confidence comments over broad speculative review.
