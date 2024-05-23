package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamProduct struct{}

func (streamProduct) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	firstIn, err := io.In.Port("first")
	if err != nil {
		return nil, err
	}

	secondIn, err := io.In.Port("second")
	if err != nil {
		return nil, err
	}

	seqOut, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var firstData, secondData []runtime.Msg
			for {
				var item map[string]runtime.Msg
				select {
				case <-ctx.Done():
					return
				case seqMsg := <-firstIn:
					item = seqMsg.Map()
				}
				if firstData = append(firstData, item["data"]); item["last"].Bool() {
					break
				}
			}
			for {
				var item map[string]runtime.Msg
				select {
				case <-ctx.Done():
					return
				case seqMsg := <-secondIn:
					item = seqMsg.Map()
				}
				if secondData = append(secondData, item["data"]); item["last"].Bool() {
					break
				}
			}

			idx := int64(0)
			for i, e1 := range firstData {
				for j, e2 := range secondData {
					select {
					case <-ctx.Done():
						return
					case seqOut <- runtime.NewMapMsg(map[string]runtime.Msg{
						"idx":  runtime.NewIntMsg(idx),
						"last": runtime.NewBoolMsg(i == len(firstData)-1 && j == len(secondData)-1),
						"data": runtime.NewMapMsg(map[string]runtime.Msg{
							"first":  e1,
							"second": e2,
						}),
					}):
						idx++
					}
				}
			}
		}
	}, nil
}
