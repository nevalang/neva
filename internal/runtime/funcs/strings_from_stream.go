package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringsFromStream struct{}

func (stringsFromStream) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		builder := strings.Builder{}

		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			item := msg.Struct()
			builder.WriteString(item.Get("data").Str())

			if !item.Get("last").Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
				return
			}

			builder.Reset()
		}
	}, nil
}
