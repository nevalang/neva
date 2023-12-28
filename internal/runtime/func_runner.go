package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrFuncNotFound    = errors.New("func not found")
	ErrFuncConstructor = errors.New("func constructor")
)

const CtxMsgKey = "msg"

type FuncRunner struct {
	registry map[string]Func
}

func MustNewFuncRunner(registry map[string]Func) FuncRunner {
	if registry == nil {
		panic(ErrNilDeps)
	}
	return FuncRunner{registry: registry}
}

func (d FuncRunner) Run(ctx context.Context, funcCalls []FuncCall) (func(), error) {
	ff := make([]func(), len(funcCalls))

	for i, call := range funcCalls {
		if call.MetaMsg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, call.MetaMsg)
		}
		constructor, ok := d.registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrFuncNotFound, call.Ref)
		}
		fun, err := constructor(ctx, call.IO)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: ref %v", ErrFuncConstructor, err, call.Ref)
		}
		ff[i] = fun
	}

	return func() {
		wg := sync.WaitGroup{}
		wg.Add(len(funcCalls))
		for _, f := range ff {
			f()
			wg.Done()
		}
		wg.Wait()
	}, nil
}
