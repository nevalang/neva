package runtime

type (
	Program struct {
		Ports       []AbsolutePortAddr
		Connections []Connection
		Effects     Effects
		StartPort   AbsolutePortAddr
	}

	AbsolutePortAddr struct {
		Path string // node path? context? scope?
		Port string
		Idx  uint8
	}

	Connection struct {
		SenderPortAddr            AbsolutePortAddr
		ReceiversConnectionPoints []ReceiverConnectionPoint
	}

	ReceiverConnectionPoint struct {
		PortAddr        AbsolutePortAddr
		Type            ReceiverConnectionPointType
		DictReadingPath []string // Only used for Type == DictReadingPath
	}

	ReceiverConnectionPointType uint8

	Effects struct {
		Operators []Operator
		Constants map[AbsolutePortAddr]Msg
	}

	Operator struct {
		Ref       OperatorRef
		PortAddrs OperatorPortAddrs
	}

	Msg struct {
		Type   MsgType
		Bool   bool
		Int    int
		Str    string
		Struct map[string]Msg
	}

	MsgType uint8

	OperatorRef struct {
		Pkg, Name string
	}

	OperatorPortAddrs struct {
		In, Out []AbsolutePortAddr
	}
)

const (
	Normal ReceiverConnectionPointType = iota + 1
	DictKeyReading
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	StructMsg
)
