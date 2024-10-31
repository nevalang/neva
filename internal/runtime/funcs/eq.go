package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type eq struct{}

func (p eq) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.Single("acc")
	if err != nil {
		return nil, err
	}

	comparedIn, err := io.In.Single("el")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var (
				wg                   sync.WaitGroup
				val1, val2           runtime.Msg
				actualOk, comparedOk bool
			)

			wg.Add(2)
			go func() {
				val1, actualOk = actualIn.Receive(ctx)
				wg.Done()
			}()
			go func() {
				val2, comparedOk = comparedIn.Receive(ctx)
				wg.Done()
			}()
			wg.Wait()

			if !actualOk || !comparedOk {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(val1.Equal(val2)),
			) {
				return
			}
		}
	}, nil
}
