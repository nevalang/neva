package runtime

import "fmt"

type (
	Program struct {
		Ports       []AbsolutePortAddr
		Connections []Connection
		Effects     Effects
		StartPort   AbsolutePortAddr
	}

	AbsolutePortAddr struct {
		Path string // node path? context? scope?
		Port string
		Idx  uint8
	}

	Connection struct {
		SenderPortAddr            AbsolutePortAddr
		ReceiversConnectionPoints []ReceiverConnectionPoint
	}

	ReceiverConnectionPoint struct {
		PortAddr        AbsolutePortAddr
		Type            ConnectionPointType
		StructFieldPath []string // Only used for Type == StructFieldReading
	}

	ConnectionPointType uint8

	Effects struct {
		Operators []Operator
		Constants map[AbsolutePortAddr]Msg
	}

	Operator struct {
		Ref       OperatorRef
		PortAddrs OperatorPortAddrs
	}

	Msg struct {
		Type   MsgType
		Bool   bool
		Int    int
		Str    string
		Struct map[string]Msg
	}

	MsgType uint8

	OperatorRef struct {
		Pkg, Name string
	}

	OperatorPortAddrs struct {
		In, Out []AbsolutePortAddr
	}
)

const (
	Normal ConnectionPointType = iota + 1
	StructFieldReading
)

const (
	IntMsg MsgType = iota + 1
	StrMsg
	BoolMsg
	StructMsg
)

func (a AbsolutePortAddr) String() string {
	return fmt.Sprintf("%s.%s[%d]", a.Path, a.Port, a.Idx)
}
