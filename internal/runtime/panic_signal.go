package runtime

import (
	"context"
	"fmt"
)

type programPanicKey struct{}

type programPanicSignal struct {
	payload string
}

func contextWithProgramPanicChan(ctx context.Context, signals chan<- Msg) context.Context {
	return context.WithValue(ctx, programPanicKey{}, signals)
}

func programPanicChanFromContext(ctx context.Context) (chan<- Msg, bool) {
	v := ctx.Value(programPanicKey{})
	if v == nil {
		return nil, false
	}
	signals, ok := v.(chan<- Msg)
	return signals, ok
}

// IsProgramPanic reports whether recovered panic value is user-level program panic.
func IsProgramPanic(recovered any) bool {
	_, ok := recovered.(programPanicSignal)
	return ok
}

// ReportProgramPanic sends user-level panic signal for current runtime call.
func ReportProgramPanic(ctx context.Context, msg Msg) {
	signals, ok := programPanicChanFromContext(ctx)
	if !ok {
		return
	}
	select {
	case signals <- msg:
	default:
	}
}

func panicWithProgramSignal(msg Msg) {
	panic(programPanicSignal{
		payload: fmt.Sprint(msg),
	})
}
