package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type getDictValue struct{}

//nolint:gocognit,gocyclo,cyclop,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (g getDictValue) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dictIn, err := io.In.Single("dict")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	keyIn, err := io.In.Single("key")
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

			key := keyMsg.Str()
			dict := dictMsg.Dict()
			switch {
			case func() bool {
				values, ok := runtime.AsDictInts(dict)
				if !ok {
					return false
				}
				value, found := values[key]
				if !found {
					return false
				}
				return resOut.Send(ctx, runtime.NewIntMsg(value))
			}():
			case func() bool {
				values, ok := runtime.AsDictStrings(dict)
				if !ok {
					return false
				}
				value, found := values[key]
				if !found {
					return false
				}
				return resOut.Send(ctx, runtime.NewStringMsg(value))
			}():
			case func() bool {
				values, ok := runtime.AsDictBools(dict)
				if !ok {
					return false
				}
				value, found := values[key]
				if !found {
					return false
				}
				return resOut.Send(ctx, runtime.NewBoolMsg(value))
			}():
			case func() bool {
				values, ok := runtime.AsDictFloats(dict)
				if !ok {
					return false
				}
				value, found := values[key]
				if !found {
					return false
				}
				return resOut.Send(ctx, runtime.NewFloatMsg(value))
			}():
			default:
				valueMsg, ok := dict.Msgs()[key]
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
		}
	}, nil
}
