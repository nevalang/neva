package mermaid

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
)

// Encoder implements the logic to write Mermaid diagram to a writer
type Encoder struct{}

//nolint:gocyclo // Encoding includes multiple rendering branches.
func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	// Parse comment for metadata (module, compiler version)
	// Format: // module=@@ main=hello_world compiler=0.34.0
	var modName, compilerVer string
	if strings.HasPrefix(prog.Comment, "//") {
		parts := strings.Fields(prog.Comment)
		for _, p := range parts {
			if strings.HasPrefix(p, "main=") {
				modName = strings.TrimPrefix(p, "main=")
			}
			if strings.HasPrefix(p, "compiler=") {
				compilerVer = strings.TrimPrefix(p, "compiler=")
			}
		}
	}

	fmt.Fprintf(w, "# Program: %s\n\n", modName)
	if compilerVer != "" {
		fmt.Fprintf(w, "**Compiler:** %s\n\n", compilerVer)
	}

	fmt.Fprintln(w, "## 1. Visual Flow")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "```mermaid")

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
	type NodeMeta struct {
		Ref string
		Msg string
	}

	type NodeInfo struct {
		Path  string
		Label string
		In    map[string]struct{}
		Out   map[string]struct{}
		Meta  NodeMeta
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

	// Second pass: Extract metadata from prog.Funcs
	// Iterate over Funcs and match them to existing nodes by looking at their used ports.
	for _, f := range prog.Funcs {
		var matchPath string

		// Check Inputs
		for _, addr := range f.IO.In {
			// Try to find component path from port address
			path := addr.Path
			if strings.HasSuffix(path, "/in") {
				path = strings.TrimSuffix(path, "/in")
			} else if strings.HasSuffix(path, "/out") {
				path = strings.TrimSuffix(path, "/out")
			}
			if path != "" {
				matchPath = path
				break
			}
		}

		// Check Outputs if not found
		if matchPath == "" {
			for _, addr := range f.IO.Out {
				path := addr.Path
				if strings.HasSuffix(path, "/in") {
					path = strings.TrimSuffix(path, "/in")
				} else if strings.HasSuffix(path, "/out") {
					path = strings.TrimSuffix(path, "/out")
				}
				if path != "" {
					matchPath = path
					break
				}
			}
		}

		// If we found a path, update the node's metadata
		if matchPath != "" {
			if n, ok := nodes[matchPath]; ok {
				n.Meta.Ref = f.Ref
				if f.Msg != nil {
					// Format message based on type
					switch f.Msg.Type {
					case ir.MsgTypeString:
						n.Meta.Msg = fmt.Sprintf("%q", f.Msg.String)
					case ir.MsgTypeInt:
						n.Meta.Msg = fmt.Sprintf("%d", f.Msg.Int)
					case ir.MsgTypeBool:
						n.Meta.Msg = fmt.Sprintf("%v", f.Msg.Bool)
					case ir.MsgTypeFloat:
						n.Meta.Msg = fmt.Sprintf("%f", f.Msg.Float)
					case ir.MsgTypeList:
						n.Meta.Msg = "[...]"
					case ir.MsgTypeDict, ir.MsgTypeStruct, ir.MsgTypeUnion:
						n.Meta.Msg = "{...}"
					}
				}
			}
		}
	}

	// Sort nodes for deterministic output
	nodePaths := make([]string, 0, len(nodes))
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
	conns := make([]conn, 0, len(prog.Connections))
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

	fmt.Fprintln(w, "```")

	fmt.Fprintln(w, "## 2. Components")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "| Node | Ref | Config | Ports |")
	fmt.Fprintln(w, "| :--- | :--- | :--- | :--- |")
	for _, path := range nodePaths {
		n := nodes[path]
		var portList []string
		for p := range n.In {
			portList = append(portList, "in:"+p)
		}
		for p := range n.Out {
			portList = append(portList, "out:"+p)
		}
		sort.Strings(portList)

		ref := n.Meta.Ref
		if ref == "" {
			ref = "-"
		} else {
			ref = "`" + ref + "`"
		}

		msg := n.Meta.Msg
		if msg == "" {
			msg = "-"
		} else {
			msg = "`" + msg + "`"
		}

		fmt.Fprintf(w, "| `%s` | %s | %s | `%s` |\n", n.Label, ref, msg, strings.Join(portList, ", "))
	}

	// 4. Metrics
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "## 3. Metrics")
	fmt.Fprintf(w, "* **Nodes:** %d\n", len(nodes))
	fmt.Fprintf(w, "* **Connections:** %d\n", len(prog.Connections))

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
