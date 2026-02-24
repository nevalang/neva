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
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			sepMsg, ok := sepIn.Receive(ctx)
			if !ok {
				return
			}

			builder := strings.Builder{}
			sep := sepMsg.Str()

			list := dataMsg.List()
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

			switch {
			case isStreamOpen(msg):
				builder.Reset()
			case isStreamData(msg):
				if builder.Len() > 0 {
					builder.WriteString(sep)
				}
				builder.WriteString(streamDataValue(msg).Str())
			case isStreamClose(msg):
				if !resOut.Send(ctx, runtime.NewStringMsg(builder.String())) {
					return
				}
				builder.Reset()
				hasSep = false
			}
		}
	}, nil
}
