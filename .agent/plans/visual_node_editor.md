# Visual Node Editor Plan (AST-Level, Mermaid-First, Readonly)

## Purpose

Design and deliver a visual node editor model for Neva that:

- reflects source-level program structure (not compiler-lowered artifacts),
- is host-agnostic (works for IDE extension and standalone app),
- starts with a readonly Mermaid backend, then evolves toward richer interactive capabilities.

## Current Architecture Baseline

- Source model: `internal/compiler/ast/flowast.go`
- Compiler flow: builder -> parser -> analyzer -> desugarer -> irgen -> backend
- LSP server: `cmd/lsp/` with indexing over analyzed AST
- LSP now includes core language features (definition/references/implementation/rename/hover/document symbols/semantic tokens/codelens) in `cmd/lsp/server`
- Existing custom visual-related route/types still present: `resolve_file` and `GetFileView*` structs in `cmd/lsp/server` (and `GetFileView` is still not wired)
- Runtime interception currently exposes `Sent` and `Received` events
- Existing IR Mermaid backend exists at `internal/compiler/backend/ir/mermaid/` (useful reference for file/output shape, but not source-level semantics)

Implication: visual representation should be projected from analyzed AST, while runtime overlays should map to `Sent/Received` events.

## Design Principles

- Compiler and runtime remain independent from UI concerns.
- Keep projection layer minimal: normalize only what renderers need, do not clone AST.
- One canonical normalized graph payload shared by Mermaid and future interactive viewer.
- Component is the primary canvas unit to keep graphs understandable.
- Source anchors are first-class for graph <-> code navigation.
- Desugared synthetic nodes remain implementation detail, not primary UX.

## Terminology

- **Normalized graph payload (V1)**: renderer-facing schema derived from analyzed AST; intentionally thin.
- **Projection layer**: backend module that converts analyzed AST into normalized graph payload with stable IDs and explicit connection semantics.
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

- Introduce minimal projection layer from analyzed AST -> `GraphDocument`
- Define stable IDs for nodes/ports/edges
- Ensure Mermaid backend consumes payload, not raw AST traversal directly
- Add golden tests for projection invariants

## Phase 2: Stabilize Current LSP Bridge

- Implement and validate `GetFileView` for current `resolve_file`
- Return enough analyzed AST context or projected payload references for readonly consumers
- Add tests for response shape and failure behavior
- Reuse existing LSP helpers introduced for core features (file resolution, location/range conversion, indexed build access) instead of duplicating logic

## Phase 3: Versioned Visual LSP Endpoints

- Add explicit visual methods, for example:
  - `neva/getGraphDocument`: returns `version`, `capabilities`, document metadata, and projected graphs.
  - `neva/getComponentGraph`: returns one `ComponentGraph` by component identity.
  - `neva/getPackageOutline`: returns package entities and navigation metadata.
- Use explicit schema versioning in payloads (`version: 1`) and additive evolution rules.
- Keep transition compatibility with `resolve_file` until new endpoints are adopted.
- Add compatibility tests for both old (`resolve_file`) and new (`neva/get*`) flows.
- Prefer implementing these endpoints through the same request-routing conventions now used by merged LSP feature handlers.

## Phase 4: Interactive Readonly Viewer MVP

- Build interactive readonly viewer (React + graph lib is one viable option, but not mandatory in this plan)
- Render component canvas with self-ports, selectors, constants, fan junctions
- Add deterministic autolayout
- Implement graph -> source and source -> graph navigation
- No write-back/editing in this milestone

Note on Mermaid:

- Mermaid is the first delivery target (Phase 0) for readonly output and validation.
- Mermaid remains useful long-term for docs/export.
- For fully interactive graph UX, a dedicated graph UI stack is still expected.

## Phase 5: Optional Runtime Overlay Layer

- Add non-canonical overlay payload keyed by stable graph IDs
- Map runtime `Sent/Received` to edge activity visualization
- Keep runtime overlays optional and decoupled from core AST snapshot

## Phase 6: Standalone Parity

- Reuse the same projector and contract in standalone host (same graph semantics and payload format).
- Keep host-specific state outside canonical payload (window layout, local preferences, shortcuts, session state).
- Parity means both hosts render the same source program consistently from the same contract.

## Dependency Note: LSP Baseline

- Core LSP feature groundwork is already merged into `main`.
- Visual-editor work should proceed immediately using this new baseline.
- Remaining dependency is visual-specific API wiring (`resolve_file` completion and/or new `neva/get*` methods), not generic LSP capabilities.

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
