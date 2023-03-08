package ir

type (
	Program struct {
		Start    PortAddr
		Ports    map[PortAddr]uint8 // Value describes buffer size
		Routines Routines
		Net      []Connection
	}

	PortAddr struct {
		Path string
		Port string
		Idx  uint8
	}

	Routines struct {
		Void      []PortAddr // TODO spawn 1 goroutine and make a lockfree round-robin
		Giver     map[PortAddr]Msg // TODO use msg refs to avoid duplications?
		Component []ComponentRef
	}

	Connection struct {
		SenderSide    ConnectionSide
		ReceiverSides []ConnectionSide
	}

	ConnectionSide struct {
		PortAddr        PortAddr
		SelectorActions []SelectorAction // len(0) means no action
	}

	SelectorAction struct {
		RecField string // "" means use ArrIdx
		ArrIdx   int
	}

	ComponentRef struct {
		Pkg, Name string
		PortAddrs CmpPortAddrs
	}

	Msg struct {
		Type  MsgType
		Bool  bool
		Int   int
		Float float64
		Str   string
		Vec   []Msg
		Map   map[string]Msg
	}

	MsgType uint8

	CmpPortAddrs struct {
		In, Out []PortAddr
	}
)

const (
	BoolMsg MsgType = iota + 1
	IntMsg
	FloatMsg
	StrMsg
	MapMsg
	Vec
)
