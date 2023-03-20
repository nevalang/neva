package ir

type (
	Program struct {
		Ports       map[PortAddr]uint8
		Routines    Routines
		Connections []Connection
	}

	PortAddr struct {
		Path, Name string
		Idx        uint8
	}

	Routines struct {
		Giver map[PortAddr]Msg
		Func  []FuncRoutine
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

	FuncRoutine struct {
		Ref FuncRef
		IO  FuncIO
	}

	FuncRef struct {
		Pkg, Name string
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

	FuncIO struct {
		In, Out []PortAddr
	}
)

const (
	IntMsg MsgType = iota + 1
)
