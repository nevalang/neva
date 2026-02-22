package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamToDict struct{}

func (streamToDict) Create(
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
		dict := map[string]runtime.Msg{}

		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			item := msg.Struct()
			entry := item.Get("data").Struct()
			key := entry.Get("key").Str()
			value := entry.Get("value")

			dict[key] = value

			if !item.Get("last").Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewDictMsg(dict)) {
				return
			}

			dict = map[string]runtime.Msg{}
		}
	}, nil
}
