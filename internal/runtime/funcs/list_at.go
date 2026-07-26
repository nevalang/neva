package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type listAt struct{}

// Create builds runtime closure for list item access by index.
//
// Invariants:
//  1. `idx` must be an integer message.
//  2. Positive and negative indexing are supported (`-1` means last element).
//  3. Out-of-bounds indexes are returned via `err` outport.
//  4. Typed scalar lists are handled first (int/string/bool/float) to avoid
//     untyped materialization in hot paths.
//  5. Untyped list fallback is used for non-scalar or mixed-value lists.
//
//nolint:cyclop,gocognit,gocyclo,varnamelen,funlen,nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (listAt) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	idxIn, err := io.In.Single("idx")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, idxMsg, ok := receive2(ctx, dataIn, idxIn)
			if !ok {
				return
			}

			item, found := messages.ListAt(dataMsg.List(), idxMsg.Int())
			if !found {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, item) {
				return
			}
		}
	}, nil
}
