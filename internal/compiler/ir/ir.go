package ir

type Program struct {
	Ports       map[PortAddr]uint8
	Routines    Routines
	Connections []Connection
}

type PortAddr struct {
	Path, Name string
	Idx        uint8
}

type Routines struct {
	Giver map[PortAddr]Msg
	Func  []FuncRoutine
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

type FuncRoutine struct {
	Ref FuncRef
	IO  FuncIO
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
	Vec   []Msg
	Map   map[string]Msg
}

type MsgType uint8

const (
	IntMsg MsgType = iota + 1
	BoolMsg
	FloatMsg
	StrMsg
	VecMsg
	MapMsg
)
