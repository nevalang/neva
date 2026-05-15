package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// switch_test.go contains unit tests for switchRouter runtime function.

// TestSwitchMatchedCaseSendsTwoCauses verifies matched case output has data+case causes.
func TestSwitchMatchedCaseSendsTwoCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	caseInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	caseRead := []<-chan runtime.OrderedMsg{caseInputs[0], caseInputs[1]}
	caseOut0 := make(chan runtime.OrderedMsg, 1)
	caseOut1 := make(chan runtime.OrderedMsg, 1)
	caseOutWrite := []chan<- runtime.OrderedMsg{caseOut0, caseOut1}
	elseOut := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(nil, runtime.NewSingleInport(tracer, dataIn, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
			"case": runtime.NewInport(runtime.NewArrayInport(tracer, caseRead, runtime.PortAddr{Path: "test/in", Port: "case"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"case": runtime.NewOutport(nil, runtime.NewArrayOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "case"}, interceptor, caseOutWrite)),
			"else": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "else"}, interceptor, elseOut), nil),
		}),
	}

	handler, err := (switchRouter{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	srcData := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "data"}, interceptor, dataIn)
	srcCase0 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "case0"}, interceptor, caseInputs[0])
	srcCase1 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "case1"}, interceptor, caseInputs[1])
	if !srcData.Send(ctx, runtime.NewStringMsg("match")) || !srcCase0.Send(ctx, runtime.NewStringMsg("nope")) || !srcCase1.Send(ctx, runtime.NewStringMsg("match")) {
		t.Fatal("failed to send switch inputs")
	}

	select {
	case out := <-caseOut1:
		if !out.Equal(runtime.NewStringMsg("match")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("match"))
		}
		assertHopCauseCount(t, tracer, out, 2)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting case[1] result")
	}

	cancel()
	<-done
}
