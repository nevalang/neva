package runtime

import (
	"context"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime/src"
)

type Decoder interface {
	Decode([]byte) (src.Program, error)
}

type (
	Builder interface {
		Build(src.Program) (Build, error)
	}

	Build struct {
		StartPort   src.AbsPortAddr
		Ports       Ports
		Connections []Connection
		Effects     Effects
	}

	Ports map[src.AbsPortAddr]chan core.Msg

	Connection struct {
		Src       src.Connection
		Sender    chan core.Msg
		Receivers []chan core.Msg
	}

	Effects struct {
		Constants []ConstantEffect
		Operators []OperatorEffect
		Triggers  []TriggerEffect
	}

	ConstantEffect struct {
		OutPort chan core.Msg
		Msg     core.Msg
	}

	OperatorEffect struct {
		Ref src.OperatorRef
		IO  core.IO
	}

	TriggerEffect struct {
		InPort  chan core.Msg
		OutPort chan core.Msg
		Msg     core.Msg
	}
)

type Executor interface {
	Exec(context.Context, Build) error
}
