package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoin struct{}

func (p stringJoin) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			builder := strings.Builder{}
			list := data.List()
			for i := range list {
				builder.WriteString(list[i].Str())
			}

			resOut <- runtime.NewStrMsg(builder.String())
		}
	}, nil
}
