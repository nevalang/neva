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

	cancel, done := runHandler(handler)
	ctx := context.Background()

	dataCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "data"}, runtime.NewStringMsg("k"), dataIn)
	_ = sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "if0"}, runtime.NewStringMsg("x"), ifInputs[0])
	ifCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "if1"}, runtime.NewStringMsg("k"), ifInputs[1])
	_ = sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "then0"}, runtime.NewStringMsg("zero"), thenInputs[0])
	thenCause := sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "then1"}, runtime.NewStringMsg("one"), thenInputs[1])
	_ = sendTracked(t, ctx, tracer, runtime.PortAddr{Path: "src/out", Port: "else"}, runtime.NewStringMsg("fallback"), elseIn)

	select {
	case out := <-resOut:
		if !out.Equal(runtime.NewStringMsg("one")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("one"))
		}
		assertHopCauseIndexes(t, tracer, out, []runtime.OrderedMsg{dataCause, ifCause, thenCause})
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
