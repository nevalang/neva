package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamForEachController serializes stream item handling for ForEach.
// It forwards Open/Close immediately and forwards Data only after done signal.
type streamForEachController struct{}

//nolint:cyclop,gocognit,gocyclo // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (streamForEachController) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	doneIn, err := input.In.Single("done")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	itemOut, err := input.Out.Single("item")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case isStreamOpen(msg):
				if !resOut.Send(ctx, msg) {
					return
				}
			case isStreamData(msg):
				item := streamDataValue(msg)
				if !itemOut.Send(ctx, item) {
					return
				}

				if _, ok := doneIn.Receive(ctx); !ok {
					return
				}

				if !resOut.Send(ctx, msg) {
					return
				}
			case isStreamClose(msg):
				if !resOut.Send(ctx, msg) {
					return
				}
			default:
				panic("stream_for_each_controller: unexpected stream tag")
			}
		}
	}, nil
}
