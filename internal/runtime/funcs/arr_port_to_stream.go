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

	return func(ctx context.Context) {
		portLen := portIn.Len()
		items := make([]runtime.OrderedMsg, portLen)

		for {
			if !resOut.Send(ctx, newStreamOpenMsg()) {
				return
			}

			if !portIn.ReceiveAll(ctx, func(idx int, msg runtime.OrderedMsg) bool {
				items[idx] = msg
				return true
			}) {
				return
			}

			for idx := range portLen {
				if !resOut.Send(ctx, newStreamDataMsg(items[idx].Msg)) {
					return
				}
			}

			if !resOut.Send(ctx, newStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
