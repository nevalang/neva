package src

type (
	Program struct {
		Ports       Ports
		Connections []Connection
		Effects     Effects
		StartPort   AbsPortAddr
	}

	Ports map[AbsPortAddr]uint8 // Ports maps address to buffer size

	AbsPortAddr struct {
		Path string
		Port string
		Idx  uint8
	}

	Connection struct {
		SenderPortAddr            AbsPortAddr
		ReceiversConnectionPoints []ReceiverConnectionPoint
	}

	ReceiverConnectionPoint struct {
		PortAddr        AbsPortAddr
		Type            ReceiverConnectionPointType
		DictReadingPath []string // Only used for DictKeyReading
	}

	ReceiverConnectionPointType uint8

	Effects struct {
		Operators []OperatorEffect
		Constants map[AbsPortAddr]Msg
		Triggers  []TriggerEffect
	}

	OperatorEffect struct {
		Ref       OperatorRef
		PortAddrs OperatorPortAddrs
	}

	TriggerEffect struct {
		Msg                     Msg
		InPortAddr, OutPortAddr AbsPortAddr
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
		In, Out []AbsPortAddr
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
