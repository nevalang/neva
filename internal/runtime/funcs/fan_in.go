package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fanIn struct{}

func (fanIn) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	data, err := io.In.Array("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			selectedMsg, ok := data.Select(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, selectedMsg.Msg) {
				return
			}
		}
	}, nil
}
