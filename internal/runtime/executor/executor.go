package executor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/core"
	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/sync/errgroup"
)

type Executor struct {
	effector NodeRunner
	router   Router
}

type (
	Router interface {
		Route(context.Context, []runtime.Connection) error
	}
	NodeRunner interface {
		RunNodes(context.Context, runtime.Nodes) error
	}
)

var (
	ErrConnector      = errors.New("connector")
	ErrEffector       = errors.New("effector")
	ErrStartPortBlock = errors.New("start port blocked")
)

func (e Executor) Exec(ctx context.Context, build runtime.Program) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.router.Route(gctx, build.Net); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.effector.RunNodes(gctx, build.Nodes); err != nil {
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

func MustNew(effector NodeRunner, router Router) Executor {
	tools.NilPanic(effector, router)
	return Executor{
		router:   router,
		effector: effector,
	}
}
