package runtime

import "github.com/emil14/neva/internal/core"

type (
	Program struct {
		Ports       []PortAddr
		Connections []Connection
		Effects     Effects
		StartPort   PortAddr
	}

	PortAddr struct {
		Path string
		Name string
		Idx  uint8
	}

	Connection struct {
		Sender    PortAddr
		Receivers []ConnectionPoint
	}

	ConnectionPoint struct {
		PortAddr        PortAddr
		Type            ConnectionPointType
		StructFieldPath []string
	}

	ConnectionPointType uint8

	Effects struct {
		Ops   []Operator
		Const map[PortAddr]ConstMsg
	}

	Operator struct {
		Ref       OpRef
		PortAddrs OpPortAddrs
	}

	ConstMsg struct {
		Type    core.Type
		BoolMsg core.BoolMsg
		IntMsg  core.IntMsg
		StrMsg  core.StrMsg
	}

	OpRef struct {
		Pkg, Name string
	}

	OpPortAddrs struct {
		In, Out []PortAddr
	}
)

const (
	Normal ConnectionPointType = iota + 1
	FieldReading
)
