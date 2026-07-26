package funcs

import (
	"context"
	"fmt"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type streamProduct struct{}

//nolint:cyclop,funlen,gocognit,gocyclo // Collecting both stream lifecycles and emitting the product are one operation.
func (streamProduct) Create(
	runtimeIO runtime.IO,
	_ messages.Msg,
) (func(ctx context.Context), error) {
	firstIn, err := runtimeIO.In.Single("first")
	if err != nil {
		return nil, fmt.Errorf("get first inport: %w", err)
	}

	secondIn, err := runtimeIO.In.Single("second")
	if err != nil {
		return nil, fmt.Errorf("get second inport: %w", err)
	}

	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("get res outport: %w", err)
	}

	// TODO: make sure it's not possible to do processing on the fly so we don't have to wait for both streams to complete
	return func(ctx context.Context) {
		for {
			var (
				firstOk, secondOk bool
				firstData         = []messages.Msg{}
				secondData        = []messages.Msg{}
			)

			var group sync.WaitGroup

			group.Go(func() {
				firstOk = waitStreamOpen(ctx, firstIn)
				if !firstOk {
					return
				}
			readFirst:
				for {
					var firstMsg messages.Msg
					firstOrdered, firstOK := firstIn.Receive(ctx)
					firstMsg, firstOk = firstOrdered.Msg, firstOK
					if !firstOk {
						return
					}

					switch {
					case messages.IsStreamData(firstMsg):
						firstData = append(firstData, messages.StreamDataValue(firstMsg))
					case messages.IsStreamClose(firstMsg):
						break readFirst
					}
				}
			})

			group.Go(func() {
				secondOk = waitStreamOpen(ctx, secondIn)
				if !secondOk {
					return
				}
			readSecond:
				for {
					var secondMsg messages.Msg
					secondOrdered, secondOK := secondIn.Receive(ctx)
					secondMsg, secondOk = secondOrdered.Msg, secondOK
					if !secondOk {
						return
					}

					switch {
					case messages.IsStreamData(secondMsg):
						secondData = append(secondData, messages.StreamDataValue(secondMsg))
					case messages.IsStreamClose(secondMsg):
						break readSecond
					}
				}
			})

			group.Wait()

			if !firstOk || !secondOk {
				return
			}

			if !resOut.Send(ctx, messages.NewStreamOpenMsg()) {
				return
			}

			for _, firstMsg := range firstData {
				for _, secondMsg := range secondData {
					if !resOut.Send(
						ctx,
						messages.NewStreamDataMsg(messages.NewStructMsg([]messages.StructField{
							messages.NewStructField("first", firstMsg),
							messages.NewStructField("second", secondMsg),
						})),
					) {
						return
					}
				}
			}

			if !resOut.Send(ctx, messages.NewStreamCloseMsg()) {
				return
			}
		}
	}, nil
}
