// Package runtime implements environment for dataflow programs execution.
package runtime

import (
	"context"
	"sync"
)

type Runtime struct {
	queue      Queue
	funcRunner FuncRunner
}

func (r Runtime) Run(ctx context.Context, prog Program) error {
	funcRun, err := r.funcRunner.Run(prog.Funcs)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	cancelableCtx, cancel := context.WithCancel(ctx)

	go func() {
		funcRun(
			context.WithValue(
				cancelableCtx,
				"cancel", //nolint:staticcheck // SA1029
				cancel,
			),
		)
		wg.Done()
	}()

	go func() {
		r.queue.Run(cancelableCtx)
		wg.Done()
	}()

	go func() {
		prog.Ports[PortAddr{Path: "in", Port: "start"}] <- &baseMsg{}
	}()

	go func() { // normal termination
		<-prog.Ports[PortAddr{Path: "out", Port: "stop"}]
		cancel()
	}()

	wg.Wait()

	return nil
}

func New(queue Queue, funcRunner FuncRunner) Runtime {
	return Runtime{
		queue:      queue,
		funcRunner: funcRunner,
	}
}
