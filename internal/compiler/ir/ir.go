package ir

type Program struct {
	Funcs []Func             // what functions to spawn and how
	Net   []Connection       // how ports are connected to each other
	Msgs  map[string]Msg     // predefined data that can be referred by funcs
	Ports map[PortAddr]uint8 // ports and their buffers size
}

type PortAddr struct {
	Path, Name string
	Idx        uint8
}

type Connection struct {
	SenderSide    ConnectionSide
	ReceiverSides []ConnectionSide
}

type ConnectionSide struct {
	PortAddr  PortAddr
	Selectors []Selector
}

type Selector struct {
	RecField string
	ArrIdx   int
}

// Func is a instantiation object that runtime will use to spawn a function
type Func struct {
	Ref FuncRef // runtime will use this reference to find the function to spawn
	IO  FuncIO  // this is the ports function will use to receive and send data
	Msg *Msg    // function can receive predefined message at instantiation time
}

type FuncRef struct {
	Pkg, Name string
}

type FuncIO struct {
	In, Out []PortAddr
}

type Msg struct {
	Type  MsgType
	Bool  bool
	Int   int
	Float float64
	Str   string
	Vec   []string          // ordered list of msg refs
	Map   map[string]string // key -> msg ref
}

type MsgType uint8

const (
	BoolMsg MsgType = iota + 1
	IntMsg
	FloatMsg
	StrMsg
	VecMsg
	MapMsg
)
