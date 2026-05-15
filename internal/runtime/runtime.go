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

func Run(ctx context.Context, prog Program, registry map[string]FuncCreator) (int, error) {
	_, exitCode, err := callOrdered(ctx, prog, registry, NewStructMsg(nil))
	return exitCode, err
}

// Call runs a single request-response round-trip using program Start/Stop.
// It sends the provided input to Start, waits for one message on Stop,
// then cancels and waits for all handlers to finish.
//
//nolint:ireturn,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func Call(ctx context.Context, prog Program, registry map[string]FuncCreator, input Msg) (Msg, int, error) {
	orderedOut, exitCode, err := callOrdered(ctx, prog, registry, input)
	if err != nil || exitCode != 0 {
		return nil, exitCode, err
	}
	return orderedOut.Msg, 0, nil
}

// callOrdered executes one runtime round-trip and returns the stop-port envelope.
func callOrdered(
	ctx context.Context,
	prog Program,
	registry map[string]FuncCreator,
	input Msg,
) (OrderedMsg, int, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	ctx = contextWithCancelFunc(ctx, cancel)

	var out OrderedMsg
	go func() {
		out, _ = prog.Stop.Receive(ctx)
		cancel(nil) // normal termination
	}()

	runFuncs, err := deferFuncCalls(prog.FuncCalls, registry)
	if err != nil {
		cancel(nil)
		return OrderedMsg{}, 0, err
	}

	funcsFinished := make(chan struct{})
	go func() {
		runFuncs(ctx) // blocks until context is cancelled
		close(funcsFinished)
	}()

	prog.Start.Send(ctx, input)

	<-funcsFinished

	if exitCode, ok := programExitCodeFromCause(context.Cause(ctx)); ok {
		return OrderedMsg{}, exitCode, nil
	}

	return out, 0, nil
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
