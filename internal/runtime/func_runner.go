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

func (d FuncRunner) Run(ctx context.Context, funcCalls []FuncCall) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	wg.Add(len(funcCalls))

	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, call := range funcCalls {
		if call.MetaMsg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, call.MetaMsg)
		}

		constructor, ok := d.registry[call.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrFuncNotFound, call.Ref)
		}

		fun, err := constructor(ctx, call.IO)
		if err != nil {
			return fmt.Errorf("%w: %v: ref %v", ErrFuncConstructor, err, call.Ref)
		}

		go func() {
			fun() // will return at ctx.Done()
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
