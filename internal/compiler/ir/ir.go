package ir

type (
	Program struct {
		StartPortAddr PortAddr
		Ports         map[PortAddr]uint8 // Value describes buffer size
		Routines      Routines
		Net           []Connection
	}

	PortAddr struct {
		Path string
		Port string
		Idx  uint8
	}

	Routines struct {
		Giver     map[PortAddr]Msg[any]
		Component []ComponentRef
	}

	Connection struct {
		SenderSide    ConnectionSide
		ReceiverSides []ConnectionSide
	}

	ConnectionSide struct {
		PortAddr  PortAddr
		Selectors []Selector // len(0) means no action
	}

	Selector struct {
		RecField string // "" means use ArrIdx
		ArrIdx   int
	}

	ComponentRef struct {
		Pkg, Name string
		PortAddrs CmpPortAddrs
	}

	// K type parameter only used for maps
	Msg[K comparable] struct {
		Type  MsgType
		Bool  bool
		Int   int
		Float float64
		Str   string
		Vec   []Msg[any]
		Map   map[K]Msg[K]
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
	VecMsg
	MapMsg
)
