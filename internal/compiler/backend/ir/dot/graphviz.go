package dot

import (
	"embed"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/nevalang/neva/internal/compiler/ir"
)

//go:embed *.tmpl
var tmplFS embed.FS

type Port struct {
	ir.PortAddr
}

func trimPortPath(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(path, "/in"), "/out")
}

func (p Port) FormatName() string {
	portStr := p.Port
	switch {
	case strings.HasSuffix(p.Path, "/in"):
		portStr += "/in"
	case strings.HasSuffix(p.Path, "/out"):
		portStr += "/out"
	}
	if p.IsArray {
		portStr = fmt.Sprintf("%s/%d", portStr, p.Idx)
	}
	return strconv.Quote(portStr)
}

func (p Port) FormatLabel() string {
	escapePort := html.EscapeString(p.Port)
	if p.IsArray {
		return html.EscapeString(fmt.Sprintf("%s[%d]", p.Port, p.Idx))
	}
	return escapePort
}

func (p Port) Format() string {
	path := p.Path
	portStr := p.Port
	switch {
	case strings.HasSuffix(p.Path, "/in"):
		path = path[:len(path)-3] // Trim /in
		portStr += "/in"
	case strings.HasSuffix(p.Path, "/out"):
		path = path[:len(path)-4] // Trim /out
		portStr += "/out"
	}
	if p.IsArray {
		portStr = fmt.Sprint(portStr, "/", p.Idx)
	}
	return fmt.Sprintf("%q:%q", path, portStr)
}

//nolint:govet // fieldalignment: graphviz layout fields grouped.
type Node struct {
	Name  string
	Extra string
	In    map[Port]struct{}
	Out   map[Port]struct{}
}

func (n Node) Format() string {
	return fmt.Sprintf("%q", n.Name)
}

func (n Node) FormatLabel() string {
	i := strings.LastIndexByte(n.Name, '/')
	if i == -1 {
		return n.Name
	}
	return n.Name[i+1:]
}

type Edge struct {
	Send Port
	Recv Port
}

//nolint:govet // fieldalignment: graphviz layout fields grouped.
type Cluster struct {
	Index    int
	Prefix   string
	Nodes    map[string]*Node
	Clusters map[string]*Cluster
}

func (c *Cluster) getOrCreateClusterNode(b *ClusterBuilder, path string) *Node {
	path = trimPortPath(path)
	return c.getOrCreateClusterNodeRec(b, path, "", path)
}

func (c *Cluster) getOrCreateClusterNodeRec(b *ClusterBuilder, path, prefix, remaining string) *Node {
	before, after, found := strings.Cut(remaining, "/")
	if !found {
		if c.Nodes == nil {
			c.Nodes = map[string]*Node{}
		}
		n, ok := c.Nodes[before]
		if ok {
			return n
		}
		n = &Node{
			Name: path,
		}
		c.Nodes[before] = n
		return n
	}
	if prefix == "" {
		prefix = before
	} else {
		prefix = prefix + "/" + before
	}
	next := c.Clusters[before]
	if next == nil {
		if c.Clusters == nil {
			c.Clusters = map[string]*Cluster{}
		}
		next = &Cluster{Index: b.nextId, Prefix: prefix}
		c.Clusters[before] = next
		b.nextId++
	}
	return next.getOrCreateClusterNodeRec(b, path, prefix, after)
}

func (c *Cluster) Label() string {
	i := strings.LastIndexByte(c.Prefix, '/')
	if i == -1 {
		return c.Prefix
	}
	return c.Prefix[i+1:]
}

//nolint:govet // fieldalignment: graphviz layout fields grouped.
type ClusterBuilder struct {
	Main  *Cluster
	Edges []Edge

	nextId int
	once   sync.Once
	tmpl   *template.Template
	err    error
}

func (b *ClusterBuilder) initTemplates() {
	b.tmpl, b.err = template.New("").ParseFS(tmplFS, "*.tmpl")
}

func (b *ClusterBuilder) insertClusterNode(addr ir.PortAddr) {
	if b.Main == nil {
		cluster := &Cluster{}
		b.Main = cluster
		b.nextId++
	}
	switch n := b.Main.getOrCreateClusterNode(b, addr.Path); {
	case strings.HasSuffix(addr.Path, "/in"):
		if n.In == nil {
			n.In = map[Port]struct{}{}
		}
		n.In[Port{addr}] = struct{}{}
	case strings.HasSuffix(addr.Path, "/out"):
		if n.Out == nil {
			n.Out = map[Port]struct{}{}
		}
		n.Out[Port{addr}] = struct{}{}
	}
}

func (b *ClusterBuilder) InsertEdge(send, recv ir.PortAddr) {
	b.insertClusterNode(send)
	b.insertClusterNode(recv)
	b.Edges = append(b.Edges, Edge{Send: Port{send}, Recv: Port{recv}})
}

func (b *ClusterBuilder) Build(w io.Writer) error {
	if b.once.Do(b.initTemplates); b.err != nil {
		return b.err
	}
	return b.tmpl.ExecuteTemplate(w, "graph.dot.tmpl", b)
}
