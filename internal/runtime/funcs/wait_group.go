package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type waitGroup struct{}

func (g waitGroup) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	addIn, err := io.In.Port("add")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.Port("sig")
	if err != nil {
		return nil, err
	}

	return g.Handle(addIn, sigOut), nil
}

func (g waitGroup) Handle(
	addIn,
	sigOut chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		var (
			mu    sync.RWMutex
			ready bool  // ready when we have received a nonnegative count, protected by mu.
			count int64 // internal count, protected by mu.
		)

		reset := func() {
			mu.Lock()
			defer mu.Unlock()
			ready = false
			count = 0
		}

		add := func(n int64) {
			mu.Lock()
			defer mu.Unlock()
			if n >= 0 {
				ready = true
			}
			count += n
		}

		done := func() bool {
			mu.RLock()
			defer mu.RUnlock()
			return ready && count <= 0
		}

		update := func() {
			for {
				select {
				case <-ctx.Done():
					return
				case n := <-addIn:
					add(n.Int())
				}
			}
		}
		go update()

		for {
			if !done() {
				select {
				case <-ctx.Done():
					return
				default:
				}
				continue
			}
			select {
			case <-ctx.Done():
				return
			case sigOut <- nil:
				reset()
			}
		}
	}
}
