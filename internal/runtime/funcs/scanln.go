package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type scanln struct{}

// TODO add `:err` outport
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (r scanln) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) { //nolint:gocognit,lll // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	sigIn, err := singleIn(rio, "sig")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(rio, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := singleOut(rio, "err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(input)) {
				return
			}
		}
	}, nil
}
