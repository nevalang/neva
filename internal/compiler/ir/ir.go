package ir

type (
	Program struct {
		Ports       map[PortAddr]uint8
		Routines    Routines
		Connections []Connection
	}

	PortAddr struct {
		Path, Port string
		Idx        uint8
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
		Selectors []Selector
	}

	Selector struct {
		RecField string
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
