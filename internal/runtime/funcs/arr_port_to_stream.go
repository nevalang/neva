package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToStream struct{}

func (arrayPortToStream) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, err := io.In.Array("port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	seqOut, err := io.Out.Single("seq")
	if err != nil {
		return nil, err
	}

	// TODO: could be optimized by using portIn.ReceiveAll()
	// but we need to handle order of sending messages to stream
	return func(ctx context.Context) {
		l := portIn.Len()

		for {
			for idx := 0; idx < l; idx++ {
				msg, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}

				item := streamItem(
					msg,
					int64(idx),
					idx == l-1,
				)

				if !seqOut.Send(ctx, item) {
					return
				}
			}
		}
	}, nil
}
