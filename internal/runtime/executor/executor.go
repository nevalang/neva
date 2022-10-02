package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/initutils"
	"github.com/emil14/neva/internal/runtime"
	"golang.org/x/sync/errgroup"
)

type Executor struct {
	connector Connector
	effector  Effector
}

type (
	Connector interface {
		Connect(context.Context, []runtime.Connection) error
	}
	Effector interface {
		Effect(context.Context, runtime.Effects) error
	}
)

var (
	ErrConnector      = errors.New("connector")
	ErrEffector       = errors.New("effector")
	ErrStartPortBlock = errors.New("start port blocked")
)

func (e Executor) Exec(ctx context.Context, build runtime.Build) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.connector.Connect(gctx, build.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.effector.Effect(gctx, build.Effects); err != nil {
			return fmt.Errorf("%w: %v", ErrEffector, err)
		}
		return nil
	})

	select {
	case <-time.After(time.Second):
		return ErrStartPortBlock
	case build.Ports[build.StartPort] <- core.NewDictMsg(nil):
		return g.Wait()
	}
}

func MustNew(effector Effector, connector Connector) Executor {
	initutils.NilPanic(effector, connector)
	return Executor{
		connector: connector,
		effector:  effector,
	}
}
