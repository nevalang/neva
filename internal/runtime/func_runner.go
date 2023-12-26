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

func (d FuncRunner) Run(ctx context.Context, funcRoutines []FuncCall) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	wg.Add(len(funcRoutines))

	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, routine := range funcRoutines {
		if routine.MetaMsg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, routine.MetaMsg)
		}

		constructor, ok := d.registry[routine.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrFuncNotFound, routine.Ref)
		}

		fun, err := constructor(ctx, routine.IO)
		if err != nil {
			return fmt.Errorf("%w: %v: ref %v", ErrFuncConstructor, err, routine.Ref)
		}

		go func() {
			fun() // will return at ctx.Done()
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
