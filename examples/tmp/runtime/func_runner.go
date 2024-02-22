package runtime

import (
	"context"
	"fmt"
	"sync"
)

type FuncRunner struct {
	registry map[string]FuncCreator
}

type FuncCreator interface {
	// Create method validates the input and builds ready to use function
	Create(funcIO FuncIO, msg Msg) (func(context.Context), error)
}

func (d FuncRunner) Run(funcCalls []FuncCall) (func(ctx context.Context), error) {
	funcs := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		creator, ok := d.registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("func creator not found: %v", call.Ref)
		}

		handler, err := creator.Create(call.IO, call.ConfigMsg)
		if err != nil {
			return nil, fmt.Errorf("create: %w: %v", err, call.Ref)
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
