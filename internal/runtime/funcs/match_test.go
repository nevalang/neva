package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// match_test.go contains unit tests for matchSelector runtime function.

// TestMatchSendsDataIfThenCauses verifies matched output has explicit data+if+then causes.
func TestMatchSendsDataIfThenCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	ifInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	thenInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	ifRead := []<-chan runtime.OrderedMsg{ifInputs[0], ifInputs[1]}
	thenRead := []<-chan runtime.OrderedMsg{thenInputs[0], thenInputs[1]}
	elseIn := make(chan runtime.OrderedMsg, 1)
	resOut := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(nil, runtime.NewSingleInport(tracer, dataIn, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
			"if":   runtime.NewInport(runtime.NewArrayInport(tracer, ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
			"then": runtime.NewInport(runtime.NewArrayInport(tracer, thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
			"else": runtime.NewInport(nil, runtime.NewSingleInport(tracer, elseIn, runtime.PortAddr{Path: "test/in", Port: "else"}, interceptor)),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resOut), nil),
		}),
	}

	handler, err := (matchSelector{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	srcData := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "data"}, interceptor, dataIn)
	srcIf0 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "if0"}, interceptor, ifInputs[0])
	srcIf1 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "if1"}, interceptor, ifInputs[1])
	srcThen0 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "then0"}, interceptor, thenInputs[0])
	srcThen1 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "then1"}, interceptor, thenInputs[1])
	srcElse := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "else"}, interceptor, elseIn)

	if !srcData.Send(ctx, runtime.NewStringMsg("k")) ||
		!srcIf0.Send(ctx, runtime.NewStringMsg("x")) ||
		!srcIf1.Send(ctx, runtime.NewStringMsg("k")) ||
		!srcThen0.Send(ctx, runtime.NewStringMsg("zero")) ||
		!srcThen1.Send(ctx, runtime.NewStringMsg("one")) ||
		!srcElse.Send(ctx, runtime.NewStringMsg("fallback")) {
		t.Fatal("failed to send match inputs")
	}

	select {
	case out := <-resOut:
		if !out.Equal(runtime.NewStringMsg("one")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("one"))
		}
		assertHopCauseCount(t, tracer, out, 3)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
