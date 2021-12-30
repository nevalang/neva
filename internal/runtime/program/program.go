package program

type Program struct {
	Nodes       map[string]Node
	Connections []Connection
	StartPort   PortAddr
}

type Node struct {
	IO       IO
	Type     NodeType
	Const    Const
	Operator OpRef
}

type IO struct {
	In, Out map[string]PortMeta
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

type Const struct {
	Int int64
}

type OpRef struct {
	Pkg, Name string
}

type Connection struct {
	From PortAddr
	To   []PortAddr
}

type PortAddr struct {
	// Type PortAddrType
	Node, Port string
	Idx        uint8
}

type PortAddrType uint8

const (
	PortTypeNorm PortAddrType = iota + 1
	ArrPortBypass
)

type Type uint8

const (
	IntType Type = iota + 1
	StrType
	BoolType
	SigType
)
