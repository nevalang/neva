package runtime

import (
	"context"
	"sync"
	"sync/atomic"
)

var counter atomic.Uint64

type Runtime struct {
	stop, start chan IndexedMsg
	funcRunner  FuncRunner
}

func (p *Runtime) Run(ctx context.Context, prog Program) error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-p.stop
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	net := NewNetwork(prog.Connections)
	go func() {
		net.Run(ctx)
		wg.Done()
	}()

	funcrun, _ := p.funcRunner.Run(prog.Funcs)
	go func() {
		ctx = context.WithValue(ctx, "cancel", cancel)
		funcrun(ctx)
		wg.Done()
	}()

	p.start <- IndexedMsg{
		index: counter.Add(1),
		data:  nil,
	}

	wg.Wait()

	return nil
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
