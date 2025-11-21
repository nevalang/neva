package mermaid

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	// The IR backend wrapper handles file creation, so we might not need to write to file here directly
	// but looking at dot backend, it seems it writes to a specific file "program.dot".
	// However, the ir backend (internal/compiler/backend/ir/backend.go) seems to call a specific encoder.
	// We will follow the pattern of `encodeDOT` in ir backend which calls `cb.Build(f)`.
	// So `Emit` here might not be used directly if we integrate it into `ir` backend package's `encodeMermaid`.
	// But to stay consistent with `dot` package having its own logic, we'll implement the core logic here.
	return nil
}

// Encoder implements the logic to write Mermaid diagram to a writer
type Encoder struct{}

func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	if _, err := fmt.Fprintln(w, "flowchart TD"); err != nil {
		return err
	}

	// Styles
	styles := `    %% Style definitions
    classDef port fill:#fff,stroke:#000,stroke-width:1px,min-width:0,padding:2px,rx:5,ry:5;
    classDef node fill:#f5f5f5,stroke:#333,stroke-width:1px,color:#333;
    classDef invisible display:none;
`
	if _, err := fmt.Fprint(w, styles); err != nil {
		return err
	}

	// 1. Collect all nodes (funcs) and their ports
	// We need to group ports by function (component) to create subgraphs
	// prog.Funcs contains the list of functions used.

	// Map from func ref (or custom name) to its definition for easy access
	// In IR, FuncCalls are linear.
	// We need to handle:
	// - Top-level ports (program inputs/outputs) -> often represented as :start, :stop or paths like in:start
	// - Component nodes (from Funcs)

	// Let's process Funcs to generate subgraphs
	for _, f := range prog.Funcs {
		// We are not using this loop to generate nodes because we don't have the instance ID easily here
		// without cross-referencing connections or understanding the Ref vs Instance ID mapping better.
		// But we can use it to improve labels if we can match the Ref to the Path.
		// For now, we skip this optimization to ensure graph correctness first.
		_ = f
	}

	// We will use a ClusterBuilder-like approach but simpler.
	// We need to identify "Nodes" (components) and their "Ports".
	// A PortAddr has Path and Port.
	// Path can be complex: "main/worker/job" -> Cluster "main", Cluster "worker", Node "job".
	// Port is the specific pin: "in" or "out_res".

	// Let's extract all unique "nodes" (leaf components) and their ports from connections.
	type NodeInfo struct {
		Path  string
		Label string
		In    map[string]struct{}
		Out   map[string]struct{}
	}
	nodes := map[string]*NodeInfo{}

	getOrCreateNode := func(path string) *NodeInfo {
		if n, ok := nodes[path]; ok {
			return n
		}
		// Default label is the last part of the path
		parts := strings.Split(path, "/")
		label := parts[len(parts)-1]
		n := &NodeInfo{
			Path:  path,
			Label: label,
			In:    make(map[string]struct{}),
			Out:   make(map[string]struct{}),
		}
		nodes[path] = n
		return n
	}

	// Helper to process a PortAddr
	processPortAddr := func(addr ir.PortAddr, isOut bool) {
		// Special handling for root ports like "in:start" or "out:stop"
		if addr.Path == "" {
			n := getOrCreateNode("global_" + sanitize(addr.Port))
			n.Label = addr.Port
			return
		}

		// Trim /in and /out from path to get component path
		path := addr.Path
		if strings.HasSuffix(path, "/in") {
			path = strings.TrimSuffix(path, "/in")
		} else if strings.HasSuffix(path, "/out") {
			path = strings.TrimSuffix(path, "/out")
		}

		n := getOrCreateNode(path)
		if isOut {
			n.Out[addr.Port] = struct{}{}
		} else {
			n.In[addr.Port] = struct{}{}
		}
	}

	// First pass: Collect nodes and ports from connections
	for sender, receiver := range prog.Connections {
		processPortAddr(sender, true)
		processPortAddr(receiver, false)
	}

	// Also check Funcs to get better labels (e.g. for constants)
	// We need to map Func Ref/Name to the paths we found.
	// Use a heuristic or map if possible.
	// If we rely solely on connections, we miss unconnected nodes (less important)
	// and we miss metadata like constant values.
	//
	// ISSUE: `prog.Funcs` in IR might not align 1:1 with paths if paths are instance IDs and Funcs are defs?
	// Or Funcs is list of calls?
	//
	// Let's rely on what DOT does: DOT completely ignores `prog.Funcs`.
	// It produces a graph of connections.
	// We will stick to that for now.

	// Sort nodes for deterministic output
	var nodePaths []string
	for p := range nodes {
		nodePaths = append(nodePaths, p)
	}
	sort.Strings(nodePaths)

	// 2. Emit Nodes (Subgraphs)
	for _, path := range nodePaths {
		n := nodes[path]
		// Generate Safe ID
		nodeID := sanitize(n.Path)

		// Check if it's a "global" node (hacky detection from processPortAddr)
		if strings.HasPrefix(path, "global_") {
			fmt.Fprintf(w, "    %s([%s]):::port\n", nodeID, n.Label)
			continue
		}

		// Subgraph for Component
		fmt.Fprintf(w, "    subgraph %s [\"%s\"]\n", nodeID, n.Label)
		fmt.Fprintln(w, "        direction TB")

		// Collect and sort ports
		var ins, outs []string
		for p := range n.In {
			ins = append(ins, p)
		}
		for p := range n.Out {
			outs = append(outs, p)
		}
		sort.Strings(ins)
		sort.Strings(outs)

		// Emit Inputs
		for _, p := range ins {
			pid := nodeID + "__" + sanitize(p)
			fmt.Fprintf(w, "        %s(%s):::port\n", pid, formatPortLabel(p))
		}

		// Emit Outputs
		for _, p := range outs {
			pid := nodeID + "__" + sanitize(p)
			fmt.Fprintf(w, "        %s(%s):::port\n", pid, formatPortLabel(p))
		}

		fmt.Fprintln(w, "    end")
		fmt.Fprintf(w, "    class %s node\n", nodeID)
		fmt.Fprintln(w, "")
	}

	// 3. Emit Connections
	// We need to sort connections for deterministic output
	type conn struct {
		From ir.PortAddr
		To   ir.PortAddr
	}
	var conns []conn
	for s, r := range prog.Connections {
		conns = append(conns, conn{s, r})
	}
	sort.Slice(conns, func(i, j int) bool {
		s1 := conns[i].From.String()
		s2 := conns[j].From.String()
		if s1 != s2 {
			return s1 < s2
		}
		return conns[i].To.String() < conns[j].To.String()
	})

	for _, c := range conns {
		fromID := getPortID(c.From)
		toID := getPortID(c.To)
		fmt.Fprintf(w, "    %s --> %s\n", fromID, toID)
	}

	return nil
}

func sanitize(s string) string {
	return strings.NewReplacer(
		"/", "_",
		":", "_",
		"[", "_",
		"]", "_",
		" ", "_",
		".", "_",
		"\"", "",
		"'", "",
		"-", "_",
	).Replace(s)
}

func formatPortLabel(s string) string {
	// Simplify port label: "in:sig" -> "sig"
	if idx := strings.LastIndex(s, ":"); idx != -1 {
		return s[idx+1:]
	}
	return s
}

func getPortID(addr ir.PortAddr) string {
	if addr.Path == "" {
		return "global_" + sanitize(addr.Port)
	}

	path := addr.Path
	if strings.HasSuffix(path, "/in") {
		path = strings.TrimSuffix(path, "/in")
	} else if strings.HasSuffix(path, "/out") {
		path = strings.TrimSuffix(path, "/out")
	}

	return sanitize(path) + "__" + sanitize(addr.Port)
}
