# Visual Node Editor Plan (AST-Level, Mermaid-First, Readonly)

## Purpose

Design and deliver a visual node editor model for Neva that:

- reflects source-level program structure (not compiler-lowered artifacts),
- is host-agnostic (works for IDE extension and standalone app),
- starts with a readonly Mermaid backend, then evolves toward richer interactive capabilities.

## Current Architecture Baseline

- Source model: `pkg/ast/flowast.go` (with shared location/types in `pkg/core`)
- Compiler flow: builder -> parser -> analyzer -> desugarer -> irgen -> backend
- LSP server: `cmd/lsp/` with indexing over analyzed AST
- LSP now includes core language features (definition/references/implementation/rename/hover/document symbols/semantic tokens/codelens) in `cmd/lsp/server`
- Existing custom visual-related route/types still present: `resolve_file` and `GetFileView*` structs in `cmd/lsp/server` (currently considered legacy/incomplete for this plan)
- Runtime interception currently exposes `Sent` and `Received` events
- Existing IR Mermaid backend exists at `internal/compiler/backend/ir/mermaid/` (useful reference for file/output shape, but not source-level semantics)

Note: AST/core were extracted from `internal/compiler/*` to `pkg/*`, which makes source-model types importable by external tooling modules.

Implication: base visual representation should be projected from raw/source AST, while semantic/runtime enrichments are layered on top when available.

## Design Principles

- Compiler and runtime remain independent from UI concerns.
- Raw AST is the base truth for visualization in v1.
- Keep projection layer minimal: normalize only what renderers need, do not clone AST.
- One canonical normalized graph payload shared by all renderers (Mermaid preview, IDE visual view, standalone app) so graph semantics are defined once.
- Component is the primary canvas unit to keep graphs understandable.
- Source anchors are first-class for graph <-> code navigation.
- Desugared synthetic nodes remain implementation detail, not primary UX.
- Reuse existing LSP language features where possible (definition/references/rename/hover/symbols/semantic tokens) instead of reimplementing them in visual-specific APIs.

## Terminology

- **Normalized graph payload (V1)**: renderer-facing schema derived primarily from raw/source AST; intentionally thin.
- **Projection layer**: backend module that converts raw/source AST into normalized graph payload with stable IDs and explicit connection semantics.
- **GraphDocument**: top-level payload for requested scope (program/package/file/component) with metadata and graphs.
- **ComponentGraph**: graph for one component (self ports, nodes, edges, anchors).
- **GraphNode**: component node instance with resolved interface/type metadata needed for rendering.
- **GraphEdge**: source-level connection representation with sender/receiver semantics.
- **SourceAnchor**: source range/location pointers used for graph-to-text and text-to-graph navigation.
- **Capabilities**: feature flags in payloads for safe progressive adoption.

## Normalized Graph Payload (V1)

Define a stable data contract with these concepts:

- `GraphDocument`: module/package/file metadata, entity list, diagnostics summary
- `ComponentGraph`: one component graph with nodes, self-ports, edges
- `GraphNode`: entity ref, resolved interface snapshot, type args, err-guard, overload metadata, DI metadata
- `GraphEdge`: sender/receiver endpoints, selectors, constant sender metadata, chain metadata
- `SourceAnchor`: source location anchors used for navigation and highlighting
- `Capabilities`: feature flags (array bypass, selectors, runtime overlay support)

### AST Mapping Rules

- `Component` -> canvas graph
- component `IO` -> boundary self-ports
- `Node` -> visual node card
- `Connection` -> visual edge(s)
- multi-sender / multi-receiver connections -> explicit visual junction semantics (merge/split)
- sender/receiver `[*]` (array bypass) -> dedicated bypass edge semantic
- `StructSelector` -> edge selector labels
- const sender -> inline constant marker

## Scope Model

- Primary graph scope: component
- File/package scope: navigation and filtering layers above component graphs

This keeps large modules usable while preserving drill-down.

## Delivery Phases

## Phase 0: Mermaid Readonly Backend (Start Here)

- Implement source-level Mermaid backend that produces markdown file(s) for a given program/package
- Reuse patterns from IR Mermaid backend where useful (encoding, file writing), but map source semantics from AST
- Support core source semantics in output:
  - component-level graphs
  - self ports
  - node instances
  - fan-in/fan-out
  - chained connections
  - selectors
  - const senders
  - array bypass (`[*]`)
- Add snapshot tests for representative examples/e2e fixtures
- Keep this output strictly readonly/documentational

## Phase 1: Projection Layer + Payload V1

- Introduce minimal projection layer from raw/source AST -> `GraphDocument`
- Define stable IDs for nodes/ports/edges
- Ensure Mermaid backend consumes payload, not raw AST traversal directly
- Add golden tests for projection invariants
- Add optional (non-blocking) semantic enrichment path that can attach analyzer-derived metadata when available

## Phase 2: Visual LSP API V1 (No Legacy Migration)

- Introduce explicit visual methods:
  - `neva/getGraphDocument`
  - `neva/getComponentGraph` (optional optimization endpoint, can be deferred)
- Do not depend on `resolve_file` compatibility for rollout
- Return raw-AST-based graph payload plus optional enrichment sections
- Add tests for request/response shape and failure behavior
- Reuse existing LSP helpers introduced for core features (file resolution, location/range conversion, indexed build access) instead of duplicating logic

## Phase 3: Versioning, Capabilities, and Optional Enrichments

- Use explicit schema versioning in payloads (`version: 1`) and additive evolution rules
- Add `capabilities` flags to advertise optional enrichments, e.g.:
  - `semanticOverlay`
  - `desugaredOverlay` (deferred / disabled in early iterations)
  - `textGraphSync`
- Keep base payload usable without any overlay
- Prefer implementing visual methods through the same request-routing conventions used by merged LSP feature handlers

## Phase 4: Interactive Readonly Viewer MVP

- Add VSCode command-mode preview first (similar to Markdown preview), e.g. `Neva: Preview Mode`
- Preview renders current file in readonly visual form:
  - file-level view with component boxes
  - each component box contains its graph view
  - option to open a component in focused/fullscreen view
- No graph editing/write-back in this milestone
- Keep implementation compatible with later evolution to richer visual modes

Note on Mermaid:

- Mermaid is the first delivery target (Phase 0) for readonly output and validation.
- Mermaid remains useful long-term for docs/export.
- For fully interactive graph UX, a dedicated graph UI stack is still expected.

## Phase 5: Optional Runtime Overlay Layer

- Add non-canonical overlay payload keyed by stable graph IDs
- Map runtime `Sent/Received` to edge activity visualization
- Keep runtime overlays optional and decoupled from core AST snapshot

## Phase 6: Standalone Parity

- Standalone app depends on Neva LSP as backend transport (not VSCode-specific).
- Reuse the same projector and contract through LSP in standalone host (same graph semantics and payload format).
- Keep host-specific state outside canonical payload (window layout, local preferences, shortcuts, session state).
- Parity means both hosts render the same source program consistently from the same contract.

## Dependency Note: LSP Baseline

- Core LSP feature groundwork is already merged into `main`.
- Visual-editor work should proceed immediately using this new baseline.
- Remaining dependency is visual-specific API wiring (`neva/get*` methods), not generic LSP capabilities.

## IDE Integration Modes (Open UX Decision)

The plan supports both IDE and standalone. Exact IDE UX architecture is intentionally left open for next iteration.
Current committed MVP direction: command-based preview mode.

Potential IDE patterns worth evaluating next:

1. Custom editor (notebook-like experience)
2. Side panel / webview synchronized with text editor
3. Inline component visualization affordances in text view (peek/minimap/toggle)

Ideas to keep for follow-up:

- side panel + selection sync
- button/code lens quick preview
- optional full visual mode

Open questions to resolve before post-MVP UX work:

- Should preview become a custom editor mode, remain a command preview, or support both?
- How should text/visual synchronization lifecycle work (cursor sync, reveal behavior, focus ownership)?
- Which transition model is best for component focus/fullscreen (inline expand, secondary editor, dedicated panel)?

## Risks and Mitigations

- Contract drift as language evolves
  - Mitigation: versioned schema + golden tests
- Graph density in larger components/packages
  - Mitigation: component-first scope + progressive disclosure
- Host coupling
  - Mitigation: thin host adapters over shared projector and contract

## Non-Goals (Initial)

- Direct graph editing / source write-back
- Showing desugared synthetic runtime/compiler nodes as primary representation
- Locking the implementation to a single UI stack before contract stabilizes

## Success Criteria

- Mermaid backend generates stable, readable markdown graphs for source-level Neva semantics
- One canonical normalized payload consumed by both Mermaid backend and later interactive viewer
- Stable graph identities across refreshes for unchanged code
- Reliable bidirectional navigation between graph and source
- Readonly outputs support full source-level connection semantics used in Neva code today
