# Review Pull Request

Run a fast, high-signal review over a pull request diff and produce a concise report.

## Tools
- `Read` for PR metadata, changed files, and diff.
- `Write` to save the review report.

## Workflow
1. Read PR context from the prompt (PR number, title, description, metadata path, diff path, output path).
2. Read `.opencode/shared/pr-review-rubric.md`.
3. Review the diff for correctness, regressions, runtime risks, and test gaps.
4. Prefer concrete findings over broad style advice.
5. For each finding, include:
   - severity (`P0`..`P3`)
   - short title
   - file path and line reference from the diff when available
   - 1-3 sentence explanation and expected impact
6. If there are no blocking issues, say that explicitly.
7. Save the final report to the exact output path requested by the caller.

## Output Format
Use this markdown structure:

1. `# PR Review`
2. `## Summary` with 2-4 bullets.
3. `## Findings`
   - If findings exist: numbered list, most severe first.
   - If no findings: `No blocking findings.`
4. `## Suggested Checks` with a short numbered list of tests/checks worth running.

## Output Rules
- Do not include code fences around the whole report.
- Do not add wrapper text like "Here is the review".
- Keep the report concise and actionable.
- Avoid repeating unchanged PR description text.
