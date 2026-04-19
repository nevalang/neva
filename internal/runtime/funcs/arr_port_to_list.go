package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type arrayPortToList struct{}

func (arrayPortToList) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(context.Context), error) {
	portIn, err := arrayIn(io, "port")
	if err != nil {
		return nil, errors.New("missing array inport 'port'")
	}

	listOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		l := portIn.Len()

		for {
			list := make([]runtime.Msg, 0, l)
			for idx := range l {
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
