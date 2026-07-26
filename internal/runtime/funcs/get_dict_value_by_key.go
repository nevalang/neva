package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// getDictValueByKey implements the internal runtime function behind the public Get component.
type getDictValueByKey struct{}

type getDictValueByKeyPorts struct {
	dictIn runtime.SingleInport
	keyIn  runtime.SingleInport
	resOut runtime.SingleOutport
	errOut runtime.SingleOutport
}

func (getDictValueByKey) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	ports, err := resolveGetDictValueByKeyPorts(runtimeIO)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runGetDictValueByKey(ctx, &ports)
	}, nil
}

func resolveGetDictValueByKeyPorts(runtimeIO runtime.IO) (getDictValueByKeyPorts, error) {
	dictIn, err := runtimeIO.In.Single("dict")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return getDictValueByKeyPorts{}, err
	}

	keyIn, err := runtimeIO.In.Single("key")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return getDictValueByKeyPorts{}, err
	}

	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return getDictValueByKeyPorts{}, err
	}

	errOut, err := runtimeIO.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return getDictValueByKeyPorts{}, err
	}

	return getDictValueByKeyPorts{
		dictIn: dictIn,
		keyIn:  keyIn,
		resOut: resOut,
		errOut: errOut,
	}, nil
}

func runGetDictValueByKey(ctx context.Context, ports *getDictValueByKeyPorts) {
	for {
		dictMsg, keyMsg, received := receive2(ctx, ports.dictIn, ports.keyIn)
		if !received {
			return
		}

		valueMsg, found := messages.DictGetValueByKey(dictMsg.Dict(), keyMsg.Str())
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
