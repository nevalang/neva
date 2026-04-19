package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type accumulator struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (a accumulator) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	initIn, err := singleIn(io, "init")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	updIn, err := singleIn(io, "upd")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	lastIn, err := singleIn(io, "last")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	curOut, err := singleOut(io, "cur")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
