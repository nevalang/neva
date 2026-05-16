package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listToStream struct{}

func (c listToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			list := listToMsgs(data.List())

			for idx := range list {
				item := streamItem(
					list[idx],
					int64(idx),
					idx == len(list)-1,
				)

				if !resOut.Send(ctx, item) {
					return
				}
			}
		}
	}, nil
}

func listToMsgs(list runtime.ListMsg) []runtime.Msg {
	if values, ok := runtime.AsListInts(list); ok {
		msgs := make([]runtime.Msg, len(values))
		for i := range values {
			msgs[i] = runtime.NewIntMsg(values[i])
		}
		return msgs
	}
	if values, ok := runtime.AsListStrings(list); ok {
		msgs := make([]runtime.Msg, len(values))
		for i := range values {
			msgs[i] = runtime.NewStringMsg(values[i])
		}
		return msgs
	}
	if values, ok := runtime.AsListBools(list); ok {
		msgs := make([]runtime.Msg, len(values))
		for i := range values {
			msgs[i] = runtime.NewBoolMsg(values[i])
		}
		return msgs
	}
	if values, ok := runtime.AsListFloats(list); ok {
		msgs := make([]runtime.Msg, len(values))
		for i := range values {
			msgs[i] = runtime.NewFloatMsg(values[i])
		}
		return msgs
	}
	return list.Msgs()
}
