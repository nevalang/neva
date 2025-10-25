package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoinStream struct{}

func (stringJoinStream) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
		builder := strings.Builder{}
		var (
			sep    string
			hasSep bool
		)

		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !hasSep {
				sepMsg, ok := sepIn.Receive(ctx)
				if !ok {
					return
				}

				sep = sepMsg.Str()
				hasSep = true
			}

			item := msg.Struct()

			if builder.Len() > 0 {
				builder.WriteString(sep)
			}

			builder.WriteString(item.Get("data").Str())

			if !item.Get("last").Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
				return
			}

			builder.Reset()
			hasSep = false
		}
	}, nil
}
