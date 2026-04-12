package funcs

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZipMany struct{}

func (streamZipMany) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Array("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
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
			var shouldStop atomic.Bool
			var aborted atomic.Bool

			var wg sync.WaitGroup
			wg.Add(streamsCount)

			for streamIdx := range streamsCount {
				idx := streamIdx

				go func() {
					defer wg.Done()

					collected := make([]runtime.Msg, 0)

					for {
						msg, ok := dataIn.Receive(ctx, idx)
						if !ok {
							aborted.Store(true)
							return
						}

						item := msg.Struct()
						collected = append(collected, item.Get("data"))

						if item.Get("last").Bool() {
							shouldStop.Store(true)
							states[idx] = streamState{data: collected}
							return
						}
					}
				}()
			}

			wg.Wait()

			if aborted.Load() {
				return
			}

			count := len(states[0].data)
			for streamIdx := 1; streamIdx < streamsCount; streamIdx++ {
				if l := len(states[streamIdx].data); l < count {
					count = l
				}
			}

			for idx := 0; idx < count; idx++ {
				zipped := make([]runtime.Msg, streamsCount)
				for streamIdx := range streamsCount {
					zipped[streamIdx] = states[streamIdx].data[idx]
				}

				if !resOut.Send(
					ctx,
					streamItem(
						runtime.NewListMsg(zipped),
						int64(idx),
						idx == count-1,
					),
				) {
					return
				}
			}

			if shouldStop.Load() || count == 0 {
				return
			}
		}
	}, nil
}
