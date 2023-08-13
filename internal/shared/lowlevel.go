package shared

// IDEA use json/protobuf strings for unmarshalling

type LLProgram struct {
	Ports map[LLPortAddr]uint8 // uint8 for slots usage
	Net   []LLConnection
	Funcs []LLFunc
}

type LLPortAddr struct {
	Path string // Path to node
	Port string // Name of the port
	Idx  uint8  // Slot index of the port
}

type LLConnection struct {
	SenderSide    LLPortAddr
	ReceiverSides []LLReceiverConnectionSide
}

type LLReceiverConnectionSide struct {
	PortAddr  LLPortAddr
	Selectors []string
}

// LLFunc is a instantiation object that runtime will use to spawn a function
type LLFunc struct {
	Ref    LLFuncRef // runtime will use this reference to find the function to spawn
	IO     LLFuncIO  // this is the ports function will use to receive and send data
	Params LLMsg     // function can receive predefined message at instantiation time
}

type LLFuncRef struct {
	Pkg, Name string
}

type LLFuncIO struct {
	In, Out []LLPortAddr
}

type LLMsg struct {
	Type  LLMsgType
	Bool  bool
	Int   int
	Float float64
	Str   string
	Vec   []string          // ordered list of msg refs
	Map   map[string]string // key -> msg ref
}

type LLMsgType uint8

const (
	LLBoolMsg LLMsgType = iota + 1
	LLIntMsg
	LLFloatMsg
	LLStrMsg
	LLVecMsg
	LLMapMsg
)
