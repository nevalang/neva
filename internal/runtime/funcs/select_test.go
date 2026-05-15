package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// select_test.go contains unit tests for selector runtime function.

// TestSelectorSendsIfCause verifies selected branch output contains explicit if-cause.
func TestSelectorSendsIfCause(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	ifInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	thenInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	ifRead := []<-chan runtime.OrderedMsg{ifInputs[0], ifInputs[1]}
	thenRead := []<-chan runtime.OrderedMsg{thenInputs[0], thenInputs[1]}
	resOutCh := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"if":   runtime.NewInport(runtime.NewArrayInport(tracer, ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
			"then": runtime.NewInport(runtime.NewArrayInport(tracer, thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resOutCh), nil),
		}),
	}

	handler, err := (selector{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	srcIf := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "if1"}, interceptor, ifInputs[1])
	srcThen0 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "then0"}, interceptor, thenInputs[0])
	srcThen1 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "then1"}, interceptor, thenInputs[1])
	if !srcIf.Send(ctx, runtime.NewBoolMsg(true)) || !srcThen0.Send(ctx, runtime.NewStringMsg("zero")) || !srcThen1.Send(ctx, runtime.NewStringMsg("one")) {
		t.Fatal("failed to send selector inputs")
	}

	select {
	case out := <-resOutCh:
		if !out.Equal(runtime.NewStringMsg("one")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("one"))
		}
		assertHopCauseCount(t, tracer, out, 1)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
