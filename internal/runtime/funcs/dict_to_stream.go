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
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			dict := dataMsg.Dict()
			if !resOut.Send(ctx, streamOpen()) {
				return
			}
			for key, valueMsg := range dict {
				entryMsg := runtime.NewStructMsg([]runtime.StructField{
					runtime.NewStructField("key", runtime.NewStringMsg(key)),
					runtime.NewStructField("value", valueMsg),
				})

				if !resOut.Send(ctx, streamData(entryMsg)) {
					return
				}
			}
			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}
