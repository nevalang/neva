package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type stringToStream struct{}

//nolint:gocognit // Stream framing and rune emission belong to one state machine.
func (stringToStream) Create(
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
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			// Ranging over a string splits it by Unicode code points (runes), not bytes.
			if !resOut.Send(ctx, messages.NewStreamOpenMsg()) {
				return
			}

			for _, runeValue := range dataMsg.Str() {
				if !resOut.Send(ctx, messages.NewStreamDataMsg(messages.NewStringMsg(string(runeValue)))) {
					return
				}
			}

			if !resOut.Send(ctx, messages.NewStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
