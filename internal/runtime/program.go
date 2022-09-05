package runtime

type (
	Program struct {
		Ports       []PortAddr
		Connections []Connection
		Effects     Effects
		StartPort   PortAddr
	}

	PortAddr struct {
		Path string
		Name string
		Idx  uint8
	}

	Connection struct {
		Sender    PortAddr
		Receivers []ConnectionPoint
	}

	ConnectionPoint struct {
		PortAddr        PortAddr
		Type            ConnectionPointType
		StructFieldPath []string
	}

	ConnectionPointType uint8

	Effects struct {
		Ops   []Operator
		Const map[PortAddr]Msg
	}

	Operator struct {
		Ref       OpRef
		PortAddrs OpPortAddrs
	}

	Msg struct {
		Type   MsgType
		Bool   bool
		Int    int
		Str    string
		Struct map[string]Msg
	}

	MsgType uint8

	OpRef struct {
		Pkg, Name string
	}

	OpPortAddrs struct {
		In, Out []PortAddr
	}
)

const (
	Normal ConnectionPointType = iota + 1
	FieldReading
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	StructMsg
)
