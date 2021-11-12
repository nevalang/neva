package program

type Program struct {
	IO    IO
	Nodes []Node
	Net   []Connection
}

type IO struct {
	In, Out []AbsPortAddr
}

type Node struct {
	Path     []string
	In, Out  map[string]PortMeta
	Type     NodeType
	Const    map[string]Const
	Operator Operator
}

type PortMeta struct {
	Slots uint8
	Buf   uint8
}

type NodeType uint8

const (
	ConstNode NodeType = iota + 1
	OperatorNode
)

type Const struct {
	Type     Type
	IntValue int
}

type Operator struct {
	Pkg, Name string
}

type Connection struct {
	From AbsPortAddr
	To   []AbsPortAddr
}

type AbsPortAddr struct {
	NodePath []string
	Port     string
	Slot     uint8
}
