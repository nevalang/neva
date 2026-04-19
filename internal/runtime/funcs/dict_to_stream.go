package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type dictToStream struct{}

func (dictToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			dict := dataMsg.Dict()
			// Go map iteration order is intentionally non-deterministic.
			size := len(dict)

			idx := 0
			for key, valueMsg := range dict {
				entryMsg := runtime.NewStructMsg([]runtime.StructField{
					runtime.NewStructField("key", runtime.NewStringMsg(key)),
					runtime.NewStructField("value", valueMsg),
				})

				if !resOut.Send(
					ctx,
					streamItem(entryMsg, int64(idx), idx == size-1),
				) {
					return
				}

				idx++
			}
		}
	}, nil
}
