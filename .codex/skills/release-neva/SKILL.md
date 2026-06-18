---
name: "release-neva"
description: "Prepare a Neva GitHub release draft from merged PRs, previous release style, and local multi-platform artifacts. Use this for monthly release preparation in nevalang/neva."
---

# Neva Release Skill

## Overview

Use this skill to prepare a monthly release for `nevalang/neva` with consistent notes and assets.
Default output is a **draft** GitHub release, not immediate publication.

## Core Tools

- `gh` (release + PR metadata)
- `git` (range and branch state)
- `make` (artifact build)
- `shasum` (asset checksums)
- GitHub web UI (optional final review/publish)

## Workflow

1. Confirm version bump policy before tagging.
   - If the bump is ambiguous, ask explicitly: `patch` or `minor`.
   - Default to `patch` for internal-only changes (CI, docs, lint, infra, refactors without user-facing behavior change).
   - Use `minor` for user-facing language/stdlib/CLI features or meaningful behavior additions.
   - `major` is out of routine flow; always ask.
2. Sync to clean, up-to-date main state.
   - `git fetch origin --tags --prune`
   - ensure working tree is clean
   - ensure build source equals `origin/main` (or use a clean worktree from it)
3. Read release style from recent history.
   - inspect 1-2 previous release notes via `gh release view ... --json body,publishedAt`
   - do not download previous binary assets
4. Collect merged PRs in the release window.
   - determine range from previous tag (for example `vX.Y.Z..origin/main`)
   - gather merged PR numbers and read each PR description (`gh pr view`)
5. Draft release notes.
   - keep style concise and factual
   - include: theme/title, summary, key sections, related PR list, full changelog compare link
   - if changes are mostly internal, say that directly
   - include examples for user-facing changes:
     - CLI changes: add at least one command snippet showing new command/flag usage
     - language changes: add at least one Nevalang code snippet
   - keep examples small; full compilable wrappers are optional
6. Build assets locally from clean main state.
   - `make build`
   - verify expected platform binaries are present
   - compute checksums with `shasum -a 256`
7. Create draft release with assets.
   - `gh release create <tag> <assets...> --target main --title <tag> --notes-file <file> --draft`
   - verify metadata and uploaded assets via `gh release view <tag> --json ...`
8. Report back.
   - release URL
   - included PRs
   - asset names and checksums
   - note whether release remains draft (default)

## Guardrails

- Do not publish the release unless explicitly requested.
- Do not invent release claims not supported by merged PR descriptions.
- Prefer `gh` over ad-hoc API calls where possible.
- Keep the skill lightweight; avoid scripts unless explicitly requested.
- Do not pad examples with unnecessary boilerplate (`def Main`, imports, full component wrappers) when a focused snippet is enough.

## Example Quality Check (Nevalang)

When release notes include Nevalang snippets, they must be concise but real.

1. Keep snippets minimal.
   - Prefer showing only the relevant node declarations + `---` + connections when that explains the change.
   - Full compile-ready files are optional.
2. Keep snippets language-true.
   - syntax must match current grammar (`internal/compiler/parser/neva.g4`)
   - semantics should follow project invariants from `AGENTS.md` (for example `Main` port shape, chaining rules, `[*]` bypass form)
3. Sanity-check against real code.
   - compare syntax patterns with 1-2 existing examples (`examples/**/*.neva`)
   - if confidence is low, rewrite to a simpler canonical form before publishing
