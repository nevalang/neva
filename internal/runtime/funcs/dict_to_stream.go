package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type dictToStream struct{}

//nolint:gocognit // Stream framing and termination handling belong to one state machine.
func (dictToStream) Create(
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

			dict := messages.DictToMessageMap(dataMsg.Dict())
			if !resOut.Send(ctx, messages.NewStreamOpenMsg()) {
				return
			}
			for key, valueMsg := range dict {
				entryMsg := messages.NewStructMsg([]messages.StructField{
					messages.NewStructField("key", messages.NewStringMsg(key)),
					messages.NewStructField("value", valueMsg),
				})

				if !resOut.Send(ctx, messages.NewStreamDataMsg(entryMsg)) {
					return
				}
			}
			if !resOut.Send(ctx, messages.NewStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
