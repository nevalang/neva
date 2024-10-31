package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsGreaterOrEqual struct{}

func (intIsGreaterOrEqual) Create(
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

			if !resOut.Send(ctx, runtime.NewBoolMsg(accMsg.Int() >= elMsg.Int())) {
				return
			}
		}
	}, nil
} 