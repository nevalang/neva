package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type stringAt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (stringAt) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
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

			result, found := messages.StringAt(dataMsg.Msg, idxMsg.Int())
			if !found {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, result) {
				return
			}
		}
	}, nil
}
