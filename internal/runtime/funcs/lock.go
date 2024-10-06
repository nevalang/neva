package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type lock struct{}

func (l lock) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var (
				wg            sync.WaitGroup
				data          runtime.Msg
				dataOk, sigOk bool
			)

			wg.Add(2)
			go func() {
				data, dataOk = dataIn.Receive(ctx)
				wg.Done()
			}()
			go func() {
				_, sigOk = sigIn.Receive(ctx)
				wg.Done()
			}()
			wg.Wait()

			if !dataOk || !sigOk {
				return
			}

			if !dataOut.Send(ctx, data) {
				return
			}
		}
	}, nil
}
