package funcs

import (
	"context"
	"unicode/utf8"

	"github.com/nevalang/neva/internal/runtime"
)

type stringAt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (stringAt) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
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
			data := dataMsg.Str()
			l := int64(utf8.RuneCountInString(data))

			if idx < -l || idx >= l {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
			}

			for i, r := range data {
				if int64(i) == idx {
					if !resOut.Send(ctx, runtime.NewStringMsg(string(r))) {
						return
					}
					break
				}
			}
		}
	}, nil
}
