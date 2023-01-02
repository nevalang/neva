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

	// there should only 2 kinds of effects: giver and operator
	Nodes struct { // should probably be "effects" because there are more nodes, not just givers and operators
		Const     []ConstNode     // should be giver
		Trigger   []TriggerNode   // should be operator
		Component []ComponentNode // should be operator too (not to be confused with compiler's operators)
		Void      []chan core.Msg // should be operator
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
