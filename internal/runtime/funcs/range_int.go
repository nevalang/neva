package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type rangeIntV1 struct{}

func (rangeIntV1) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fromIn, err := io.In.Single("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			fromMsg, ok := fromIn.Receive(ctx)
			if !ok {
				return
			}

			toMsg, ok := toIn.Receive(ctx)
			if !ok {
				return
			}

			var (
				from = fromMsg.Int()
				to   = toMsg.Int()

				idx  int64 = 0
				last bool  = false
				data int64 = from
			)

			if from < to {
				for !last {
					if data == to-1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
					data++
				}
			} else {
				for !last {
					if data == toMsg.Int()+1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
					data--
				}
			}
		}
	}, nil
}

// Add this new struct and function after the existing rangeInt implementation

type rangeIntV2 struct{}

func (rangeIntV2) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fromIn, err := io.In.Single("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		return nil, err
	}

	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			// Wait for signal before processing
			_, ok := sigIn.Receive(ctx)
			if !ok {
				return
			}

			fromMsg, ok := fromIn.Receive(ctx)
			if !ok {
				return
			}

			toMsg, ok := toIn.Receive(ctx)
			if !ok {
				return
			}

			var (
				from = fromMsg.Int()
				to   = toMsg.Int()

				idx  int64 = 0
				last bool  = false
				data int64 = from
			)

			if from < to { // example: 0..10
				for !last {
					if data == to-1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
					data++
				}
			} else { // example: 10..0
				for !last {
					if data == toMsg.Int()+1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
					data--
				}
			}
		}
	}, nil
}
