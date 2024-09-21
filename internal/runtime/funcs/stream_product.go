package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamProduct struct{}

func (streamProduct) Create(
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

	seqOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			firstData := []runtime.Msg{}
			for {
				seqMsg, ok := firstIn.Receive(ctx)
				if !ok {
					return
				}

				item := seqMsg.Map()
				firstData = append(firstData, item["data"])

				if item["last"].Bool() {
					break
				}
			}

			secondData := []runtime.Msg{}
			for {
				seqMsg, ok := secondIn.Receive(ctx)
				if !ok {
					return
				}

				item := seqMsg.Map()
				secondData = append(secondData, item["data"])

				if item["last"].Bool() {
					break
				}
			}

			for i, msg1 := range firstData {
				for j, msg2 := range secondData {
					seqOut.Send(
						ctx,
						streamItem(
							runtime.NewMapMsg(map[string]runtime.Msg{
								"first":  msg1,
								"second": msg2,
							}),
							int64(i),
							i == len(firstData)-1 && j == len(secondData)-1,
						),
					)
				}
			}
		}
	}, nil
}
