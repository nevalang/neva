package runtime

import (
	"context"
	"errors"
	"fmt"
)

type Runtime struct {
	connector Connector
	runner    RoutineRunner
}

func NewRuntime(c Connector, r RoutineRunner) Runtime {
	return Runtime{
		connector: c,
		runner:    r,
	}
}

type (
	Connector interface {
		Connect(context.Context, []Connection) error
	}
	RoutineRunner interface {
		Run(context.Context, Routines) error
	}
)

var (
	ErrStartPortNotFound = errors.New("start port not found")
	ErrExitPortNotFound  = errors.New("exit port not found")
	ErrConnector         = errors.New("connector")
	ErrRoutineRunner     = errors.New("routine runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) (code int, err error) {
	startPort, ok := prog.Ports[PortAddr{Name: "start"}]
	if !ok {
		return 0, ErrStartPortNotFound
	}

	exitPort, ok := prog.Ports[PortAddr{Name: "exit"}]
	if !ok {
		return 0, ErrExitPortNotFound
	}

	ctx, cancel := context.WithCancel(ctx)
	g, gctx := WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(gctx, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.runner.Run(gctx, prog.Routines); err != nil {
			return fmt.Errorf("%w: %v", ErrRoutineRunner, err)
		}
		return nil
	})

	go func() { startPort <- nil }()

	var exitCode int
	go func() {
		exitCode = (<-exitPort).Int()
		cancel()
	}()

	return exitCode, g.Wait()
}
