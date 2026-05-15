package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// fan_in_test.go contains unit tests for fanIn runtime function.

// TestFanInSendsSingleExplicitCause verifies fan_in propagates one explicit cause.
func TestFanInSendsSingleExplicitCause(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	dataRead := []<-chan runtime.OrderedMsg{dataInputs[0], dataInputs[1]}
	resOutCh := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(runtime.NewArrayInport(tracer, dataRead, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resOutCh), nil),
		}),
	}

	handler, err := (fanIn{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()

	src := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "src/out", Port: "data1"}, interceptor, dataInputs[1])
	if !src.Send(ctx, runtime.NewStringMsg("v")) {
		t.Fatal("failed to send source data")
	}

	select {
	case out := <-resOutCh:
		if !out.Equal(runtime.NewStringMsg("v")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("v"))
		}
		assertHopCauseCount(t, tracer, out, 1)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
