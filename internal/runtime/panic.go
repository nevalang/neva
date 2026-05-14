package runtime

import (
	"context"
)

// Panic control is a runtime-internal bridge for user-level panic semantics:
// components can request graceful program stop without using Go panic transport.
type programCancelCauseKey struct{}
type programExitCode int

func (c programExitCode) Error() string {
	return "program exit code"
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
	exitCode, ok := cause.(programExitCode)
	if !ok {
		return 0, false
	}
	return int(exitCode), true
}

func FailProgram(ctx context.Context) {
	mustProgramCancelCause(ctx)(programExitCode(1))
}
