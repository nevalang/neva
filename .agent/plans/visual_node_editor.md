# Visual Node Editor Plan (Finalized)

## Purpose

Design and deliver a visual node editor for Neva that:

- reflects source-level program structure (not compiler-lowered artifacts),
- is host-agnostic (works for IDE extension and standalone app),
- starts with a readonly Mermaid backend consuming a defined graph model, then evolves toward richer interactive capabilities,
- provides a foundation for a next-generation hybrid programming experience combining textual and visual editing.

## Current Architecture Baseline

- Source model: `pkg/ast/flowast.go` (with shared location/types in `pkg/core`)
- Compiler flow: builder -> parser -> analyzer -> desugarer -> irgen -> backend
- LSP server: `cmd/lsp/` with indexing over analyzed AST
- LSP includes core language features (definition/references/implementation/rename/hover/document symbols/semantic tokens/codelens) in `cmd/lsp/server`
- Existing custom visual-related route/types: `resolve_file` and `GetFileView*` structs in `cmd/lsp/server` (legacy/incomplete, but the `Extra.NodesPorts` pattern of attaching resolved interfaces to nodes is a useful predecessor to the graph model)
- Runtime interception exposes `Sent` and `Received` events via `Interceptor` interface in `internal/runtime/interceptors.go`
- Trace log: `--emit-trace` CLI flag produces `trace.log` with append-only `sent | <port> | <message>` and `recv | <port> | <message>` event records via `DebugInterceptor`
- Existing IR-level visualization backends:
  - Mermaid: `internal/compiler/backend/ir/mermaid/` (markdown with flowcharts from lowered IR)
  - Three.js: `internal/compiler/backend/ir/threejs/` (interactive 3D HTML visualization from lowered IR)
  - DOT/Graphviz: `internal/compiler/backend/ir/dot/`
- VSCode extension: separate repository `vscode-neva` (marketplace: `nevalang.vscode-nevalang`), contains TextMate grammar and LSP client wiring
- Analyzer supports library mode: can analyze all packages without a specific entry point (`mainPkgName == ""` path in `analyzer.go`), which the LSP indexer uses for workspace-wide analysis

Note: AST/core were extracted from `internal/compiler/*` to `pkg/*`, making source-model types importable by external tooling.

Implication: base visual representation is projected from raw/source AST, while semantic and runtime enrichments are layered on top when available.

## Design Principles

- Compiler and runtime remain independent from UI concerns.
- Raw AST is the base truth for visualization. The program at AST level is nested and recursive (modules contain packages, packages contain files, files contain entities, components contain nodes and connections). This is fundamentally different from IR where the whole program is a single flat graph. Visualization must respect this nested structure.
- Keep projection layer minimal: normalize only what renderers need, do not clone AST.
- One canonical source graph model shared by all renderers (Mermaid output, IDE visual view, standalone app, Three.js) so graph semantics are defined once. Individual renderers consume this model and may drop information they cannot represent. The model describes Neva semantics (ports, nodes, connections, fan-in, fan-out) -- not renderer syntax (Mermaid subgraphs, React Flow handles, Three.js meshes). Renderer-specific concerns belong in renderer adapters, not in the model.
- Component is the primary canvas unit.
- Source anchors are first-class for graph <-> code navigation.
- Desugared synthetic nodes are implementation detail, not primary UX.
- Reuse existing LSP language features where possible instead of reimplementing them in visual-specific APIs.
- Visualize exactly as the programmer writes: AST-level constructs (fan-in, fan-out, chained connections, struct selectors, const senders, array bypass) are shown as-is, not decomposed into desugared form.
- Use consistent Neva-native naming throughout: fan-in, fan-out, chained connection, struct selector, const sender, array bypass.

## Terminology

- **Source graph model**: renderer-facing schema derived from raw/source AST; the canonical data contract consumed by all renderers.
- **Projection layer**: module that converts raw/source AST into source graph model with stable IDs and explicit connection semantics.
- **GraphDocument**: top-level payload for a requested scope (package/file/component) with metadata and component graphs.
- **ComponentGraph**: graph for one component (self-ports, nodes, edges, anchors).
- **GraphNode**: component node instance with resolved interface/type metadata needed for rendering.
- **GraphEdge**: source-level connection representation preserving fan-in, fan-out, chained connection, and selector semantics as written in source.
- **SourceAnchor**: source range/location pointers used for graph-to-text and text-to-graph navigation.
- **Capabilities**: feature flags in payloads for safe progressive adoption.

## Source Graph Model

Define a stable data contract with these concepts:

- `GraphDocument`: package/file metadata, list of component graphs, diagnostics summary
- `ComponentGraph`: one component graph with self-ports, nodes, edges, source anchors
- `GraphNode`: entity ref, resolved interface snapshot, type args, error guard flag, overload metadata, DI metadata. Represents a node instance (a reference to another component or interface).
- `GraphLiteral`: inline constant value (bool, int, float, string, list, dict/struct, enum) attached directly to a receiver port. Distinguished from `GraphNode` in the model: literals are not independent graph entities with their own ports and edges -- they are value tags on the port they feed into. This distinction dramatically improves signal-to-noise ratio in the visual output by eliminating unnecessary boxes and edges for simple values like `'Hello, World!'` or `42`.
- `GraphEdge`: sender/receiver endpoints preserving source-level semantics:
  - Fan-in: multiple senders to one or more receivers (as written)
  - Fan-out: one sender to multiple receivers (as written)
  - Chained connections: nested sender -> receiver -> sender -> receiver chains (as written)
  - Struct selectors: field path labels on edges
  - Const senders: inline constant markers with value (rendered as port labels / value tags on the receiving port, not as separate node boxes -- see `GraphLiteral`)
  - Array bypass (`[*]`): dedicated bypass edge semantic. This is a unique Neva feature where multiple values flow through a single logical edge. Renderers should visually distinguish bypass edges from regular edges (e.g., a thicker bus-style stroke, a double line, or other visual weight) to communicate that this is a multi-value channel, not a single message wire.
- `SourceAnchor`: source location tied to every graph element for bidirectional navigation
- `Capabilities`: feature flags (array bypass, selectors, runtime overlay, semantic overlay)

### AST Mapping Rules

- `Component` -> canvas graph (the `ComponentGraph`)
- Component `IO` (interface) -> boundary self-ports on the canvas edges
- `Node` -> visual node card
- `Connection` -> visual edge(s), preserving source-level shape:
  - `len(Senders) > 1` -> fan-in visual (multiple incoming edges converging)
  - `len(Receivers) > 1` -> fan-out visual (edge splitting to multiple targets)
  - `ConnectionReceiver.ChainedConnection != nil` -> chain is flattened into sequential edge segments during projection (the model stores the linear sequence, not the recursive nesting)
  - `ConnectionSender.StructSelector` -> selector labels attached to edge
  - `ConnectionSender.Const` -> `GraphLiteral` attached to the receiving port (not a separate node box with an edge). Inline literal values (`'Hello'`, `42`, `true`) become port labels/value tags. Const references (`$myConst`) may be rendered as small reference badges.
  - `PortAddr.Idx == ArrayBypassIdx (255)` -> array bypass edge with distinct visual weight (bus-style stroke to indicate multi-value channel)

### Stable ID Strategy

IDs are not user-visible but must be stable across refreshes for unchanged code. Scheme based on AST path:

- Node ID: `{package}/{file}/{component}/{nodeName}`
- Port ID: `{package}/{file}/{component}/{nodeName}:{portName}` (with `[{idx}]` suffix for array ports)
- Self-port ID: `{package}/{file}/{component}:in:{portName}` / `:out:{portName}`
- Edge ID: content-based hash of sender and receiver endpoint IDs (stable across connection reordering in source)

The nested/recursive nature of the AST means IDs naturally scope to their containing component -- there is no global flat namespace.

## Scope Model

- Primary graph scope: **component** (the canvas)
- File scope: contains multiple component graphs and entity declarations
- Package scope: contains multiple files; used for navigation and overview

### Overview / Index View

At the file or package level, an overview shows entities as summary cards:

- Components and interfaces shown as boxes with name and input/output port names
- No internal graph detail at this level
- Acts as navigation index: links/drill-down into full component graphs
- Applicable to both Mermaid output (as a summary diagram or table of contents at the top of the markdown file) and interactive viewers

### Progressive Zoom Levels (Interactive Viewer)

For the interactive viewer (not Mermaid), progressive detail on zoom:

1. **Far zoom**: only entity boxes with names visible
2. **Medium zoom**: port names become visible on entity boxes
3. **Close zoom / click-to-open**: component box becomes transparent, revealing internal graph (nodes, edges, full connection detail)
4. **Hover interactions**: component description tooltip (from comments, when available), port description tooltip

## Delivery Phases

### Phase 0: Source Graph Model + Mermaid Backend

Start with the source graph model from day one. Mermaid is the first renderer consuming it.

**Source graph model:**
- Implement the projection layer from raw/source AST to `GraphDocument`
- Define stable IDs for nodes, ports, edges
- Support all source-level connection semantics: fan-in, fan-out, chained connections, struct selectors, const senders, array bypass
- Add golden tests for projection invariants using representative examples from `examples/` and `e2e/` fixtures

**Mermaid renderer:**
- Consumes the source graph model (not raw AST traversal)
- Reuse encoding/sanitization patterns from IR Mermaid backend where useful
- Output structure: **directory tree mirroring package structure**, with markdown files containing Mermaid diagrams
  - One markdown file per package (or per file -- to be determined by what works best for navigation)
  - Each markdown file contains: package/file overview (entity index with component/interface names and ports), then one Mermaid flowchart per component showing its full graph
  - Markdown files are interlinked: overview links to component sections, component nodes that reference other components link to their definitions
- Output is strictly readonly/documentational
- Add snapshot tests

**Optional semantic enrichment path:**
- When analyzer results are available, attach resolved type metadata to graph nodes (non-blocking: the model is usable from raw AST alone)

### Phase 1: VSCode Extension Integration + Mermaid Preview

- In the `vscode-neva` extension, add a command: `Neva: Preview Graph (Mermaid)`
- Preview renders the Mermaid output for the current file/package using VSCode's built-in Markdown preview capabilities
- Validates the source graph model produces useful, readable output that real users can navigate
- Extension infrastructure (command registration, LSP client wiring) is exercised early

### Phase 2: Interactive Readonly Viewer MVP

- Add a webview-based visual preview in `vscode-neva`: `Neva: Preview Mode`
- Preview renders current file in readonly visual form:
  - File-level view with component boxes (overview)
  - Click a component to see its full graph view
  - Option to open a component in focused/fullscreen view
- **Drill-down navigation**: double-clicking a node that is an instance of another component swaps the current canvas to that component's `ComponentGraph`. This should feel seamless -- the viewer fetches (or has pre-fetched) the target component's graph and transitions to it, with a breadcrumb or back-navigation to return. This "infinite drill-down" is what makes visual programming feel like navigating an operating system, not staring at a static diagram. It is the core interaction that gives the viewer real utility beyond documentation.
- No graph editing/write-back in this phase
- Technology choice for webview graph rendering to be evaluated (React Flow/XYFlow, Cytoscape.js, D3, or custom SVG -- the source graph model is renderer-agnostic, with a thin adapter per renderer)
- Source graph model is served via LSP (see Phase 3), or initially via direct projection in the extension if simpler for bootstrapping

### Phase 3: Visual LSP API

- Introduce explicit LSP methods shaped by what the interactive viewer actually needed:
  - `neva/getGraphDocument` (file or package scope)
  - `neva/getComponentGraph` (single component, optional optimization)
- Return source-graph-model-based payload with optional semantic enrichment sections
- Reuse existing LSP helpers (file resolution, location/range conversion, indexed build access)
- Add tests for request/response shape and failure behavior
- Do not depend on legacy `resolve_file` for rollout

### Phase 4: Versioning, Capabilities, and Enrichment Layers

- Explicit schema versioning in payloads (`version: 1`) with additive evolution rules
- `capabilities` flags:
  - `semanticOverlay` (analyzer-derived type information)
  - `desugaredOverlay` (deferred/disabled initially)
  - `textGraphSync` (bidirectional navigation/selection sync)
- Base payload remains usable without any overlay
- Enrichment layers from compiler stages (analyzed AST, desugared AST, IR) can optionally attach to the source graph model for advanced use cases -- the idea is that the final super-detailed view could merge information from all stages into a single rich picture

### Phase 5: Runtime Overlay + Trace Playback

**Runtime overlay:**
- Non-canonical overlay payload keyed by stable graph IDs
- Map runtime `Sent`/`Received` events to edge activity visualization
- Keep runtime overlays optional and decoupled from core AST snapshot

**Trace playback visualization:**
- The trace log (`trace.log` from `--emit-trace`) is an append-only event journal recording all message sends and receives during execution
- Visualize trace as a "movie": step forward and backward through events, seeing messages flow along edges in the graph
- This is playback only (not live debugging) -- the user can scrub through the recorded execution timeline
- Neva programs can have side effects, so this is not a time machine (you cannot re-execute from an arbitrary point), but the visual replay of the recorded trace is valuable for understanding program behavior
- Consider: could potentially work with live interception too (not just log files), showing real-time message flow during execution

### Phase 6: Standalone Parity

- Standalone app depends on Neva LSP as backend transport (not VSCode-specific)
- Reuse the same projection layer and source graph model through LSP
- Keep host-specific state outside canonical payload (window layout, preferences, shortcuts)
- Parity: both hosts render the same source program consistently from the same contract

## Future Explorations (Not In Scope, Worth Mentioning)

### Three.js as Alternative Renderer

The existing IR-level Three.js backend (`internal/compiler/backend/ir/threejs/`) generates interactive 3D HTML visualizations with force-directed layout, orbit controls, and animated connections. This is significant prior art.

Future possibility: a source-level Three.js renderer consuming the same source graph model. Unlike Mermaid (directory of markdown files for filesystem-based navigation) or the webview viewer (panel inside IDE), a Three.js renderer would produce a single HTML file representing an entire program as one navigable 3D space. This is a fundamentally different interaction paradigm -- the user doesn't think in terms of files, but navigates a spatial environment.

This could serve as:
- An alternative readonly visualization for documentation/exploration
- A potential foundation for an alternative interactive editor (if the 3D paradigm proves compelling)
- An export format for sharing program visualizations outside the IDE

Deferred until after the core Mermaid and interactive viewer are stable, but the source graph model should be designed to support this renderer without schema changes.

### Custom File Renderer / Hybrid Text-Visual Mode

VSCode's extension API supports custom editors and notebook-like renderers. A compelling future direction:

- A `.neva` file rendered as a sequence of sections, where each section is either text (connection network, declarations) or a visual graph (rendered component)
- Toggle between text and visual representation per component within the same view
- Both representations could eventually be editable, with changes syncing bidirectionally

This is distinct from notebooks: Neva files are not notebooks (no cell execution model, no interleaved output). But the visual pattern of mixing text blocks and rendered visual blocks in a single document view is similar and could reuse some of the same VSCode APIs (custom editor provider).

Important VSCode platform constraints to keep in mind:

- A Custom Editor is a webview running inside an iframe. It does **not** inherit LSP features (hover, autocomplete, diagnostics squiggles, semantic tokens) that the standard text editor gets for free. Any such features needed in the visual view must be reimplemented in the webview.
- AI inline autocomplete (e.g., Copilot's "ghost text") is tied to the active text editor's cursor and text buffer. If a Custom Editor completely hides the text, inline completion goes blind -- there is no standard cursor in a text buffer for it to attach to.
- AI chat tools (Cursor Composer, Copilot Chat, etc.) are unaffected: they can read the backing `.neva` file from the workspace regardless of what the user is looking at. The visual mode does not interfere with file-level AI assistance.
- The real risk with a full Custom Editor is not technical interference but **information overload**: presenting too many modes, too many toggles, and too many representations at once. Simplicity of the user's mental model matters more than feature density.

Key open questions:
- Is the file the right unit, or should the visual mode abstract away the filesystem?
- How does bidirectional text <-> visual sync work at the editing level?
- What is the transition model between text editing and visual editing within one component?

### Structured Documentation Comments

Comments in `.neva` files are currently parsed but not stored in the AST. Future work:

- Define a structured comment format (similar to JSDoc/Javadoc) for documenting components and ports
- Store parsed documentation in the AST
- Surface in visual editor: hover over a component to see its full description, hover over a port to see its specific documentation
- Starting point: just component-level descriptions, then port-level

### Debugger Protocol Integration

A separate initiative exists for implementing VSCode's Debug Adapter Protocol (DAP) for Neva. Until that work is done, the exact debugging UX in the visual editor is an open question. The trace playback visualization (Phase 5) provides debugging value without requiring DAP, but a full debugger would enable:

- Breakpoints on nodes/edges in the visual graph
- Step-through execution with live state inspection
- Variable/message inspection on ports

This is intentionally deferred and will be designed when the DAP work progresses.

### Full Enrichment View (All Compiler Stages Merged)

The ultimate vision for a feature-rich visual debugging experience: overlay information from every compiler stage onto the source graph:

- Raw AST: the user's code as written (base layer, always present)
- Analyzed AST: resolved types, validated semantics
- Desugared AST: synthetic nodes inserted by the compiler (fan-in nodes, field selector nodes, const emitter nodes) -- shown as ghost/secondary nodes distinct from user-written nodes
- IR: lowered representation showing the actual runtime graph
- Runtime trace: recorded message flow

Each layer toggleable. This gives a complete picture from "what I wrote" through "what the compiler sees" to "what actually happened at runtime."

## IDE Integration Model

### Committed MVP Direction

**Side panel preview** (webview synchronized with text editor), activated via command (`Neva: Preview Mode`). This is the proven pattern used by Markdown Preview, PlantUML, and similar extensions. It:

- Does not take ownership of the document lifecycle
- Supports bidirectional navigation (click code -> highlight in graph, click graph -> reveal in code)
- Coexists naturally with the text editor
- Does not require custom editor provider complexity

### Note: "Tethered UI" Concept

An idea worth keeping in mind (attributed to a Google Gemini 3 Pro review of this plan): highlighting a line of code in the text editor should make the corresponding node or edge in the visual preview "glow," and clicking a node in the graph should scroll the text editor to that node's definition. This "tethered UI" feel is what builds user trust in the visual model -- the two representations are visibly, continuously linked, not just navigable on explicit action. This is not a planned deliverable yet but captures the aspiration for text-graph synchronization quality.

### Post-MVP Patterns to Evaluate

Once the side-panel preview is proven:

1. **Custom editor mode**: full-screen visual editor replacing the text view for a component (toggle-able)
2. **Code lens / inline affordances**: small visual previews or "open graph" buttons inline in the text editor
3. **Selection sync lifecycle**: cursor position in text maps to highlighted element in graph, and vice versa
4. **Focus/drill-down model**: how clicking a node that represents another component transitions to that component's graph (inline expand vs. navigation vs. secondary panel)

## Risks and Mitigations

- **Contract drift as language evolves**
  - Mitigation: versioned schema + golden tests against real example programs
- **Graph density in larger components**
  - Mitigation: component-first scope + progressive zoom levels + overview/index at package level
- **Host coupling**
  - Mitigation: thin host adapters over shared projection layer and source graph model
- **Mermaid limitations for validating interactive UX**
  - Mitigation: Mermaid validates the data model and AST mapping, not the interactive experience. Move to interactive viewer (Phase 2) promptly.
- **Source graph model drifting toward renderer-specific concerns**
  - Risk: since Mermaid is the first consumer, the model could unconsciously absorb Mermaid idioms (subgraph nesting, arrow styles, node shape keywords). Later renderers (React-based viewer, Three.js) would then fight the model instead of consuming it naturally.
  - Mitigation: the model must describe Neva semantics only (components, nodes, ports, edges, fan-in, fan-out, literals, array bypass). Anything renderer-specific (Mermaid syntax encoding, React Flow handle positions, Three.js mesh types) belongs strictly in the renderer adapter layer. Code review and golden tests should enforce this boundary.
- **Layout quality in interactive viewer**
  - Mitigation: source graph model should convey port ordering (canonical order derived from AST declaration order) and graph topology in a layout-algorithm-friendly format. Exact layout engine choice is a Phase 2 decision.
- **Custom Editor mode losing LSP and AI autocomplete capabilities**
  - Risk: a Custom Editor (webview iframe) does not inherit LSP features (hover, autocomplete, diagnostics) or AI inline completion (Copilot ghost text) from the standard text editor. Moving to a full visual mode could feel like a capability downgrade.
  - Mitigation: the side-panel preview MVP preserves the text editor and all its capabilities. Custom Editor mode (if pursued post-MVP) should be a toggle, not a replacement -- the user can always switch back to text. AI chat tools (Composer, Copilot Chat) remain fully functional regardless of visual mode since they read the workspace files directly.
- **Information overload from too many modes and representations**
  - Risk: offering text editing, side-panel preview, custom editor, zoom levels, enrichment overlays, and trace playback simultaneously could overwhelm users rather than help them.
  - Mitigation: progressive disclosure -- start with the simplest mode (side-panel preview) and let users opt into richer modes. Each mode should feel like a natural extension, not a separate tool to learn.

## Non-Goals (Initial Phases)

- Direct graph editing / source write-back
- Showing desugared synthetic nodes as primary representation
- Locking to a single UI rendering stack before the source graph model stabilizes
- Full debugger integration (separate DAP initiative)
- Abstracting away the filesystem in the visual editor (deferred open question)

## Success Criteria

- Source graph model correctly represents all source-level Neva connection semantics (fan-in, fan-out, chained connections, struct selectors, const senders, array bypass)
- Mermaid backend generates stable, readable, interlinked markdown with component graphs
- One canonical source graph model consumed by Mermaid renderer and later interactive viewer
- Stable graph identities across refreshes for unchanged code
- Reliable bidirectional navigation between graph elements and source locations
- Interactive viewer (Phase 2) provides useful readonly visualization that real users find valuable for understanding Neva programs

## Dependency Summary

- **LSP baseline**: core language features already merged, visual-specific API (`neva/get*` methods) is the remaining wiring
- **VSCode extension**: `vscode-neva` repository exists with LSP client and TextMate grammar; new visual features are added there
- **Trace log**: `--emit-trace` flag and `DebugInterceptor` already implemented in runtime
- **Analyzer library mode**: indexer can analyze all packages without entry point, supporting workspace-wide visualization
- **Debugger protocol (DAP)**: separate initiative, not a blocker for any phase in this plan
