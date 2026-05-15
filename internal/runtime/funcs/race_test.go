package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// race_test.go contains unit tests for race runtime function.

// TestRaceSendsDataAndCaseCauses verifies race emits data payload with data+case causes.
func TestRaceSendsDataAndCaseCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	caseInputs := []chan runtime.OrderedMsg{make(chan runtime.OrderedMsg, 1), make(chan runtime.OrderedMsg, 1)}
	caseRead := []<-chan runtime.OrderedMsg{caseInputs[0], caseInputs[1]}
	caseOut0 := make(chan runtime.OrderedMsg, 1)
	caseOut1 := make(chan runtime.OrderedMsg, 1)
	caseOutWrite := []chan<- runtime.OrderedMsg{caseOut0, caseOut1}

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(nil, runtime.NewSingleInport(tracer, dataIn, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
			"case": runtime.NewInport(runtime.NewArrayInport(tracer, caseRead, runtime.PortAddr{Path: "test/in", Port: "case"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"case": runtime.NewOutport(nil, runtime.NewArrayOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "case"}, interceptor, caseOutWrite)),
		}),
	}

	handler, err := (race{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	srcData := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "data"}, interceptor, dataIn)
	srcCase1 := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "case1"}, interceptor, caseInputs[1])
	if !srcData.Send(ctx, runtime.NewIntMsg(42)) || !srcCase1.Send(ctx, runtime.NewStringMsg("pick-1")) {
		t.Fatal("failed to send race inputs")
	}

	select {
	case out := <-caseOut1:
		if !out.Equal(runtime.NewIntMsg(42)) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewIntMsg(42))
		}
		assertHopCauseCount(t, tracer, out, 2)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting case[1] result")
	}

	cancel()
	<-done
}
