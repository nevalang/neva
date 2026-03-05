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

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	// TODO: could be optimized by using portIn.ReceiveAll()
	// but we need to handle order of sending messages to stream
	return func(ctx context.Context) {
		l := portIn.Len()

		for {
			if !resOut.Send(ctx, streamOpen()) {
				return
			}

			for idx := range l {
				msg, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}

				if !resOut.Send(ctx, streamData(msg)) {
					return
				}
			}

			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}
