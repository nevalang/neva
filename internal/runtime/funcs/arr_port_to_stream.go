package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToStream struct{}

func (arrayPortToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, err := arrayIn(io, "port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	// TODO: could be optimized by using portIn.ReceiveAll()
	// but we need to handle order of sending messages to stream
	return func(ctx context.Context) {
		//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		l := portIn.Len()

		for {
			for idx := range l {
				msg, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}

				item := streamItem(
					msg,
					int64(idx),
					idx == l-1,
				)

				if !resOut.Send(ctx, item) {
					return
				}
			}
		}
	}, nil
}
