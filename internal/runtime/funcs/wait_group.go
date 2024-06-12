package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type waitGroup struct{}

func (g waitGroup) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	countIn, err := io.In.SingleInport("count")
	if err != nil {
		return nil, err
	}

	sigIn, err := io.In.SingleInport("sig")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.SingleOutport("sig")
	if err != nil {
		return nil, err
	}

	return g.Handle(countIn, sigIn, sigOut), nil
}

func (waitGroup) Handle(
	countIn,
	sigIn,
	sigOut chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			var wg sync.WaitGroup
			var count int64

			select {
			case n := <-countIn:
				count = n.Int()
				wg.Add(int(count))
			case <-ctx.Done():
				return
			}

			go func() {
				for i := int64(0); i < count; i++ {
					select {
					case <-sigIn:
						wg.Done()
					case <-ctx.Done():
						return
					}
				}
			}()
			wg.Wait()

			select {
			case sigOut <- nil:
			case <-ctx.Done():
				return
			}
		}
	}
}
