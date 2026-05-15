package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// int_add_test.go contains unit tests for intAdd runtime function.

// TestIntAddProducesExpectedValue checks arithmetic behavior.
func TestIntAddProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, intAdd{}, runtime.NewIntMsg(7), runtime.NewIntMsg(5), runtime.NewIntMsg(12))
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
	leftCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "left"}, runtime.NewIntMsg(20), leftInput)
	rightCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "right"}, runtime.NewIntMsg(22), rightInput)

	select {
	case out := <-resultOutput:
		if !out.Equal(runtime.NewIntMsg(42)) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewIntMsg(42))
		}
		assertHopCauseIndexes(t, tracer, out, []runtime.OrderedMsg{leftCause, rightCause})
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
