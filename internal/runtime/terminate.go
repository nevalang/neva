package runtime

import (
	"context"
	"errors"
)

// Termination control is a runtime-internal bridge for user-level stop semantics:
// components can request graceful program termination with an explicit exit code.
type programCancelCauseKey struct{}
type programExitError struct {
	exitCode int
}

func (e programExitError) Error() string {
	return "program terminated"
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

func ProgramExitCodeFromCause(cause error) (int, bool) {
	var exitErr programExitError
	ok := errors.As(cause, &exitErr)
	if !ok {
		return 0, false
	}
	return exitErr.exitCode, true
}

// Terminate requests graceful runtime termination with the provided process exit code.
func Terminate(ctx context.Context, exitCode int) {
	mustProgramCancelCause(ctx)(programExitError{exitCode: exitCode})
}
