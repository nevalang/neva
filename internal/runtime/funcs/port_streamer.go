package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type portStreamer struct{}

func (portStreamer) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portsIn, ok := io.In["ports"] // slots of the "ports" array inport
	if !ok {
		return nil, errors.New("missing port 'ports'")
	}

	streamOut, err := io.Out.Port("stream")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		// naive implementation that will performe poorly in case
		// messages arrive in different order
		// better implementation should have buffer
		// so we don't block the senders
		// but still emit messages to stream outport in order
		for _, slot := range portsIn {
			streamOut <- <-slot
		}
	}, nil
}
