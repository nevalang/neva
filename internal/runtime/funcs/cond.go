package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type cond struct{}

//nolint:gocognit,varnamelen
func (c cond) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	ifIn, err := io.In.Single("if")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	thenOut, err := io.Out.Single("then")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	elseOut, err := io.Out.Single("else")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var dataMsg, ifMsg runtime.Msg
			var dataOk, ifOk bool

			//nolint:varnamelen
			var wg sync.WaitGroup

			wg.Go(func() {
				dataMsg, dataOk = dataIn.Receive(ctx)
			})

			wg.Go(func() {
				ifMsg, ifOk = ifIn.Receive(ctx)
			})

			wg.Wait()

			if !dataOk || !ifOk {
				return
			}

			var out runtime.SingleOutport
			if ifMsg.Bool() {
				out = thenOut
			} else {
				out = elseOut
			}

			if !out.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
