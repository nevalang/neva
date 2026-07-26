package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// dictGetByKey implements the internal runtime function behind the public Get component.
type dictGetByKey struct{}

//nolint:gocognit,gocyclo,cyclop,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (dictGetByKey) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
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
				dictMsg, keyMsg messages.Msg
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

			valueMsg, found := dictValueByKey(dictMsg.Dict(), keyMsg.Str())
			if !found {
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

// dictValueByKey reads one value without converting a typed dict into map[string]Msg.
//
//nolint:ireturn // Runtime function output is expressed as Msg.
func dictValueByKey(dict messages.DictMsg, key string) (messages.Msg, bool) {
	if values, ok := messages.AsDictInts(dict); ok {
		value, found := values[key]
		return messages.NewIntMsg(value), found
	}
	if values, ok := messages.AsDictStrings(dict); ok {
		value, found := values[key]
		return messages.NewStringMsg(value), found
	}
	if values, ok := messages.AsDictBools(dict); ok {
		value, found := values[key]
		return messages.NewBoolMsg(value), found
	}
	if values, ok := messages.AsDictFloats(dict); ok {
		value, found := values[key]
		return messages.NewFloatMsg(value), found
	}
	value, found := dict.Untyped()[key]
	return value, found
}
