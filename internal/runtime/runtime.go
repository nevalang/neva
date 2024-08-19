package runtime

import (
	"context"
	"sync"
	"sync/atomic"
)

var counter atomic.Uint64

type Runtime struct {
	funcRunner FuncRunner
}

// Run runs the program.
func (p *Runtime) Run(ctx context.Context, prog Program) error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-prog.Stop
		cancel() // this is how we normally stop the program
	}()

	// here we create runtime function instances to call them later
	funcrun, err := p.funcRunner.Run(prog.Funcs)
	if err != nil {
		return err
	}

	// we gonna use wg to make sure we don't terminate until either
	// program is successfully executed or context is closed by some runtime function
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ctx = context.WithValue(ctx, "cancel", cancel) //nolint:staticcheck // SA1029
		funcrun(ctx)                                   // funcrun blocks until context is closed
		wg.Done()
	}()

	// block until the program starts
	// by some node actually receiving from start port
	prog.Start <- IndexedMsg{
		index: counter.Add(1),
		data:  &baseMsg{},
	}

	// basically wait for the context to be closed
	// by either writing to the stop port
	// or calling cancel func by some runtime function like `panic`
	wg.Wait()

	return nil
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
