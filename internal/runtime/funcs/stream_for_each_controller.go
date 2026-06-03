package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamForEachController serializes stream item handling for ForEach.
// It forwards each item only after the handler signals completion.
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
			msg, dataOK := dataIn.Receive(ctx)
			if !dataOK {
				return
			}

			item := msg.Struct()
			if !itemOut.Send(ctx, item.Get("data"), msg) {
				return
			}

			if _, doneOK := doneIn.Receive(ctx); !doneOK {
				return
			}

			if !resOut.Send(ctx, streamItem(item.Get("data"), item.Get("idx").Int(), item.Get("last").Bool()), msg) {
				return
			}
		}
	}, nil
}
