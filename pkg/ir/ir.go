package ir

// Program represents the main structure containing ports, connections, and funcs.
type Program struct {
	Ports       []*PortInfo
	Connections []*Connection
	Funcs       []*Func
}

// PortInfo contains information about each port.
type PortInfo struct {
	PortAddr *PortAddr
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
	SenderSide    *PortAddr
	ReceiverSides []*ReceiverConnectionSide
}

// ReceiverConnectionSide represents the receiver side of a connection.
type ReceiverConnectionSide struct {
	PortAddr *PortAddr
}

// Func represents a function within the program.
type Func struct {
	Ref string
	Io  *FuncIO
	Msg *Msg
}

// FuncIO represents the input/output ports of a function.
type FuncIO struct {
	Inports  []*PortAddr
	Outports []*PortAddr
}

// Msg represents a message.
type Msg struct {
	Type  MsgType
	Bool  bool
	Int   int64
	Float float64
	Str   string
	List  []*Msg
	Map   map[string]*Msg
}

// MsgType is an enumeration of message types.
type MsgType int32

const (
	MSG_TYPE_UNSPECIFIED MsgType = 0
	MSG_TYPE_BOOL        MsgType = 1
	MSG_TYPE_INT         MsgType = 2
	MSG_TYPE_FLOAT       MsgType = 3
	MSG_TYPE_STR         MsgType = 4
	MSG_TYPE_LIST        MsgType = 5
	MSG_TYPE_MAP         MsgType = 6
)
