package funcs

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToStream struct{}

//nolint:gocognit // Stream framing and per-port forwarding belong to one state machine.
func (arrayPortToStream) Create(
	runtimeIO runtime.IO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, err := runtimeIO.In.Array("port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("get res outport: %w", err)
	}

	// TODO: could be optimized by using portIn.ReceiveAll()
	// but we need to handle order of sending messages to stream
	return func(ctx context.Context) {
		portLen := portIn.Len()

		for {
			if !resOut.Send(ctx, runtime.NewStreamOpenMsg()) {
				return
			}

			for idx := range portLen {
				msg, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}

				if !resOut.Send(ctx, runtime.NewStreamDataMsg(msg.Msg)) {
					return
				}
			}

			if !resOut.Send(ctx, runtime.NewStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
