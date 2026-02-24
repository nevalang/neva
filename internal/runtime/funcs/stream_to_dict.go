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
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			// Static typing guarantees stream payload is streams.Entry<T>.
			streamItemMsg := dataMsg.Struct()
			entryMsg := streamItemMsg.Get("data").Struct()
			key := entryMsg.Get("key").Str()
			valueMsg := entryMsg.Get("value")

			// Duplicate key policy: last message for the key wins.
			dict[key] = valueMsg

			if !streamItemMsg.Get("last").Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewDictMsg(dict)) {
				return
			}

			dict = map[string]runtime.Msg{}
		}
	}, nil
}
