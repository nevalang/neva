package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listAt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (listAt) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	idxIn, err := singleIn(io, "idx")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := singleOut(io, "err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			idxMsg, ok := idxIn.Receive(ctx)
			if !ok {
				return
			}

			idx := idxMsg.Int()
			data := dataMsg.List()

			l := int64(len(data))
			if idx < -l || idx >= l {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
			}

			if idx < 0 {
				// support negative indexing:
				//	$l = [1, 2, 3]
				//	$l[-1] // 3
				idx += int64(len(data))
			}

			if !resOut.Send(ctx, data[idx]) {
				return
			}
		}
	}, nil
}
