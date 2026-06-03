# Visual Editor Session Snapshot

## Date

- 2026-05-10

## Done

- Added runtime Dataflow Trace core in `neva`:
  - trace envelope on runtime messages,
  - sent/recv hop recording,
  - path reconstruction API,
  - formatted cause trace output for panic.
- Updated panic signaling semantics:
  - runtime call now tracks and returns program panic errors,
  - generated executable exits with code `1` on panic without duplicate runtime-error line.
- Added machine-readable trace event output v1 in `DebugInterceptor` (JSON lines).
- Added tests for trace ancestry and panic error semantics in `internal/runtime`.
- Created P0 master tracker issue: #1118.
- Updated `.agent/plans/visual_node_editor.md` with split-repo/LSP-externalized baseline corrections.
- Added snapshot template and updated `docs/qna.md` terminology (`Dataflow Trace` as canonical).

## Decisions

- `Dataflow Trace` is canonical term; `Graph trace` remains alias.
- Runtime trace core is implemented first; DAP remains a later adapter concern.
- Cross-repo execution model is explicit (`neva` + `neva-lsp` + `vscode-nevalang`).

## Risks

- Current ancestry model is single-parent lineage; fan-in causality trees need follow-up semantics.
- Source graph contract + compiler sidecar mapping not yet implemented in this session.
- LSP and VSCode tracks are out-of-repo and still pending.

## Next 3 Steps

1. Design and implement source-level GraphDocument/ComponentGraph contract in `neva` with stable IDs.
2. Add compiler debug sidecar mapping `SourceEdgeID <-> runtime edges` and tests for mapping invariants.
3. Open companion tracker issues in `neva-lsp` and `vscode-nevalang`, linked back to #1118.

## Linked Issues / PRs

- #1118 (master tracker)
- #970
- #977
- #595
- #1050
- #792
- #1046

