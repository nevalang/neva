package ir

// Program represents the main structure containing ports, connections, and funcs.
type Program struct {
	Ports       []PortInfo
	Connections []Connection
	Funcs       []Func
}

// PortInfo contains information about each port.
type PortInfo struct {
	PortAddr PortAddr
	BufSize  uint32
}

// PortAddr represents the address of a port.
type PortAddr struct {
	Path string
	Port string
	Idx  uint32
}

// Connection represents connections between ports.
type Connection struct {
	SenderSide    PortAddr
	ReceiverSides []ReceiverConnectionSide
}

// ReceiverConnectionSide represents the receiver side of a connection.
type ReceiverConnectionSide struct {
	PortAddr PortAddr
}

// Func represents a function within the program.
type Func struct {
	Ref string
	IO  FuncIO
	Msg *Msg
}

// FuncIO represents the input/output ports of a function.
type FuncIO struct {
	In  []PortAddr
	Out []PortAddr
}

// Msg represents a message.
type Msg struct {
	Type  MsgType
	Bool  bool
	Int   int64
	Float float64
	Str   string
	List  []Msg
	Map   map[string]Msg
}

// MsgType is an enumeration of message types.
type MsgType int32

const (
	MsgTypeUnspecified MsgType = 0
	MsgTypeBool        MsgType = 1
	MsgTypeInt         MsgType = 2
	MsgTypeFloat       MsgType = 3
	MsgTypeString      MsgType = 4
	MsgTypeList        MsgType = 5
	MsgTypeMap         MsgType = 6
)
