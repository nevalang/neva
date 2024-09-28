package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamToList struct{}

func (s streamToList) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		list := []runtime.Msg{}

		for {
			msg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			item := msg.Struct()

			list = append(list, item.Get("data"))

			if !item.Get("last").Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewListMsg(list)) {
				return
			}

			list = []runtime.Msg{}
		}
	}, nil
}
