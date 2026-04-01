# PR Review Rubric

Use this rubric to keep PR reviews consistent and actionable.

## Priorities
- `P0` Critical: data loss, security issue, crash, or hard correctness break.
- `P1` High: likely production bug or clear behavior regression.
- `P2` Medium: important maintainability/testability gap with realistic risk.
- `P3` Low: minor issue or optional improvement.

## What To Look For
- Incorrect behavior changes vs intended semantics.
- Missing error handling or unsafe runtime assumptions.
- Concurrency/race/deadlock risks.
- Backward compatibility concerns.
- Missing or weak tests for changed behavior.

## What To Avoid
- Pure style nits unless they hide a real risk.
- Rewriting project conventions in each review.
- Vague comments without impact explanation.

## Reporting Rules
- Prefer fewer, high-confidence findings.
- Always include concrete file references.
- If uncertain, call it out explicitly.
