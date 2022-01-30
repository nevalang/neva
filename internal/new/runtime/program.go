package runtime

type (
	Program struct {
		Nodes       map[string]Node
		Connections []Connection
		StartPort   AbsPortAddr
	}

	Node struct {
		Type        NodeType
		IO          IO
		OperatorRef OperatorRef
		ConstOut    map[string]ConstValue
	}

	Connection struct {
		From AbsPortAddr
		To   []AbsPortAddr
	}

	AbsPortAddr struct {
		Node string
		Port string
		Idx  uint8
	}

	NodeType uint8

	IO struct {
		In, Out map[RelPortAddr]Port
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

	RelPortAddr struct {
		Port string
		Idx  uint8
	}

	Port struct {
		ArrSize, Buf uint8
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
