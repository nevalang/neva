package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type ternarySelector struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p ternarySelector) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifIn, err := singleIn(io, "if")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenIn, err := singleIn(io, "then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elseIn, err := singleIn(io, "else")
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
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			dataMsg, ok := ifIn.Receive(ctx)
			if !ok {
				return
			}

			thenMsg, ok := thenIn.Receive(ctx)
			if !ok {
				return
			}

			elseMsg, ok := elseIn.Receive(ctx)
			if !ok {
				return
			}

			var resMsg runtime.Msg
			if dataMsg.Bool() {
				resMsg = thenMsg
			} else {
				resMsg = elseMsg
			}

			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
