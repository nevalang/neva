package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type arrayPortToList struct{}

func (arrayPortToList) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ messages.Msg,
) (func(context.Context), error) {
	portIn, err := io.In.Array("port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	listOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		l := portIn.Len()

		for {
			list := make([]messages.Msg, 0, l)
			for idx := range l {
				ordered, ok := portIn.Receive(ctx, idx)
				if !ok {
					return
				}
				list = append(list, ordered.Msg)
			}

			if !listOut.Send(ctx, messages.NewListMsg(list)) {
				return
			}
		}
	}, nil
}
