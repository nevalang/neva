package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime/src"
	"github.com/emil14/neva/pkg/tools"
)

type Runtime struct {
	decoder  Decoder
	builder  Builder
	executor Executor
}

var (
	ErrDecoder     = errors.New("decoder")
	ErrNoStartPort = errors.New("no start port")
	ErrBuilder     = errors.New("builder")
	ErrExecutor    = errors.New("executor")
)

func (r Runtime) Run(ctx context.Context, prog src.Program) error {
	if _, ok := prog.Ports[prog.StartPortAddr]; !ok {
		return fmt.Errorf("%w: %v", ErrNoStartPort, prog.StartPortAddr)
	}

	build, err := r.builder.Build(prog)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuilder, err)
	}

	return fmt.Errorf(
		"%w: %v",
		ErrExecutor,
		r.executor.Exec(ctx, build),
	)
}

func MustNew(
	decoder Decoder,
	builder Builder,
	executor Executor,
) Runtime {
	tools.NilPanic(decoder, builder, executor)

	return Runtime{
		decoder:  decoder,
		builder:  builder,
		executor: executor,
	}
}
