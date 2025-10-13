package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type accumulator struct{}

func (a accumulator) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	initIn, err := io.In.Single("init")
	if err != nil {
		return nil, err
	}

	updIn, err := io.In.Single("upd")
	if err != nil {
		return nil, err
	}

	lastIn, err := io.In.Single("last")
	if err != nil {
		return nil, err
	}

	curOut, err := io.Out.Single("cur")
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
				acc  runtime.Msg
				last = false
			)

			initMsg, initOk := initIn.Receive(ctx)
			if !initOk {
				return
			}

			if !curOut.Send(ctx, initMsg) {
				return
			}

			acc = initMsg

			for !last {
				var dataMsg, lastMsg runtime.Msg
				var dataOk, lastOk bool

				var wg sync.WaitGroup

				wg.Go(func() {
					dataMsg, dataOk = updIn.Receive(ctx)
				})

				wg.Go(func() {
					lastMsg, lastOk = lastIn.Receive(ctx)
				})

				wg.Wait()

				if !dataOk || !lastOk {
					return
				}

				if !curOut.Send(ctx, dataMsg) {
					return
				}

				acc = dataMsg
				last = lastMsg.Bool()
			}

			if !resOut.Send(ctx, acc) {
				return
			}
		}
	}, nil
}
