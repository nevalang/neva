package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type NodeMeta struct {
	Ref string
	Msg string
}

type Node struct {
	Path  string
	Label string
	In    map[string]struct{}
	Out   map[string]struct{}
	Meta  NodeMeta
}

type Conn struct {
	From ir.PortAddr
	To   ir.PortAddr
}

type Report struct {
	ModName     string
	CompilerVer string
	Nodes       map[string]*Node
	NodePaths   []string
	Conns       []Conn
}

func Build(prog *ir.Program) Report {
	modName, compilerVer := parseCommentMetadata(prog.Comment)

	nodes := map[string]*Node{}

	getOrCreateNode := func(path string) *Node {
		if n, ok := nodes[path]; ok {
			return n
		}
		parts := strings.Split(path, "/")
		label := parts[len(parts)-1]
		n := &Node{
			Path:  path,
			Label: label,
			In:    make(map[string]struct{}),
			Out:   make(map[string]struct{}),
		}
		nodes[path] = n
		return n
	}

	processPortAddr := func(addr ir.PortAddr, isOut bool) {
		if addr.Path == "" {
			path := "global_" + addr.Port
			n := getOrCreateNode(path)
			n.Label = addr.Port
			return
		}

		path := trimInOutSuffix(addr.Path)
		n := getOrCreateNode(path)
		if isOut {
			n.Out[addr.Port] = struct{}{}
		} else {
			n.In[addr.Port] = struct{}{}
		}
	}

	for sender, receiver := range prog.Connections {
		processPortAddr(sender, true)
		processPortAddr(receiver, false)
	}

	for _, f := range prog.Funcs {
		matchPath := matchFuncPath(f)
		if matchPath == "" {
			continue
		}
		if n, ok := nodes[matchPath]; ok {
			n.Meta.Ref = f.Ref
			if f.Msg != nil {
				n.Meta.Msg = formatMsg(f.Msg)
			}
		}
	}

	nodePaths := make([]string, 0, len(nodes))
	for p := range nodes {
		nodePaths = append(nodePaths, p)
	}
	sort.Strings(nodePaths)

	conns := make([]Conn, 0, len(prog.Connections))
	for s, r := range prog.Connections {
		conns = append(conns, Conn{From: s, To: r})
	}
	sort.Slice(conns, func(i, j int) bool {
		s1 := conns[i].From.String()
		s2 := conns[j].From.String()
		if s1 != s2 {
			return s1 < s2
		}
		return conns[i].To.String() < conns[j].To.String()
	})

	return Report{
		ModName:     modName,
		CompilerVer: compilerVer,
		Nodes:       nodes,
		NodePaths:   nodePaths,
		Conns:       conns,
	}
}

func WriteHeader(w io.Writer, r Report) error {
	if _, err := fmt.Fprintf(w, "# Program: %s\n\n", r.ModName); err != nil {
		return err
	}
	if r.CompilerVer != "" {
		if _, err := fmt.Fprintf(w, "**Compiler:** %s\n\n", r.CompilerVer); err != nil {
			return err
		}
	}
	return nil
}

func WriteComponentsTable(w io.Writer, r Report) error {
	if _, err := fmt.Fprintln(w, "## 2. Components"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Node | Ref | Config | Ports |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| :--- | :--- | :--- | :--- |"); err != nil {
		return err
	}

	for _, path := range r.NodePaths {
		n := r.Nodes[path]
		portList := SortedPorts(n)

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

		if _, err := fmt.Fprintf(w, "| `%s` | %s | %s | `%s` |\n", n.Label, ref, msg, strings.Join(portList, ", ")); err != nil {
			return err
		}
	}

	return nil
}

func WriteMetrics(w io.Writer, r Report) error {
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "## 3. Metrics"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "* **Nodes:** %d\n", len(r.Nodes)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "* **Connections:** %d\n", len(r.Conns)); err != nil {
		return err
	}
	return nil
}

func SortedPorts(n *Node) []string {
	portList := make([]string, 0, len(n.In)+len(n.Out))
	for p := range n.In {
		portList = append(portList, "in:"+p)
	}
	for p := range n.Out {
		portList = append(portList, "out:"+p)
	}
	sort.Strings(portList)
	return portList
}

func parseCommentMetadata(comment string) (string, string) {
	var modName, compilerVer string
	if strings.HasPrefix(comment, "//") {
		parts := strings.Fields(comment)
		for _, p := range parts {
			if strings.HasPrefix(p, "main=") {
				modName = strings.TrimPrefix(p, "main=")
			}
			if strings.HasPrefix(p, "compiler=") {
				compilerVer = strings.TrimPrefix(p, "compiler=")
			}
		}
	}
	return modName, compilerVer
}

func trimInOutSuffix(path string) string {
	if strings.HasSuffix(path, "/in") {
		return strings.TrimSuffix(path, "/in")
	}
	if strings.HasSuffix(path, "/out") {
		return strings.TrimSuffix(path, "/out")
	}
	return path
}

func matchFuncPath(f ir.FuncCall) string {
	for _, addr := range f.IO.In {
		path := trimInOutSuffix(addr.Path)
		if path != "" {
			return path
		}
	}
	for _, addr := range f.IO.Out {
		path := trimInOutSuffix(addr.Path)
		if path != "" {
			return path
		}
	}
	return ""
}

func formatMsg(msg *ir.Message) string {
	switch msg.Type {
	case ir.MsgTypeString:
		return fmt.Sprintf("%q", msg.String)
	case ir.MsgTypeInt:
		return fmt.Sprintf("%d", msg.Int)
	case ir.MsgTypeBool:
		return fmt.Sprintf("%v", msg.Bool)
	case ir.MsgTypeFloat:
		return fmt.Sprintf("%f", msg.Float)
	case ir.MsgTypeList:
		return "[...]"
	case ir.MsgTypeDict, ir.MsgTypeStruct, ir.MsgTypeUnion:
		return "{...}"
	default:
		return "-"
	}
}
