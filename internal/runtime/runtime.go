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

func (p *Runtime) Run(ctx context.Context, prog Program) error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-prog.Stop
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	net := NewNetwork(prog.Connections)
	go func() {
		net.Run(ctx)
		wg.Done()
	}()

	funcrun, err := p.funcRunner.Run(prog.Funcs)
	if err != nil {
		return err
	}

	go func() {
		ctx = context.WithValue(ctx, "cancel", cancel)
		funcrun(ctx)
		wg.Done()
	}()

	prog.Start <- IndexedMsg{
		index: counter.Add(1),
		data:  &baseMsg{},
	}

	wg.Wait()

	return nil
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
