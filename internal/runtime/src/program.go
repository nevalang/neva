package src

type (
	Program struct {
		Ports       map[AbsolutePortAddr]uint8 // Ports maps address to buffer size
		Connections []Connection
		Effects     Effects
		StartPort   AbsolutePortAddr
	}

	AbsolutePortAddr struct {
		Path string
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
		DictReadingPath []string // Only used for DictKeyReading
	}

	ReceiverConnectionPointType uint8

	Effects struct {
		Operators []OperatorEffect
		Constants map[AbsolutePortAddr]Msg
		Triggers  []TriggerEffect
	}

	OperatorEffect struct {
		Ref       OperatorRef
		PortAddrs OperatorPortAddrs
	}

	TriggerEffect struct {
		Msg                     Msg
		InPortAddr, OutPortAddr AbsolutePortAddr
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
	DictReading
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	StructMsg
)
