package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type getDictValue struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (g getDictValue) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dictIn, err := singleIn(io, "dict")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	keyIn, err := singleIn(io, "key")
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
			var (
				dictMsg, keyMsg runtime.Msg
				dictOk, keyOk   bool
			)

			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			wg := sync.WaitGroup{}
			wg.Go(func() {
				dictMsg, dictOk = dictIn.Receive(ctx)
			})
			wg.Go(func() {
				keyMsg, keyOk = keyIn.Receive(ctx)
			})
			wg.Wait()
			if !dictOk || !keyOk {
				return
			}

			valueMsg, ok := dictMsg.Dict()[keyMsg.Str()]
			if !ok {
				if !errOut.Send(ctx, errFromString("Key not found in dictionary")) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, valueMsg) {
				return
			}
		}
	}, nil
}
