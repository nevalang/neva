package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// dictGetByKey implements the internal runtime function behind the public Get component.
type dictGetByKey struct{}

type dictGetPorts struct {
	dictIn runtime.SingleInport
	keyIn  runtime.SingleInport
	resOut runtime.SingleOutport
	errOut runtime.SingleOutport
}

func (dictGetByKey) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	ports, err := resolveDictGetPorts(runtimeIO)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runDictGet(ctx, &ports)
	}, nil
}

func resolveDictGetPorts(runtimeIO runtime.IO) (dictGetPorts, error) {
	dictIn, err := runtimeIO.In.Single("dict")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return dictGetPorts{}, err
	}

	keyIn, err := runtimeIO.In.Single("key")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return dictGetPorts{}, err
	}

	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return dictGetPorts{}, err
	}

	errOut, err := runtimeIO.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return dictGetPorts{}, err
	}

	return dictGetPorts{
		dictIn: dictIn,
		keyIn:  keyIn,
		resOut: resOut,
		errOut: errOut,
	}, nil
}

func runDictGet(ctx context.Context, ports *dictGetPorts) {
	for {
		dictMsg, keyMsg, received := receive2(ctx, ports.dictIn, ports.keyIn)
		if !received {
			return
		}

		valueMsg, found := messages.DictGet(dictMsg.Dict(), keyMsg.Str())
		if !found {
			if !ports.errOut.Send(ctx, errFromString("Key not found in dictionary")) {
				return
			}
			continue
		}
		if !ports.resOut.Send(ctx, valueMsg) {
			return
		}
	}
}
