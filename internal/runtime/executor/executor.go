package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/sync/errgroup"
)

type Executor struct {
	fx  Effector
	net Connector
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

func (e Executor) Exec(ctx context.Context, build runtime.Executable) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.net.Connect(gctx, build.Net); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.fx.Effect(gctx, build.Fx); err != nil {
			return fmt.Errorf("%w: %v", ErrEffector, err)
		}
		return nil
	})

	select {
	case <-time.After(time.Second):
		return ErrStartPortBlock
	case build.Ports[build.Start] <- core.NewDictMsg(nil):
		return g.Wait()
	}
}

func MustNew(effector Effector, connector Connector) Executor {
	tools.PanicWithNil(effector, connector)
	return Executor{
		net: connector,
		fx:  effector,
	}
}
