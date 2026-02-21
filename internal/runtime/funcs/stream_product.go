package funcs

import (
	"context"
	"sync"

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

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	// TODO: make sure it's not possible to do processing on the fly so we don't have to wait for both streams to complete
	return func(ctx context.Context) {
		for {
			var (
				firstOk, secondOk bool
				firstData         = []runtime.Msg{}
				secondData        = []runtime.Msg{}
			)

			var wg sync.WaitGroup

			wg.Go(func() {
				for {
					var firstMsg runtime.Msg
					firstMsg, firstOk = firstIn.Receive(ctx)
					if !firstOk {
						return
					}

					streamItem := firstMsg.Struct()
					firstData = append(firstData, streamItem.Get("data"))

					if streamItem.Get("last").Bool() {
						break
					}
				}
			})

			wg.Go(func() {
				for {
					var secondMsg runtime.Msg
					secondMsg, secondOk = secondIn.Receive(ctx)
					if !secondOk {
						return
					}

					streamItem := secondMsg.Struct()
					secondData = append(secondData, streamItem.Get("data"))

					if streamItem.Get("last").Bool() {
						break
					}
				}
			})

			wg.Wait()

			if !firstOk || !secondOk {
				return
			}

			for i, firstMsg := range firstData {
				for j, secondMsg := range secondData {
					resOut.Send(
						ctx,
						streamItem(
							runtime.NewStructMsg([]runtime.StructField{
								runtime.NewStructField("first", firstMsg),
								runtime.NewStructField("second", secondMsg),
							}),
							int64(i*len(secondData)+j),
							i == len(firstData)-1 && j == len(secondData)-1,
						),
					)
				}
			}
		}
	}, nil
}
