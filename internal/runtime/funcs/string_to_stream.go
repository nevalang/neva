package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type stringToStream struct{}

func (stringToStream) Create(
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

			// We split by Unicode code points (runes), not bytes.
			// Byte iteration would break multibyte UTF-8 chars into fragments.
			runes := []rune(dataMsg.Str())
			if !resOut.Send(ctx, streamOpen()) {
				return
			}

			for _, runeValue := range runes {
				if !resOut.Send(ctx, streamData(runtime.NewStringMsg(string(runeValue)))) {
					return
				}
			}

			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}
