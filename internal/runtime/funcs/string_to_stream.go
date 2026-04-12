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
			// []rune also gives stable element count for idx/last stream metadata.
			runes := []rune(dataMsg.Str())
			for idx, runeValue := range runes {
				streamItemMsg := streamItem(
					runtime.NewStringMsg(string(runeValue)),
					int64(idx),
					idx == len(runes)-1,
				)

				if !resOut.Send(ctx, streamItemMsg) {
					return
				}
			}
		}
	}, nil
}
