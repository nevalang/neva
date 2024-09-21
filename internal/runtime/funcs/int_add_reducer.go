package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intAddReducer struct{}

func (intAddReducer) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	firstIn, err := io.In.Single("first")
	if err != nil {
		return nil, err
	}

	secondIn, err := io.In.Single("second")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			firstMsg, ok := firstIn.Receive(ctx)
			if !ok {
				return
			}

			secondMsg, ok := secondIn.Receive(ctx)
			if !ok {
				return
			}

			resMsg := runtime.NewIntMsg(firstMsg.Int() + secondMsg.Int())
			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
