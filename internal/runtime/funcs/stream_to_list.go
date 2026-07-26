package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type streamToList struct{}

func (s streamToList) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ messages.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("data")
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
		// Fully materializes one stream batch before emitting resulting list.
		list := make([]messages.Msg, 0, 1)

		for {
			msg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case messages.IsStreamOpen(msg.Msg):
				list = make([]messages.Msg, 0, 1)
				continue
			case messages.IsStreamData(msg.Msg):
				list = append(list, messages.StreamDataValue(msg.Msg))
				continue
			case !messages.IsStreamClose(msg.Msg):
				continue
			}

			if !resOut.Send(ctx, messages.ListFromMessages(list)) {
				return
			}
		}
	}, nil
}
