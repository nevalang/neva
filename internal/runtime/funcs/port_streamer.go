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
	portIn, err := io.In.Array("port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	seqOut, err := io.Out.SingleOutport("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			l := portIn.Len()
			if !portIn.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				item := streamItem(
					msg,
					int64(idx),
					idx == l-1,
				)
				return seqOut.Send(ctx, item)
			}) {
				return
			}
		}
	}, nil
}
