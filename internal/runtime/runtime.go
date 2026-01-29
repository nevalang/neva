package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

//nolint:gochecknoglobals // global monotonic counter shared by all runtime outports.
var counter atomic.Uint64

type cancelFuncKey struct{}

func contextWithCancelFunc(ctx context.Context, cancel context.CancelFunc) context.Context {
	return context.WithValue(ctx, cancelFuncKey{}, cancel)
}

// CancelFuncFromContext returns the cancel function stored by Call, if present.
func CancelFuncFromContext(ctx context.Context) (context.CancelFunc, bool) {
	v := ctx.Value(cancelFuncKey{})
	if v == nil {
		return nil, false
	}
	cancel, ok := v.(context.CancelFunc)
	return cancel, ok
}

type FuncCreator interface {
	Create(IO, Msg) (func(context.Context), error)
}

func Run(ctx context.Context, prog Program, registry map[string]FuncCreator) error {
	_, err := Call(ctx, prog, registry, NewStructMsg(nil))
	return err
}

// Call runs a single request-response round-trip using program Start/Stop.
// It sends the provided input to Start, waits for one message on Stop,
// then cancels and waits for all handlers to finish.
func Call(ctx context.Context, prog Program, registry map[string]FuncCreator, in Msg) (Msg, error) {
	var out Msg
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		out, _ = prog.Stop.Receive(ctx)
		cancel() // normal termination
	}()

	runFuncs, err := deferFuncCalls(prog.FuncCalls, registry)
	if err != nil {
		cancel()
		return nil, err
	}

	funcsFinished := make(chan struct{})

	go func() {
		// runFuncs blocks until context is cancelled (by the stop port or by panic)
		runFuncs(contextWithCancelFunc(ctx, cancel))
		close(funcsFinished)
	}()

	prog.Start.Send(ctx, in)

	<-funcsFinished

	return out, nil
}

func deferFuncCalls(
	funcCalls []FuncCall,
	registry map[string]FuncCreator,
) (func(ctx context.Context), error) {
	handlers, err := createHandlers(funcCalls, registry)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		for i := range handlers {
			routine := handlers[i]
			wg.Go(func() {
				routine(ctx)
			})
		}
		wg.Wait()
	}, nil
}

func createHandlers(
	funcCalls []FuncCall,
	registry map[string]FuncCreator,
) ([]func(context.Context), error) {
	funcs := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		creator, ok := registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("func creator not found: %v", call.Ref)
		}

		handler, err := creator.Create(call.IO, call.Config)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", call.Ref, err)
		}

		funcs[i] = handler
	}

	return funcs, nil
}
