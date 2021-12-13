package program

type Program struct {
	Nodes       map[string]Node
	Connections []Connection
	RootNode    string
}

type Node struct {
	In, Out map[string]PortMeta
	Type    NodeType
	Const   map[string]ConstValue
	OpRef   OpRef
}

type PortMeta struct {
	Slots uint8
	Buf   uint8
}

type NodeType uint8

const (
	ConstNode NodeType = iota + 1
	OperatorNode
	ModuleNode
)

type ConstValue struct {
	Type     Type
	IntValue int
}

type OpRef struct {
	Pkg, Name string
}

type Connection struct {
	From PortAddr
	To   []PortAddr
}

type PortAddr struct {
	Node string
	Port string
	Slot uint8
}
