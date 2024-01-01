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
	Func func(FuncIO, Msg) (func(context.Context), error)
)

var (
	ErrStartPortNotFound = errors.New("enter port not found")
	ErrExitPortNotFound  = errors.New("exit port not found")
	ErrConnector         = errors.New("connector")
	ErrFuncRunner        = errors.New("func runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) (err error) {
	enter := prog.Ports[PortAddr{Path: "in", Port: "enter"}]
	if enter == nil {
		return ErrStartPortNotFound
	}

	exit := prog.Ports[PortAddr{Path: "out", Port: "exit"}]
	if exit == nil {
		return ErrExitPortNotFound
	}

	funcRun, err := r.funcRunner.Run(prog.Funcs)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFuncRunner, err)
	}

	ctx2, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		funcRun(ctx2)
		wg.Done()
	}()
	go func() {
		r.connector.Connect(ctx2, prog.Connections)
		wg.Done()
	}()

	go func() {
		enter <- emptyMsg{}
	}()

	go func() {
		<-exit
		cancel()
	}()

	wg.Wait()

	return nil
}
