package ir

type (
	Program struct {
		Ports       map[PortAddr]uint8
		Routines    Routines
		Connections []Connection
	}

	PortAddr struct {
		Path, Port string // TODO rename Port to Name
		Idx        uint8
	}

	Routines struct {
		Giver     map[PortAddr]Msg
		Component []ComponentRef // Rename to Components
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

	ComponentRef struct { // TODO rename to ComponentRoutine
		Pkg, Name string
		PortAddrs ComponentPortAddrs
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

	ComponentPortAddrs struct {
		In, Out []PortAddr
	}
)

const (
	IntMsg MsgType = iota + 1
)
