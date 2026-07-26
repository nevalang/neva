package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type streamToDict struct{}

func (streamToDict) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ messages.Msg,
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
		dict := map[string]messages.Msg{}

		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case isStreamOpen(dataMsg.Msg):
				dict = map[string]messages.Msg{}
				continue
			case isStreamData(dataMsg.Msg):
				entryMsg := streamDataValue(dataMsg.Msg).Struct()
				key := messages.StructGet(entryMsg, "key").Str()
				valueMsg := messages.StructGet(entryMsg, "value")

				// Duplicate key policy: last message for the key wins.
				dict[key] = valueMsg
				continue
			case !isStreamClose(dataMsg.Msg):
				continue
			}

			if !resOut.Send(ctx, messages.NewDictMsg(dict)) {
				return
			}
		}
	}, nil
}
