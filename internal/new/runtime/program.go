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
		Type        NodeType
		IO          IO
		OperatorRef OperatorRef
		Const       map[string]ConstValue
	}

	NodeType uint8

	IO struct {
		In, Out map[string]Port
	}

	Port struct {
		ArrSize, Buf uint8
	}

	OperatorRef struct {
		Pkg, Name string
	}

	ConstValue struct {
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
	SimpleNode NodeType = iota + 1
	OperatorNode
	ConstNode
)
