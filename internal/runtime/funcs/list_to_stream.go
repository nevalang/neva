package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listToStream struct{}

//nolint:gocognit // Stream framing and termination handling belong to one state machine.
func (c listToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			list := data.List()
			if !resOut.Send(ctx, newStreamOpenMsg()) {
				return
			}

			for idx := range list {
				if !resOut.Send(ctx, newStreamDataMsg(list[idx])) {
					return
				}
			}

			if !resOut.Send(ctx, newStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
