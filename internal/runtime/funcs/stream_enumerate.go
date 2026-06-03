package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamEnumerate struct{}

func (streamEnumerate) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			item := msg.Struct()
			indexed := runtime.NewStructMsg([]runtime.StructField{
				runtime.NewStructField("idx", item.Get("idx")),
				runtime.NewStructField("data", item.Get("data")),
			})

			if !resOut.Send(ctx, streamItem(indexed, item.Get("idx").Int(), item.Get("last").Bool())) {
				return
			}
		}
	}, nil
}
