package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoin struct{}

func (p stringJoin) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
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

			builder := strings.Builder{}
			list := data.List()
			for i := range list {
				builder.WriteString(list[i].Str())
			}

			if !resOut.Send(ctx, runtime.NewStrMsg(builder.String())) {
				return
			}
		}
	}, nil
}
