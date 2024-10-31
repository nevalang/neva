package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type intPow struct{}

func (intPow) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	accIn, err := io.In.Single("acc")
	if err != nil {
		return nil, err
	}

	elIn, err := io.In.Single("el")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var accMsg, elMsg runtime.Msg
			var accOk, elOk bool

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				accMsg, accOk = accIn.Receive(ctx)
			}()

			go func() {
				defer wg.Done()
				elMsg, elOk = elIn.Receive(ctx)
			}()

			wg.Wait()

			if !accOk || !elOk {
				return
			}

			// Calculate power using repeated multiplication
			base := accMsg.Int()
			exp := elMsg.Int()
			result := int64(1)
			
			for i := int64(0); i < exp; i++ {
				result *= base
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(result)) {
				return
			}
		}
	}, nil
}
