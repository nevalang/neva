package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type portSequencer struct{}

func (portSequencer) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, ok := io.In["port"]
	if !ok {
		return nil, errors.New("missing array inport 'port'")
	}

	streamOut, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var msg runtime.Msg

		for {
			for _, slot := range portIn {
				select {
				case <-ctx.Done():
					return
				case msg = <-slot:
				}

				select {
				case <-ctx.Done():
					return
				case streamOut <- msg:
				}
			}

			select {
			case <-ctx.Done():
				return
			case streamOut <- nil:
			}
		}
	}, nil
}
