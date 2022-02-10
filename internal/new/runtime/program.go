package runtime

type (
	Program struct {
		Ports       map[FullPortAddr]Port
		Connections []Connection
		Effects     Effects
		StartPort   FullPortAddr
	}

	FullPortAddr struct {
		Path string
		Port string
		Idx  uint8
	}

	Port struct {
		Buf uint8
	}

	Connection struct {
		From FullPortAddr
		To   []FullPortAddr
	}

	Effects struct {
		Ops   []Operator
		Const map[FullPortAddr]ConstMsg
	}

	Operator struct {
		Ref OperatorRef
		IO  OperatorIO
	}

	ConstMsg struct {
		Type MsgType
		Int  int
		Str  string
		Bool bool
	}

	OperatorRef struct {
		Pkg, Name string
	}

	OperatorIO struct {
		In, Out []FullPortAddr
	}

	MsgType uint8
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	SigMsg
)
