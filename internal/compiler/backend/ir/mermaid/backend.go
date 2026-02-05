package mermaid

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler/backend/ir/report"
	"github.com/nevalang/neva/internal/compiler/ir"
)

// Encoder implements the logic to write Mermaid diagram to a writer
type Encoder struct{}

//nolint:gocyclo // Encoding includes multiple rendering branches.
func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	rep := report.Build(prog)
	if err := report.WriteHeader(w, rep); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "## 1. Visual Flow"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "```mermaid"); err != nil {
		return err
	}
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
	nodes := rep.Nodes
	nodePaths := rep.NodePaths

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
		ins := make([]string, 0, len(n.In))
		outs := make([]string, 0, len(n.Out))
		for p := range n.In {
			ins = append(ins, p)
		}
		for p := range n.Out {
			outs = append(outs, p)
		}
		sort.Strings(ins)
		sort.Strings(outs)

		for _, p := range ins {
			pid := nodeID + "__" + sanitize(p)
			fmt.Fprintf(w, "        %s(%s):::port\n", pid, formatPortLabel(p))
		}

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
	for _, c := range rep.Conns {
		fromID := getPortID(c.From)
		toID := getPortID(c.To)
		fmt.Fprintf(w, "    %s --> %s\n", fromID, toID)
	}

	fmt.Fprintln(w, "```")

	if err := report.WriteComponentsTable(w, rep); err != nil {
		return err
	}
	return report.WriteMetrics(w, rep)
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

func EncodeFlowchart(prog *ir.Program) (string, error) {
	var buf bytes.Buffer

	rep := report.Build(prog)
	if _, err := fmt.Fprintln(&buf, "flowchart TD"); err != nil {
		return "", err
	}

	styles := `    %% Style definitions
    classDef port fill:#fff,stroke:#000,stroke-width:1px,min-width:0,padding:2px,rx:5,ry:5;
    classDef node fill:#f5f5f5,stroke:#333,stroke-width:1px,color:#333;
    classDef invisible display:none;
`
	if _, err := fmt.Fprint(&buf, styles); err != nil {
		return "", err
	}

	nodes := rep.Nodes
	for _, path := range rep.NodePaths {
		n := nodes[path]
		nodeID := sanitize(n.Path)

		if strings.HasPrefix(path, "global_") {
			fmt.Fprintf(&buf, "    %s([%s]):::port\n", nodeID, n.Label)
			continue
		}

		fmt.Fprintf(&buf, "    subgraph %s [\"%s\"]\n", nodeID, n.Label)
		fmt.Fprintln(&buf, "        direction TB")

		ins := make([]string, 0, len(n.In))
		outs := make([]string, 0, len(n.Out))
		for p := range n.In {
			ins = append(ins, p)
		}
		for p := range n.Out {
			outs = append(outs, p)
		}
		sort.Strings(ins)
		sort.Strings(outs)

		for _, p := range ins {
			pid := nodeID + "__" + sanitize(p)
			fmt.Fprintf(&buf, "        %s(%s):::port\n", pid, formatPortLabel(p))
		}
		for _, p := range outs {
			pid := nodeID + "__" + sanitize(p)
			fmt.Fprintf(&buf, "        %s(%s):::port\n", pid, formatPortLabel(p))
		}

		fmt.Fprintln(&buf, "    end")
		fmt.Fprintf(&buf, "    class %s node\n", nodeID)
		fmt.Fprintln(&buf, "")
	}

	for _, c := range rep.Conns {
		fromID := getPortID(c.From)
		toID := getPortID(c.To)
		fmt.Fprintf(&buf, "    %s --> %s\n", fromID, toID)
	}

	return buf.String(), nil
}
