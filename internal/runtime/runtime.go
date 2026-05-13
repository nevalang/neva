package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

//nolint:gochecknoglobals // global monotonic counter shared by all runtime outports.
var counter atomic.Uint64

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
//
//nolint:ireturn,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func Call(ctx context.Context, prog Program, registry map[string]FuncCreator, in Msg) (Msg, error) {
	var out Msg
	tracer := NewTracer()
	ctx, cancel := context.WithCancelCause(ctx)
	ctx = contextWithTracer(ctx, tracer)
	ctx = contextWithProgramCancelCause(ctx, cancel)
	ctx = contextWithTraceActivation(ctx)
	go func() {
		out, _ = prog.Stop.Receive(ctx)
		cancel(nil) // normal termination
	}()

	runFuncs, err := deferFuncCalls(prog.FuncCalls, registry)
	if err != nil {
		cancel(nil)
		return nil, err
	}

	funcsFinished := make(chan struct{})

	go func() {
		// runFuncs blocks until context is cancelled (by the stop port or by panic)
		runFuncs(ctx)
		close(funcsFinished)
	}()

	prog.Start.Send(ctx, in)

	<-funcsFinished

	if cause := context.Cause(ctx); IsProgramPanicError(cause) {
		return nil, fmt.Errorf("program panic: %w", cause)
	}

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
		//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		wg := sync.WaitGroup{}
		for i := range handlers {
			routine := handlers[i]
			wg.Go(func() {
				routine(contextWithTraceActivation(ctx))
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

	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
