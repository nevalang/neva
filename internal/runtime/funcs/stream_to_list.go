package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamToList struct{}

func (s streamToList) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("data")
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
		// Fully materializes one stream batch before emitting resulting list.
		list := []runtime.Msg{}

		for {
			msg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case runtime.IsStreamOpen(msg):
				list = list[:0]
				continue
			case runtime.IsStreamData(msg):
				list = append(list, runtime.StreamDataValue(msg))
				continue
			case !runtime.IsStreamClose(msg):
				continue
			}

			if !resOut.Send(ctx, runtime.NewListMsg(list)) {
				return
			}
		}
	}, nil
}
