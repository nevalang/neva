package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime/src"
)

type Runtime struct {
	decoder  Decoder
	builder  Builder
	executor Executor
}

type (
	Decoder interface {
		Decode([]byte) (src.Program, error)
	}

	Builder interface {
		Build(src.Program) (Build, error)
	}
	Build struct {
		StartPort   src.AbsPortAddr
		Ports       Ports
		Connections []Connection
		Effects     Effects
	}
	Ports      map[src.AbsPortAddr]chan core.Msg
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

	Executor interface {
		Exec(context.Context, Build) error
	}
)

var (
	ErrDecoder     = errors.New("decoder")
	ErrNoStartPort = errors.New("no start port")
	ErrBuilder     = errors.New("builder")
	ErrExecutor    = errors.New("executor")
)

func (r Runtime) Run(ctx context.Context, bb []byte) error {
	prog, err := r.decoder.Decode(bb)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	if _, ok := prog.Ports[prog.StartPort]; !ok {
		return fmt.Errorf("%w: %v", ErrNoStartPort, prog.StartPort)
	}

	build, err := r.builder.Build(prog)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuilder, err)
	}

	return fmt.Errorf("%w: %v", ErrExecutor, r.executor.Exec(ctx, build))
}

func MustNew(
	decoder Decoder,
	builder Builder,
	executor Executor,
) Runtime {
	utils.NilPanic(decoder, builder, executor)

	return Runtime{
		decoder:  decoder,
		builder:  builder,
		executor: executor,
	}
}
