package program

type Program struct {
	Nodes map[string]Node
	Net   []Connection
	IORef IORef
}

type IORef struct {
	In, Out []FullPortAddr
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
	From FullPortAddr
	To   []FullPortAddr
}

type FullPortAddr struct {
	Node string
	Port string
	Slot uint8
}
