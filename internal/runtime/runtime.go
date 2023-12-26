// Package runtime implements environment for dataflow programs execution.
package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/runtime/errgroup"
)

type Runtime struct {
	connector  Connector
	funcRunner FuncRunner
}

var ErrNilDeps = errors.New("runtime deps nil")

func New(connector Connector, funcRunner FuncRunner) Runtime {
	return Runtime{
		connector:  connector,
		funcRunner: funcRunner,
	}
}

type (
	Func func(context.Context, FuncIO) (func(), error)
)

var (
	ErrStartPortNotFound = errors.New("enter port not found")
	ErrExitPortNotFound  = errors.New("exit port not found")
	ErrConnector         = errors.New("connector")
	ErrFuncRunner        = errors.New("func runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) (code int, err error) {
	enter := prog.Ports[PortAddr{Path: "main/in", Port: "enter"}]
	if enter == nil {
		return 0, ErrStartPortNotFound
	}

	exit := prog.Ports[PortAddr{Path: "main/out", Port: "exit"}]
	if exit == nil {
		return 0, ErrExitPortNotFound
	}

	ctx, cancel := context.WithCancel(ctx)
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(gctx, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.funcRunner.Run(gctx, prog.Funcs); err != nil {
			return fmt.Errorf("%w: %v", ErrFuncRunner, err)
		}
		return nil
	})

	go func() { // kick
		enter <- emptyMsg{}
	}()

	var exitCode int64
	go func() {
		exitCode = (<-exit).Int()
		cancel()
	}()

	return int(exitCode), g.Wait()
}
