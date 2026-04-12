package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listToStream struct{}

func (c listToStream) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			list := data.List()

			for idx := range list {
				item := streamItem(
					list[idx],
					int64(idx),
					idx == len(list)-1,
				)

				if !resOut.Send(ctx, item) {
					return
				}
			}
		}
	}, nil
}
