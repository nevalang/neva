# Papercuts

Private maintainer journal for small repository and engineering-harness
frictions encountered by LLMs. Entries are generated with the repository-local
skill that owns this file.

---

<!-- Add papercuts below this line. Put new entries first. -->

## 2026-08-02 12:05 +07 — Codex

Validating this documentation change with `pre-commit run --files ...` ->
`InvalidConfigError: .pre-commit-config.yaml is not a file`, despite the
engineering guide saying the stop hook runs pre-commit after a changed turn.
Used `git diff --check` for this documentation-only change. Likely fix: make
the hook detect a repository configuration before invoking pre-commit, or
document the intended configuration location.

## 2026-08-02 12:08 +07 — Codex

Pushing this documentation-only branch -> the pre-push `make lint` gate failed
on lint violations unchanged from `origin/main`; it also reported files from a
sibling worktree. The same gate's unit tests and vulnerability scan passed.
Likely fix: restore a clean lint baseline and isolate the lint scan to the
current worktree.
