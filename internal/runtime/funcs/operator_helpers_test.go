package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// operator_helpers_test.go contains helper-path and operator baseline unit tests.

// TestBinaryOperatorHelperSendsTwoCauses verifies binary helper sets left+right causes.
func TestBinaryOperatorHelperSendsTwoCauses(t *testing.T) {
	t.Parallel()

	io, leftInput, rightInput, resultOutput := newBinaryIO()
	handler, err := (intAdd{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	tracer := runtime.TracerFromIO(io)
	srcLeft := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "left"}, runtime.NoEffectInterceptor{}, leftInput)
	srcRight := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "right"}, runtime.NoEffectInterceptor{}, rightInput)
	if !srcLeft.Send(ctx, runtime.NewIntMsg(20)) || !srcRight.Send(ctx, runtime.NewIntMsg(22)) {
		t.Fatal("failed to send operator inputs")
	}

	select {
	case out := <-resultOutput:
		if !out.Equal(runtime.NewIntMsg(42)) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewIntMsg(42))
		}
		assertHopCauseCount(t, tracer, out, 2)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}

// TestBinaryOperatorsBehavior checks representative binary operators per runtime func contract.
func TestBinaryOperatorsBehavior(t *testing.T) {
	t.Parallel()

	t.Run("int_add", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, intAdd{}, runtime.NewIntMsg(7), runtime.NewIntMsg(5), runtime.NewIntMsg(12))
	})
	t.Run("float_div", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, floatDiv{}, runtime.NewFloatMsg(9.0), runtime.NewFloatMsg(2.0), runtime.NewFloatMsg(4.5))
	})
	t.Run("string_add", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, stringAdd{}, runtime.NewStringMsg("ne"), runtime.NewStringMsg("va"), runtime.NewStringMsg("neva"))
	})
	t.Run("eq", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, eq{}, runtime.NewStringMsg("same"), runtime.NewStringMsg("same"), runtime.NewBoolMsg(true))
	})
	t.Run("not_eq", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, notEq{}, runtime.NewIntMsg(1), runtime.NewIntMsg(2), runtime.NewBoolMsg(true))
	})
}

// TestUnaryOperatorsBehavior checks representative unary operators per runtime func contract.
func TestUnaryOperatorsBehavior(t *testing.T) {
	t.Parallel()

	t.Run("int_inc", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intInc{}, runtime.NewIntMsg(41), runtime.NewIntMsg(42))
	})
	t.Run("int_dec", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intDec{}, runtime.NewIntMsg(41), runtime.NewIntMsg(40))
	})
	t.Run("int_neg", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intNeg{}, runtime.NewIntMsg(8), runtime.NewIntMsg(-8))
	})
}
