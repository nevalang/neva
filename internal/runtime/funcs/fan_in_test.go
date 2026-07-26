package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
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

	cancel, done := runHandler(handler)
	ctx := context.Background()
	cause := sendTracked(
		t,
		ctx,
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "data1"},
		messages.NewStringMsg("v"),
		dataInputs[1],
	)

	select {
	case out := <-resOutCh:
		if !messages.Equal(out.Msg, messages.NewStringMsg("v")) {
			t.Fatalf("payload = %v, want %v", out, messages.NewStringMsg("v"))
		}
		assertHopCauseIndexes(t, tracer, out, []runtime.OrderedMsg{cause})
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}
