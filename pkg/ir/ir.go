package ir

// Program represents the main structure containing ports, connections, and funcs.
type Program struct {
	Ports       []*PortInfo   `protobuf:"bytes,1,rep,name=ports" json:"ports,omitempty"`
	Connections []*Connection `protobuf:"bytes,2,rep,name=connections" json:"connections,omitempty"`
	Funcs       []*Func       `protobuf:"bytes,3,rep,name=funcs" json:"funcs,omitempty"`
}

// PortInfo contains information about each port.
type PortInfo struct {
	PortAddr *PortAddr `protobuf:"bytes,1,req,name=port_addr,json=portAddr" json:"port_addr,omitempty"`
	BufSize  uint32    `protobuf:"varint,2,opt,name=buf_size,json=bufSize" json:"buf_size,omitempty"`
}

// PortAddr represents the address of a port.
type PortAddr struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Port string `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
	Idx  uint32 `protobuf:"varint,3,opt,name=idx" json:"idx,omitempty"`
}

// Connection represents connections between ports.
type Connection struct {
	SenderSide    *PortAddr                 `protobuf:"bytes,1,req,name=sender_side,json=senderSide" json:"sender_side,omitempty"`
	ReceiverSides []*ReceiverConnectionSide `protobuf:"bytes,2,rep,name=receiver_sides,json=receiverSides" json:"receiver_sides,omitempty"`
}

// ReceiverConnectionSide represents the receiver side of a connection.
type ReceiverConnectionSide struct {
	PortAddr *PortAddr `protobuf:"bytes,1,req,name=port_addr,json=portAddr" json:"port_addr,omitempty"`
}

// Func represents a function within the program.
type Func struct {
	Ref string  `protobuf:"bytes,1,opt,name=ref" json:"ref,omitempty"`
	Io  *FuncIO `protobuf:"bytes,2,opt,name=io" json:"io,omitempty"`
	Msg *Msg    `protobuf:"bytes,3,opt,name=msg" json:"msg,omitempty"`
}

// FuncIO represents the input/output ports of a function.
type FuncIO struct {
	Inports  []*PortAddr `protobuf:"bytes,1,rep,name=inports" json:"inports,omitempty"`
	Outports []*PortAddr `protobuf:"bytes,2,rep,name=outports" json:"outports,omitempty"`
}

// Msg represents a message.
type Msg struct {
	Type  MsgType         `protobuf:"varint,1,opt,name=type,enum=MsgType" json:"type,omitempty"`
	Bool  bool            `protobuf:"varint,2,opt,name=bool" json:"bool,omitempty"`
	Int   int64           `protobuf:"varint,3,opt,name=int" json:"int,omitempty"`
	Float float64         `protobuf:"fixed64,4,opt,name=float" json:"float,omitempty"`
	Str   string          `protobuf:"bytes,5,opt,name=str" json:"str,omitempty"`
	List  []*Msg          `protobuf:"bytes,6,rep,name=list" json:"list,omitempty"`
	Map   map[string]*Msg `protobuf:"bytes,7,rep,name=map" json:"map,omitempty"`
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
