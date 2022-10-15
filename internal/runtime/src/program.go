package src

type (
	Program struct {
		Fx          Fx
		Start       PortAddr
		Ports       PortSet
		Connections []Connection
	}

	PortSet map[PortAddr]uint8 // value is buf size

	PortAddr struct {
		Path string
		Port string
		Idx  uint8
	}

	Connection struct {
		SenderSide    ConnectionSide
		ReceiverSides []ConnectionSide
	}

	ConnectionSide struct {
		PortAddr PortAddr
		Action   ConnectorAction
		Payload  ConnectorPayload
	}

	ConnectorAction uint8

	ConnectorPayload struct {
		ReadDict []string
	}

	Fx struct {
		Func    []FuncFx
		Const   map[PortAddr]Msg
		Trigger []TriggerFx
	}

	FuncFx struct {
		Ref   FuncRef
		Ports PortAddrs
	}

	TriggerFx struct {
		In  PortAddr
		Out PortAddr
		Msg Msg
	}

	Msg struct {
		Type MsgType
		Bool bool
		Int  int
		Str  string
		List []Msg
		Dict map[string]Msg
	}

	MsgType uint8

	FuncRef struct {
		Pkg, Name string
	}

	PortAddrs struct {
		In, Out []PortAddr
	}
)

const (
	Nothing ConnectorAction = iota + 1
	ReadDict
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	DictMsg
	List
)
