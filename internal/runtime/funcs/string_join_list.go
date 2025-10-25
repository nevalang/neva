package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoinList struct{}

func (stringJoinList) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	sepIn, err := io.In.Single("sep")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			sepMsg, ok := sepIn.Receive(ctx)
			if !ok {
				return
			}

			sep := sepMsg.Str()

			builder := strings.Builder{}
			list := data.List()
			for i := range list {
				if i > 0 {
					builder.WriteString(sep)
				}
				builder.WriteString(list[i].Str())
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
				return
			}
		}
	}, nil
}
