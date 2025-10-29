package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

func (streamZip) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	firstIn, err := io.In.Single("first")
	if err != nil {
		return nil, err
	}

	secondIn, err := io.In.Single("second")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	// TODO optimize (read 1 message at a time from each inport, then send 1 message to outport
	// close out stream as soon as one of the two input messages are last, but be careful with the
	// rest messages in the second stream, you also need to read them)
	return func(ctx context.Context) {
		for {
			firstData := []runtime.Msg{}
			for {
				seqMsg, ok := firstIn.Receive(ctx)
				if !ok {
					return
				}
				item := seqMsg.Struct()
				firstData = append(firstData, item.Get("data"))
				if item.Get("last").Bool() {
					break
				}
			}

			secondData := []runtime.Msg{}
			for {
				seqMsg, ok := secondIn.Receive(ctx)
				if !ok {
					return
				}
				item := seqMsg.Struct()
				secondData = append(secondData, item.Get("data"))
				if item.Get("last").Bool() {
					break
				}
			}

			n := len(firstData)
			if m := len(secondData); m < n {
				n = m
			}

			for i := 0; i < n; i++ {
				if !dataOut.Send(
					ctx,
					streamItem(
						runtime.NewStructMsg([]runtime.StructField{
							runtime.NewStructField("first", firstData[i]),
							runtime.NewStructField("second", secondData[i]),
						}),
						int64(i),
						i == n-1,
					),
				) {
					return
				}
			}
		}
	}, nil
}
