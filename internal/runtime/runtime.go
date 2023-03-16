package runtime

import (
	"context"
	"errors"
	"fmt"
)

type Runtime struct {
	connector     Connector
	routineRunner RoutineRunner
}

func NewRuntime(c Connector, r RoutineRunner) Runtime {
	return Runtime{
		connector:     c,
		routineRunner: r,
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

// Run returns exit code of the program when it terminates or non nil error if something goes wrong in the Runtime.
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
		if err := r.routineRunner.Run(gctx, prog.Routines); err != nil {
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
