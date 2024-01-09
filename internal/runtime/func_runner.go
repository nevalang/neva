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

type FuncRunner struct {
	registry map[string]FuncCreator
}

type FuncCreator interface {
	// Create method validates the input and builds ready to use function
	Create(FuncIO, Msg) (func(context.Context), error)
}

func (d FuncRunner) Run(funcCalls []FuncCall) (func(ctx context.Context), error) {
	funcs := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		creator, ok := d.registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrFuncNotFound, call.Ref)
		}

		handler, err := creator.Create(call.IO, call.MetaMsg)
		if err != nil {
			return nil, fmt.Errorf("%w: %v: ref %v", ErrFuncConstructor, err, call.Ref)
		}

		funcs[i] = handler
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(len(funcs))
		for i := range funcs {
			routine := funcs[i]
			go func() {
				routine(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}, nil
}

func MustNewFuncRunner(registry map[string]FuncCreator) FuncRunner {
	if registry == nil {
		panic(ErrNilDeps)
	}
	return FuncRunner{registry: registry}
}
