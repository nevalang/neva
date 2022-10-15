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
		Build(src.Program) (Executable, error)
	}

	Executable struct {
		Start src.PortAddr
		Ports Ports
		Net   []Connection
		Fx    Effects
	}

	Ports map[src.PortAddr]chan core.Msg

	Connection struct {
		Src       src.Connection
		Sender    chan core.Msg
		Receivers []chan core.Msg
	}

	Effects struct {
		Const   []ConstFx
		Trigger []TriggerFx
		Func    []FuncFx
		VoidFx  []chan core.Msg
	}

	ConstFx struct {
		OutPort chan core.Msg
		Msg     core.Msg
	}

	TriggerFx struct {
		InPort  chan core.Msg
		OutPort chan core.Msg
		Msg     core.Msg
	}

	FuncFx struct {
		Ref src.FuncRef
		IO  core.IO
	}
)

type Executor interface {
	Exec(context.Context, Executable) error
}
