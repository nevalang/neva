// Package runtime implements environment for dataflow programs execution.
package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	funcRun, err := r.funcRunner.Run(ctx, prog.Funcs)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFuncRunner, err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		funcRun()
		wg.Done()
	}()
	go func() {
		r.connector.Connect(ctx, prog.Connections)
		wg.Done()
	}()

	go func() {
		enter <- emptyMsg{}
	}()

	var exitCode int64
	go func() {
		exitCode = (<-exit).Int()
		cancel()
	}()

	wg.Wait() // wait for connector and funcs to finish

	return int(exitCode), nil
}
