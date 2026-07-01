package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamMapController preserves stream item ordering for Map.
// For each item it waits for the mapped payload before forwarding the result.
type streamMapController struct{}

//nolint:cyclop,gocognit,gocyclo // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (streamMapController) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	mappedIn, err := input.In.Single("mapped")
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

			mappedMsg, ok := mappedIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, streamItem(mappedMsg.Msg, item.Get("idx").Int(), item.Get("last").Bool()), msg, mappedMsg) {
				return
			}
		}
	}, nil
}
