// Package runtime implements environment for dataflow programs execution.
package runtime

import (
	"context"
	"errors"
	"sync"
)

type Runtime struct {
	funcRunner FuncRunner
}

func (r Runtime) Run(ctx context.Context, prog Program) error {
	stop, ok := prog.Ports[PortAddr{Path: "out", Port: "stop"}]
	if !ok {
		return errors.New("stop outport not found")
	}

	funcRun, err := r.funcRunner.Run(prog.Funcs)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	cancelableCtx, cancel := context.WithCancel(ctx)

	go func() {
		<-stop
		cancel()
	}()

	qch := make(chan QueueItem)

	queue := NewQueue(
		qch,
		prog.Connections,
		prog.Ports,
	)

	go func() {
		queue.Run(cancelableCtx)
		wg.Done()
	}()

	qch <- QueueItem{
		Sender: PortAddr{Path: "in", Port: "start"},
		Msg:    &baseMsg{},
	}

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

	wg.Wait()

	return nil
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
