package runtime

import (
	"context"

	"github.com/emil14/neva/internal/runtime/core"
	"github.com/emil14/neva/internal/runtime/src"
)

type Decoder interface {
	Decode([]byte) (src.Program, error)
}

type Executor interface {
	Exec(context.Context, Program) error
}

type (
	Builder interface {
		Build(src.Program) (Program, error)
	}

	Program struct {
		Start src.Ports
		Ports Ports
		Net   []Connection
		Nodes Nodes
	}

	Ports map[src.Ports]chan core.Msg

	Connection struct {
		Src       src.Connection
		Sender    chan core.Msg
		Receivers []chan core.Msg
	}

	Nodes struct {
		Const     []ConstNode
		Trigger   []TriggerNode
		Component []ComponentNode
		Void      []chan core.Msg
	}

	ConstNode struct {
		OutPort chan core.Msg
		Msg     core.Msg
	}

	TriggerNode struct {
		InPort  chan core.Msg
		OutPort chan core.Msg
		Msg     core.Msg
	}

	ComponentNode struct {
		Ref src.FuncRef
		IO  core.IO
	}
)
