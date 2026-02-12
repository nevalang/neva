# Visual Editor for Neva AST-Level Program Representation

## Current State

The compiler already has:

- **AST** (`internal/compiler/ast/flowast.go`) with full JSON tags for LSP consumption
- **IR visualization backends**: DOT, Mermaid, ThreeJS (in `internal/compiler/backend/ir/`)
- **LSP server** (`cmd/lsp/`) with a partially defined `resolve_file` custom method and `GetFileViewResponse` type
- **Desugarer** that transforms syntactic sugar into canonical nodes/connections

Key insight from AGENTS.md: "Avoid composite/control-flow syntax to preserve 1:1 visual/text graph mapping."

## AST-to-Visual Mapping Model

The core question: how do Neva AST abstractions map to visual node editor elements?

### Canvas Scope

Each canvas represents a **single Component** (a `def` block). This is the natural unit because:

- A component is self-contained: it has an interface (IO), nodes, and a network
- It maps to a single graph (nodes + edges)
- It aligns with the "block diagram" concept from LabVIEW and other tools
- Navigation between components works via drill-down (double-click a node to open its component definition)

A package or file view becomes a **list/tree of components**, not a canvas.

### Element Mapping

```
AST Concept          -> Visual Element
----------------------------------------------
Component            -> Canvas (the graph)
Component.Interface  -> "Self" boundary (input/output port strip on canvas edges)
  IO.In ports        -> Left-edge port handles (incoming data into the component)
  IO.Out ports       -> Right-edge port handles (outgoing data from the component)
Node                 -> Box/Card in the canvas
  Node.EntityRef     -> Box label (e.g. "split: strings.Split")
  Node.TypeArgs      -> Badge/subtitle on box (e.g. "<int>")
  Node.ErrGuard      -> Visual indicator (e.g. "?" icon)
  Node.DIArgs        -> Collapsible sub-section or secondary ports
Connection           -> Edge/Wire between ports
  Sender.PortAddr    -> Edge source handle
  Receiver.PortAddr  -> Edge target handle
  Sender.Const       -> Inline constant badge on edge or mini-node
  StructSelector     -> Small intermediary badge on edge (e.g. ".body")
Port (on a Node)     -> Handle dot on left/right side of box
  Port.IsArray       -> Multiple indexed handles or array indicator
  Port.TypeExpr      -> Tooltip or type badge
Const (literal)      -> Small distinct node shape or inline annotation
Import               -> Not shown on canvas (available via node palette)
TypeDef              -> Not shown on canvas (tooltip on port hover)
```

### Node Anatomy (Detail)

Each node in the canvas is a card with:

- **Header**: `localName: pkg.ComponentName<TypeArgs>` (e.g., `split: strings.Split`)
- **Left side**: Input port handles (labeled, colored by type)
- **Right side**: Output port handles (labeled, colored by type)
- **Error guard indicator**: `?` badge if `ErrGuard == true`
- **Overload indicator**: If `OverloadIndex != nil`

Special cases:

- **Self ports** (`:portName`): Rendered as port strips on the canvas boundary, not as regular nodes
- **Constants/literals** in connections: Small inline nodes or edge labels (like Unreal Blueprint's "literal" pins)
- **Fan-in/Fan-out**: Shown as visual merging/splitting of wires (NOT as synthetic `FanIn`/`FanOut` nodes -- those are desugarer artifacts)

### Connection Rendering

Connections map directly from `[]Connection` in the AST Component:

- Single sender -> single receiver: simple edge
- Multiple senders (fan-in `[a, b] -> c`): visual merge point
- Multiple receivers (fan-out `a -> [b, c]`): visual split point
- Chained connections (`a -> b -> c`): shown as sequential edges
- Struct selectors (`.field`): small label on the edge or intermediary diamond
- Constant senders (`42 -> node:port`): small constant pill attached to edge origin

### What is NOT on the Canvas

- Imports (belong to file-level metadata, shown in a separate panel or palette)
- Type definitions (shown as tooltips or in a types panel)
- Constants (only shown when used as connection senders)
- Interface entities (shown in a separate view or as component signature)

## Architecture: Layered Separation

```
                    +-----------------------+
                    |   VSCode Extension    |
                    | (or standalone shell) |
                    +---------+-------------+
                              |
                    +---------v-------------+
                    |   Visual Editor UI    |
                    |  (React + xyflow)     |
                    |  - Renders graph      |
                    |  - Pan/zoom/select    |
                    |  - Node palette       |
                    +---------+-------------+
                              | postMessage / JSON
                    +---------v-------------+
                    |  Editor Protocol      |
                    |  (JSON-based API)     |
                    |  Request/Response     |
                    +---------+-------------+
                              |
                    +---------v-------------+
                    |   LSP Server          |
                    |  (Go, extends LSP)    |
                    |  - resolve_file       |
                    |  - resolve_component  |
                    +---------+-------------+
                              |
                    +---------v-------------+
                    |   Compiler Frontend   |
                    |  - Parser -> AST      |
                    |  - Analyzer           |
                    +-----------------------+
```

Key principles:

- **Compiler knows nothing about the editor**. It produces AST with JSON tags.
- **LSP server** is the bridge: it runs the compiler frontend, holds the analyzed AST, and serves it to the editor via custom JSON-RPC methods.
- **Editor Protocol** is a JSON message format between the webview and the LSP. This protocol is editor-agnostic.
- **Visual Editor UI** is a React app using `@xyflow/react` (React Flow). It can run inside a VSCode webview or as a standalone web app.

### AST Backend (New Compiler Component)

By analogy with existing IR backends (DOT, Mermaid, ThreeJS), we introduce an **AST-level visualization backend**. However, unlike IR backends that flatten everything, this operates at the source-code level:

- **Input**: Analyzed `ast.Build` (post-analyzer, pre-desugarer)
- **Output**: A JSON format optimized for the visual editor (nodes with positions, edges with routing, resolved port interfaces)
- **Location**: New package, e.g. `internal/compiler/sourcecode/visualize/` or simply served by LSP

The existing `GetFileViewResponse` type in `cmd/lsp/server/get_file_view.go` already outlines this:

```go
type GetFileViewResponce struct {
    File  src.File
    Extra Extra // info not in the file but needed for rendering
}

type Extra struct {
    NodesPorts map[string]map[string]src.Interface // flows -> nodes -> interface
}
```

This needs extension to also include:

- Resolved type information for ports (for type-based coloring and tooltips)
- Node positions (if stored; initially auto-layout)
- Component signatures for referenced nodes (for rendering port handles)

### What the LSP Needs to Serve

Custom JSON-RPC methods (extending existing `resolve_file`):

1. **`neva/getComponentView`** - Returns a component's full visual representation:
   - All nodes with their resolved interfaces (port names and types)
   - All connections
   - Component's own interface (self-ports)
   - Position hints (if available)

2. **`neva/getPackageOutline`** - Returns list of entities in a package (for navigation/palette)

3. **`neva/getNodeSignature`** - Returns resolved interface for a node reference (for rendering ports when dragging from palette -- future)

## Technology Choices

- **Graph library**: `@xyflow/react` (React Flow 12) -- MIT license, 35k+ stars, active development, used by Stripe/Typeform, supports custom nodes/edges, minimap, controls, elkjs layout integration
- **Layout engine**: `elkjs` (Eclipse Layout Kernel) -- best automatic layout for dataflow graphs, supports layered/hierarchical layout, port-aware
- **UI framework**: React 18+ with TypeScript
- **Bundler**: Vite (fast HMR, good webview compat)
- **Styling**: CSS modules or Tailwind (TBD)
- **Communication**: VSCode webview `postMessage` API wrapping JSON-RPC to LSP

## Phased Roadmap

### Phase 0: AST Visualization Protocol (Compiler Side)

- Implement `neva/getComponentView` in the LSP server
- Use analyzed AST (post-analyzer, pre-desugarer) as source
- Resolve node interfaces (port names + types) for all nodes in a component
- Return JSON with: component interface, nodes (with resolved ports), connections
- Test with the existing fizzbuzz and advanced_error_handling examples

### Phase 1: Readonly Viewer MVP (VSCode Extension)

- Scaffold React + xyflow webview inside VSCode extension
- Custom node component: shows name, ref, input/output port handles
- Custom edge component: shows connection type, handles fan-in/fan-out visually
- Self-port strips on canvas edges for component IO
- Auto-layout via elkjs (layered left-to-right)
- Sync with active `.neva` file (open component view when file changes)
- Pan, zoom, minimap, select node (highlight corresponding text)

### Phase 2: Enhanced Readonly Features

- Type-based port coloring (different colors for int, string, any, error, etc.)
- Constant/literal sender visualization (inline pills or small nodes)
- Struct selector visualization (`.field` labels on edges)
- Fan-in/fan-out visual merge/split points
- Error guard `?` visual indicator
- Tooltips on ports (type info), nodes (component signature)
- Click node to jump to text definition (bidirectional navigation)
- Node palette panel (list available components from imports)

### Phase 3: Interactive Editing (Future)

- Drag nodes to set positions (stored as metadata)
- Add nodes from palette
- Draw connections by dragging between ports
- Delete nodes/connections
- Sync changes back to `.neva` text files via LSP
- Undo/redo

### Phase 4: Standalone App (Future)

- Extract React app from webview into standalone Electron/Tauri shell
- Communicate with LSP server over stdio/TCP
- Same visual editor, different host

## Reference Analysis Summary

| Tool                    | Key Lesson for Neva                                                                                             |
| ----------------------- | --------------------------------------------------------------------------------------------------------------- |
| **Unreal Blueprints**   | Pin-to-pin connections with type coloring, drag-to-empty creates node, context menu shows compatible nodes only |
| **LabVIEW**             | Dual panel (front panel + block diagram), component IO as boundary ports                                        |
| **n8n / Node-RED**      | Clean node cards with input/output handles, subflow collapse, drag-on-wire to insert                            |
| **React Flow / xyflow** | Best library for custom node editors in React, supports all required features                                   |
| **Pipe**                | Comments as grey boxes on canvas, compact multi-input representation                                            |
| **Flyde**               | Closest reference -- VSCode extension with React Flow webview, custom editor API                                |
| **Mermaid (existing)**  | ELK layout works well for Neva dataflow graphs (see issue #970 screenshot)                                      |

## Key Design Decisions

1. **Pre-desugarer AST** as visualization source (users see their code, not compiler artifacts)
2. **Component = Canvas** (not file or package)
3. **Self-ports as boundary strips** (like LabVIEW front panel terminals)
4. **Fan-in/fan-out as visual merges** (not synthetic nodes)
5. **Readonly first** (viewer before editor)
6. **xyflow for rendering** (best React graph library, elkjs for layout)
7. **LSP as bridge** (compiler stays editor-agnostic)
8. **JSON protocol** (webview-agnostic, works in VSCode and standalone)

## Appendix: AST Structure Reference

### Component (the core visual unit)

```go
type Component struct {
    Interface  Interface              // Component interface (IO ports)
    Directives map[Directive]string   // Compiler directives
    Nodes      map[string]Node        // Component nodes (instances)
    Net        []Connection           // Network of connections
    Meta       core.Meta              // Source location
}
```

### Node (box on canvas)

```go
type Node struct {
    Directives    map[Directive]string // Compiler directives
    EntityRef     core.EntityRef       // Reference to component/interface
    TypeArgs      TypeArgs             // Type arguments for generics
    ErrGuard      bool                 // Whether node uses `?` operator
    DIArgs        map[string]Node      // Dependency injection arguments
    OverloadIndex *int                 // Index for overloaded components
    Meta          core.Meta            // Source location
}
```

### Connection (edge on canvas)

```go
type Connection struct {
    Senders   []ConnectionSender   // Sending ports/constants
    Receivers []ConnectionReceiver // Receiving ports
    Meta      core.Meta            // Source location
}

type ConnectionSender struct {
    PortAddr       *PortAddr  // Port address (if sending from port)
    Const          *Const     // Constant value (if sending constant)
    StructSelector []string   // Field selector path for struct access
    Meta           core.Meta  // Source location
}

type ConnectionReceiver struct {
    PortAddr          *PortAddr   // Port address (if receiving to port)
    ChainedConnection *Connection // Chained connection (if forwarding)
    Meta              core.Meta   // Source location
}
```

### IO and Ports

```go
type IO struct {
    In   map[string]Port // Input ports (keyed by port name)
    Out  map[string]Port // Output ports (keyed by port name)
    Meta core.Meta
}

type Port struct {
    TypeExpr ts.Expr   // Type expression
    IsArray  bool      // Whether port is an array
    Meta     core.Meta
}

type PortAddr struct {
    Node string      // Node name (empty for self ports)
    Port string      // Port name
    Idx  *uint8      // Array index (nil if not array, 255 for bypass)
    Meta core.Meta
}
```

### Build Hierarchy

```
Build
  └── Modules (map[ModuleRef]Module)
      └── Module
          ├── Manifest (ModuleManifest)
          └── Packages (map[string]Package)
              └── Package (map[string]File)
                  └── File
                      ├── Imports (map[string]Import)
                      └── Entities (map[string]Entity)
                          └── Entity.Component []Component
```

## Appendix: Desugarer Transformations (What Visual Editor Hides)

The desugarer creates virtual nodes that the visual editor should NOT show directly:

- `__const__{N}` -- virtual constants for literal senders
- `__new__{N}` -- `New<T>` nodes for constant emitters
- `__newv2__{N}` -- `NewV2<T>` nodes for chained constant senders
- `__field__{N}` -- `Field<T>` nodes for struct selectors
- `__fan_in__{N}` -- `FanIn` nodes for multi-sender connections
- `__fan_out__{N}` -- `FanOut` nodes for multi-receiver connections
- `__del__` -- `Del` node for unused outports

The visual editor works with pre-desugarer AST to show the user's original code structure.
