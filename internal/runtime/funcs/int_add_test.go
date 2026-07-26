package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// int_add_test.go contains unit tests for intAdd runtime function.

// TestIntAddProducesExpectedValue checks arithmetic behavior.
func TestIntAddProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, intAdd{}, messages.NewIntMsg(7), messages.NewIntMsg(5), messages.NewIntMsg(12))
}

// TestIntAddSendsTwoCauses verifies binary helper path stores left/right causes.
func TestIntAddSendsTwoCauses(t *testing.T) {
	t.Parallel()

	io, leftInput, rightInput, resultOutput := newBinaryIO()
	handler, err := (intAdd{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	ctx := context.Background()
	tracer := runtime.TracerFromIO(io)
	leftCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "left"}, messages.NewIntMsg(20), leftInput)
	rightCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "right"}, messages.NewIntMsg(22), rightInput)

	select {
	case out := <-resultOutput:
		if !messages.Equal(out.Msg, messages.NewIntMsg(42)) {
			t.Fatalf("payload = %v, want %v", out, messages.NewIntMsg(42))
		}
		assertHopCauseIndexes(t, tracer, out, []runtime.OrderedMsg{leftCause, rightCause})
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
