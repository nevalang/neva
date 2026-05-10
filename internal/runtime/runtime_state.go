package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type runtimeStateKey struct{}

type runtimeState struct {
	panicErr error
	mu       sync.RWMutex
}

func contextWithRuntimeState(ctx context.Context, state *runtimeState) context.Context {
	return context.WithValue(ctx, runtimeStateKey{}, state)
}

func runtimeStateFromContext(ctx context.Context) (*runtimeState, bool) {
	v := ctx.Value(runtimeStateKey{})
	if v == nil {
		return nil, false
	}
	state, ok := v.(*runtimeState)
	return state, ok
}

type programPanicError struct {
	payload string
}

func (e programPanicError) Error() string {
	return "panic: " + e.payload
}

// IsProgramPanic reports whether the runtime returned a program-level panic error.
func IsProgramPanic(err error) bool {
	var programPanicError programPanicError
	ok := errors.As(err, &programPanicError)
	return ok
}

// ReportProgramPanic stores panic signal in current runtime call context.
func ReportProgramPanic(ctx context.Context, msg Msg) {
	state, ok := runtimeStateFromContext(ctx)
	if !ok {
		return
	}
	state.mu.Lock()
	defer state.mu.Unlock()

	if state.panicErr != nil {
		return
	}
	state.panicErr = programPanicError{
		payload: fmt.Sprint(UnwrapTraceMsg(msg)),
	}
}

func panicErrorFromContext(ctx context.Context) error {
	state, ok := runtimeStateFromContext(ctx)
	if !ok {
		return nil
	}

	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.panicErr
}
