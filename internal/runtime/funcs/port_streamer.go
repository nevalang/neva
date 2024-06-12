package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToStream struct{}

func (arrayPortToStream) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, ok := io.In["port"]
	if !ok {
		return nil, errors.New("missing array inport 'port'")
	}

	streamOut, err := io.Out.SingleOutport("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			l := len(portIn)

			for i, slot := range portIn {
				var msg runtime.Msg
				select {
				case <-ctx.Done():
					return
				case msg = <-slot:
				}

				item := streamItem(
					msg,
					int64(i),
					i == l-1,
				)

				select {
				case <-ctx.Done():
					return
				case streamOut <- item:
				}
			}
		}
	}, nil
}
