package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToList struct{}

func (arrayPortToList) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, err := io.In.Array("port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	listOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		l := portIn.Len()

		for {
			list := make([]runtime.Msg, 0, l)
			for idx := 0; idx < l; idx++ {
				msg, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}
				list = append(list, msg)
			}

			if !listOut.Send(ctx, runtime.NewListMsg(list)) {
				return
			}
		}
	}, nil
}
