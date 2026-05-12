package runtime

import (
	"context"
	"errors"
)

// Panic control is a runtime-internal bridge for user-level panic semantics:
// components can request graceful program stop without using Go panic transport.
type programCancelCauseKey struct{}

type ProgramPanicError struct{}

func (ProgramPanicError) Error() string {
	return "program panicked"
}

func contextWithProgramCancelCause(
	ctx context.Context,
	cancel context.CancelCauseFunc,
) context.Context {
	return context.WithValue(ctx, programCancelCauseKey{}, cancel)
}

func mustProgramCancelCause(ctx context.Context) context.CancelCauseFunc {
	cancel, ok := ctx.Value(programCancelCauseKey{}).(context.CancelCauseFunc)
	if !ok || cancel == nil {
		panic("runtime invariant: program cancel cause func is missing")
	}
	return cancel
}

func IsProgramPanicError(err error) bool {
	var target ProgramPanicError
	return errors.As(err, &target)
}

func FailProgram(ctx context.Context) {
	mustProgramCancelCause(ctx)(ProgramPanicError{})
}
