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
	Create(FuncIO, Msg) (func(context.Context), error)
}

// Run returns a function that runs all runtime functions with their configurations.
// Each runtime function runs in a goroutine, but returned function blocks until all functions finish.
// You can cancel the context that you pass to returned function, to finish all runtime functions.
// It returns error if any of the runtime functions fails to start because of invalid configuration.
func (d FuncRunner) Run(funcCalls []FuncCall) (func(ctx context.Context), error) {
	handlers, err := d.createHandlers(funcCalls)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(len(handlers))
		for i := range handlers {
			routine := handlers[i]
			go func() {
				routine(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}, nil
}

func (d FuncRunner) createHandlers(funcCalls []FuncCall) ([]func(context.Context), error) {
	funcs := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		creator, ok := d.registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("func creator not found: %v", call.Ref)
		}

		handler, err := creator.Create(call.IO, call.ConfigMsg)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", call.Ref, err)
		}

		funcs[i] = handler
	}

	return funcs, nil
}

func NewFuncRunner(registry map[string]FuncCreator) FuncRunner {
	return FuncRunner{registry: registry}
}
