package funcs

import (
	"context"
	"strconv"

	"github.com/nevalang/neva/internal/runtime"
)

type formatInt struct{}

func (formatInt) Create(
	//nolint:varnamelen
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	baseIn, err := io.In.Single("base")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			base, ok := baseIn.Receive(ctx)
			if !ok {
				return
			}

			res := strconv.FormatInt(data.Int(), int(base.Int()))
			if !resOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}
