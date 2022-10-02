package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/pkg/initutils"
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
	initutils.NilPanic(decoder, builder, executor)

	return Runtime{
		decoder:  decoder,
		builder:  builder,
		executor: executor,
	}
}
