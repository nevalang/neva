package funcs

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZipMany struct{}

//nolint:cyclop,funlen,gocognit,gocyclo // Coordinated collection and emission across all streams form one operation.
func (streamZipMany) Create(
	runtimeIO runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := runtimeIO.In.Array("data")
	if err != nil {
		return nil, fmt.Errorf("get data inport: %w", err)
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		streamsCount := dataIn.Len()
		if streamsCount == 0 {
			return
		}

		for {
			type streamState struct {
				data []runtime.Msg
			}

			states := make([]streamState, streamsCount)
			var aborted atomic.Bool

			var group sync.WaitGroup
			group.Add(streamsCount)

			for streamIdx := range streamsCount {
				idx := streamIdx

				go func() {
					defer group.Done()

					collected := make([]runtime.Msg, 0)

					if !waitStreamOpen(ctx, &dataInSlot{arr: dataIn, idx: idx}) {
						aborted.Store(true)
						return
					}

					for {
						msg, ok := dataIn.Receive(ctx, idx)
						if !ok {
							aborted.Store(true)
							return
						}

						switch {
						case runtime.IsStreamData(msg):
							collected = append(collected, runtime.StreamDataValue(msg))
						case runtime.IsStreamClose(msg):
							states[idx] = streamState{data: collected}
							return
						}
					}
				}()
			}

			group.Wait()

			if aborted.Load() {
				return
			}

			count := len(states[0].data)
			for streamIdx := 1; streamIdx < streamsCount; streamIdx++ {
				if l := len(states[streamIdx].data); l < count {
					count = l
				}
			}

			if !resOut.Send(ctx, runtime.NewStreamOpenMsg()) {
				return
			}

			for idx := range count {
				zipped := make([]runtime.Msg, streamsCount)
				for streamIdx := range streamsCount {
					zipped[streamIdx] = states[streamIdx].data[idx]
				}

				if !resOut.Send(ctx, runtime.NewStreamDataMsg(runtime.NewListMsg(zipped))) {
					return
				}
			}

			if !resOut.Send(ctx, runtime.NewStreamCloseMsg()) {
				return
			}
		}
	}, nil
}

type dataInSlot struct {
	arr runtime.ArrayInport
	idx int
}

func (slot *dataInSlot) Receive(ctx context.Context) (runtime.OrderedMsg, bool) {
	return slot.arr.Receive(ctx, slot.idx)
}
