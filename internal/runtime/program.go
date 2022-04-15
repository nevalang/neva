package runtime

type (
	Program struct {
		Ports       []PortAddr
		Connections []Connection // replace with map? (avoid possible duplicates)
		Effects     Effects
		StartPort   PortAddr
	}

	PortAddr struct {
		Path string // IDEA: rename to Node for consistency with compiler?
		Name string
		Idx  uint8
	}

	Connection struct {
		From PortAddr
		To   []PortAddr
	}

	Effects struct {
		Ops   []Operator
		Const map[PortAddr]ConstMsg
	}

	Operator struct {
		Ref       OpRef
		PortAddrs OpPortAddrs
	}

	ConstMsg struct {
		Type MsgType
		Int  int
		Str  string
		Bool bool
	}

	OpRef struct {
		Pkg, Name string
	}

	OpPortAddrs struct {
		In, Out []PortAddr
	}

	MsgType uint8
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	SigMsg
)
