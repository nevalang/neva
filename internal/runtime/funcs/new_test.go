package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// New must preserve the trigger signal as a cause of emitted constant.
func TestNewV2_TracksSignalCause(t *testing.T) {
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}

	sigCh := make(chan runtime.OrderedMsg, 1)
	resCh := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"sig": runtime.NewInport(
				nil,
				runtime.NewSingleInport(tracer, sigCh, runtime.PortAddr{Path: "__newv2__1/in", Port: "sig"}, interceptor),
			),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(
				runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "__newv2__1/out", Port: "res"}, interceptor, resCh),
				nil,
			),
		}),
	}

	handler, err := newV2{}.Create(io, runtime.NewIntMsg(10))
	if err != nil {
		t.Fatalf("create newV2 handler: %v", err)
	}

	ctx := t.Context()
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	startOut := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "main/in", Port: "start"},
		interceptor,
		sigCh,
	)
	if !startOut.Send(ctx, runtime.NewStructMsg(nil)) {
		t.Fatalf("send start signal")
	}

	var out runtime.OrderedMsg
	select {
	case out = <-resCh:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting newV2 output")
	}

	outHop, ok := tracer.HopByOrderedMsg(out)
	if !ok {
		t.Fatal("missing output hop")
	}
	if len(outHop.CauseIndexes) != 1 {
		t.Fatalf("expected exactly one cause, got %v", outHop.CauseIndexes)
	}

	parents := tracer.HopsByCauseIndexes(outHop.CauseIndexes)
	if len(parents) != 1 {
		t.Fatalf("expected exactly one parent hop, got %d", len(parents))
	}
	if parents[0].Sender == nil {
		t.Fatal("parent sender is nil")
	}
	if parents[0].Sender.Port != "start" {
		t.Fatalf("expected parent sender port start, got %q", parents[0].Sender.Port)
	}
}
