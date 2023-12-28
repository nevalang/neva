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

func (d FuncRunner) Run(funcCalls []FuncCall) (func(ctx context.Context), error) {
	routines := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		constructor, ok := d.registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrFuncNotFound, call.Ref)
		}
		handler, err := constructor(call.IO, call.MetaMsg)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: ref %v", ErrFuncConstructor, err, call.Ref)
		}
		routines[i] = handler
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(len(routines))
		for i := range routines {
			routine := routines[i]
			go func() {
				routine(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}, nil
}
