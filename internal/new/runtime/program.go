package runtime

type (
	Program struct {
		Nodes       map[string]Node
		Connections []Connection
		StartPort   PortAddr
	}

	Connection struct {
		From PortAddr
		To   []PortAddr
	}

	PortAddr struct {
		Node, Port string
		Idx        uint8
	}

	Node struct {
		Type  NodeType
		IO    NodeIO
		OpRef OpRef
		Const map[string]Msg
	}

	NodeType uint8

	NodeIO struct {
		In, Out map[string]PortMeta
	}

	PortMeta struct {
		Slots, Buf uint8
	}

	OpRef struct {
		Pkg, Name string
	}

	Msg struct {
		Type MsgType
		Int  int
		Str  string
		Bool bool
	}

	MsgType uint8
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	SigMsg
)

const (
	ModuleNode NodeType = iota + 1
	OperatorNode
	ConstNode
)
