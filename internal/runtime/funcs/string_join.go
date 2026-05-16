package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringJoinList struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (stringJoinList) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	sepIn, err := io.In.Single("sep")
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
			dataMsg, sepMsg, ok := receive2(ctx, dataIn, sepIn)
			if !ok {
				return
			}

			builder := strings.Builder{}
			sep := sepMsg.Str()

			list := dataMsg.List()
			if stringsList, ok := runtime.AsListStrings(list); ok {
				for i := range stringsList {
					if i > 0 {
						builder.WriteString(sep)
					}
					builder.WriteString(stringsList[i])
				}
			} else {
				msgList := list.Msgs()
				for i := range msgList {
					if i > 0 {
						builder.WriteString(sep)
					}
					builder.WriteString(msgList[i].Str())
				}
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
				return
			}
		}
	}, nil
}

type stringJoinStream struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (stringJoinStream) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	sepIn, err := io.In.Single("sep")
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
		builder := strings.Builder{}
		var (
			sep    string
			hasSep bool
		)

		for {
			var msg runtime.Msg
			var ok bool

			if !hasSep {
				var sepMsg runtime.Msg
				msg, sepMsg, ok = receive2(ctx, dataIn, sepIn)
				if !ok {
					return
				}

				sep = sepMsg.Str()
			} else {
				msg, ok = dataIn.Receive(ctx)
				if !ok {
					return
				}
			}

			if !appendAndFlushJoinItem(ctx, resOut, &builder, sep, msg.Struct()) {
				return
			}
			hasSep = builder.Len() > 0
		}
	}, nil
}

func appendAndFlushJoinItem(
	ctx context.Context,
	resOut runtime.SingleOutport,
	builder *strings.Builder,
	sep string,
	item runtime.StructMsg,
) bool {
	if builder.Len() > 0 {
		builder.WriteString(sep)
	}
	builder.WriteString(item.Get("data").Str())

	if !item.Get("last").Bool() {
		return true
	}

	if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
		return false
	}
	builder.Reset()
	return true
}
