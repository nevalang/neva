package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	decoder Decoder
	builder Builder
	// executor Executor
	connector Connector
	effector  Effector
}

var (
	ErrDecoder          = errors.New("decoder")
	ErrNoStartPort      = errors.New("no start port")
	ErrBuilder          = errors.New("builder")
	ErrConnector        = errors.New("connector")
	ErrEffector         = errors.New("effector")
	ErrStartPortBlocked = errors.New("start port blocked")
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

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(gctx, build.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.effector.Effect(gctx, build.Effects); err != nil {
			return fmt.Errorf("%w: %v", ErrEffector, err)
		}
		return nil
	})

	select {
	case <-time.After(time.Second):
		return errors.New("timeout")
	case build.Ports[prog.StartPort] <- core.NewDictMsg(nil):
		return g.Wait()
	}
}

func MustNew(
	decoder Decoder,
	builder Builder,
	connector Connector,
	effector Effector,
) Runtime {
	utils.NilPanic(decoder, effector, connector)

	return Runtime{
		decoder:   decoder,
		builder:   builder,
		effector:  effector,
		connector: connector,
	}
}
