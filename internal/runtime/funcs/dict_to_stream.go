package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type dictToStream struct{}

func (dictToStream) Create(
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

			dict := data.Dict()
			size := len(dict)

			idx := 0
			for key, value := range dict {
				entry := runtime.NewStructMsg([]runtime.StructField{
					runtime.NewStructField("key", runtime.NewStringMsg(key)),
					runtime.NewStructField("value", value),
				})

				if !resOut.Send(ctx, streamItem(entry, int64(idx), idx == size-1)) {
					return
				}

				idx++
			}
		}
	}, nil
}
