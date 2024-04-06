package graphviz

import (
	"fmt"
	"io"
	"strings"

	"github.com/nevalang/neva/internal/runtime/ir"
)

type Node struct {
	ir.PortAddr
}

func (n Node) Label() string { return n.Port }

func (n Node) Name() string {
	path := strings.TrimSuffix(strings.TrimSuffix(n.Path, "/in"), "/out")
	return fmt.Sprint(path, ":", n.Port)
}

type Edge struct {
	Send Node
	Recv Node
}

func (e Edge) Label() string {
	switch send, recv := e.Send.Idx, e.Recv.Idx; {
	case send == 0 && recv == 0:
		return ""
	default:
		return fmt.Sprintf("%d->%d", send, recv)
	}
}

type Cluster struct {
	Index    int
	Prefix   string
	Nodes    map[Node]struct{}
	Clusters map[string]*Cluster
}

func (c *Cluster) Label() string {
	i := strings.LastIndexByte(c.Prefix, '/')
	if i == -1 {
		return c.Prefix
	}
	return c.Prefix[i+1:]
}

type ClusterBuilder struct {
	Main   *Cluster
	Edges  []Edge
	nextId int
}

func (b *ClusterBuilder) insertClusterNode(addr ir.PortAddr) {
	if b.Main == nil {
		cluster := &Cluster{}
		b.Main = cluster
		b.nextId++
	}
	cluster := b.Main
	prefix := ""
	for path := addr.Path; ; {
		before, after, found := strings.Cut(path, "/")
		if !found {
			if cluster.Nodes == nil {
				cluster.Nodes = map[Node]struct{}{}
			}
			cluster.Nodes[Node{addr}] = struct{}{}
			break
		}
		if prefix == "" {
			prefix = before
		} else {
			prefix = prefix + "/" + before
		}
		next := cluster.Clusters[before]
		if next == nil {
			if cluster.Clusters == nil {
				cluster.Clusters = map[string]*Cluster{}
			}
			next = &Cluster{Index: b.nextId, Prefix: prefix}
			cluster.Clusters[before] = next
			b.nextId++
		}
		path = after
		cluster = next
	}
}

func (b *ClusterBuilder) InsertEdge(send, recv ir.PortAddr) {
	b.insertClusterNode(send)
	b.insertClusterNode(recv)
	b.Edges = append(b.Edges, Edge{Send: Node{send}, Recv: Node{recv}})
}

func (b *ClusterBuilder) Build(w io.Writer) error {
	fmt.Fprintln(w, "digraph G {")
	recursiveBuild(w, "  ", b.Main)
	for _, e := range b.Edges {
		fmt.Fprintf(w, "  %q -> %q", e.Send.Name(), e.Recv.Name())
		if label := e.Label(); label != "" {
			fmt.Fprintf(w, "[label = %q;]", label)
		}
		fmt.Fprintln(w, ";")
	}
	fmt.Fprintln(w, "}")
	return nil
}

func recursiveBuild(w io.Writer, indent string, c *Cluster) error {
	fprintlnIndent := func(f string, a ...any) {
		fmt.Fprintln(w, indent, fmt.Sprintf(f, a...))
	}
	fprintlnIndent("subgraph cluster_%d {", c.Index)
	fprintlnIndent("  label = \"%s\";", c.Label())
	for n := range c.Nodes {
		fprintlnIndent("  %q [label = \"%s\";];", n.Name(), n.Label())
	}
	for _, sub := range c.Clusters {
		recursiveBuild(w, indent+"  ", sub)
	}
	fprintlnIndent("}")
	return nil
}
