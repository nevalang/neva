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

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-stop
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	queue := NewQueue(
		prog.QueueChan,
		prog.Connections,
		prog.Ports,
	)

	go func() {
		queue.Run(ctx)
		wg.Done()
	}()

	prog.QueueChan <- QueueItem{
		Sender: PortAddr{Path: "in", Port: "start"},
		Msg:    &baseMsg{},
	}

	go func() {
		funcRun(
			context.WithValue(
				ctx,
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
