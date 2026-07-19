package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type stringToStream struct{}

//nolint:gocognit // Stream framing and rune emission belong to one state machine.
func (stringToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
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
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			// We split by Unicode code points (runes), not bytes.
			// Byte iteration would break multibyte UTF-8 chars into fragments.
			runes := []rune(dataMsg.Str())
			if !resOut.Send(ctx, newStreamOpenMsg()) {
				return
			}

			for _, runeValue := range runes {
				if !resOut.Send(ctx, newStreamDataMsg(runtime.NewStringMsg(string(runeValue)))) {
					return
				}
			}

			if !resOut.Send(ctx, newStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
