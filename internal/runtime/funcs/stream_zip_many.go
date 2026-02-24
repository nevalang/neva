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
			var aborted atomic.Bool

			var wg sync.WaitGroup
			wg.Add(streamsCount)

			for streamIdx := range streamsCount {
				idx := streamIdx

				go func() {
					defer wg.Done()

					collected := make([]runtime.Msg, 0)

					if !waitStreamOpen(ctx, dataInSlot{arr: dataIn, idx: idx}) {
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
						case isStreamData(msg):
							collected = append(collected, streamDataValue(msg))
						case isStreamClose(msg):
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

			if !resOut.Send(ctx, streamOpen()) {
				return
			}

			for idx := 0; idx < count; idx++ {
				zipped := make([]runtime.Msg, streamsCount)
				for streamIdx := range streamsCount {
					zipped[streamIdx] = states[streamIdx].data[idx]
				}

				if !resOut.Send(ctx, streamData(runtime.NewListMsg(zipped))) {
					return
				}
			}

			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}

type dataInSlot struct {
	arr runtime.ArrayInport
	idx int
}

func (d dataInSlot) Receive(ctx context.Context) (runtime.Msg, bool) {
	return d.arr.Receive(ctx, d.idx)
}
