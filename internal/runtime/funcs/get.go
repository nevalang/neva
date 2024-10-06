package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type getDictValue struct{}

func (g getDictValue) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dictIn, err := io.In.Single("dict")
	if err != nil {
		return nil, err
	}

	keyIn, err := io.In.Single("key")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var (
				dictMsg, keyMsg runtime.Msg
				dictOk, keyOk   bool
			)

			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				dictMsg, dictOk = dictIn.Receive(ctx)
			}()
			go func() {
				defer wg.Done()
				keyMsg, keyOk = keyIn.Receive(ctx)
			}()
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
