---
name: review-pull-request
description: Use when a GitHub pull request should be reviewed, either automatically from CI on PR open/update or explicitly by a reviewer.
---

# Review Pull Request

This skill marks the entrypoint for GitHub pull request review.

Use it when the task is to review a pull request in GitHub.

Do not keep the detailed review methodology here.

Responsibility split:
- `.opencode/agents/review-pull-request.md` owns orchestration and shared review policy.
- `.opencode/agents/correctness.md` owns product and logic review.
- `.opencode/agents/readability.md` owns clarity and maintainability review.
- `.opencode/agents/performance.md` owns performance and benchmark review.
- `.opencode/agents/security.md` owns security review.
