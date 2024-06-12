package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

func (streamZip) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	firstIn, err := io.In.SingleInport("first")
	if err != nil {
		return nil, err
	}

	secondIn, err := io.In.SingleInport("second")
	if err != nil {
		return nil, err
	}

	seqOut, err := io.Out.SingleOutport("seq")
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

			n := len(firstData)
			if m := len(secondData); m < n {
				n = m
			}

			idx := int64(0)
			for i := 0; i < n; i++ {
				e1, e2 := firstData[i], secondData[i]
				select {
				case <-ctx.Done():
					return
				case seqOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"idx":  runtime.NewIntMsg(idx),
					"last": runtime.NewBoolMsg(i == n-1),
					"data": runtime.NewMapMsg(map[string]runtime.Msg{
						"first":  e1,
						"second": e2,
					}),
				}):
					idx++
				}
			}
		}
	}, nil
}
