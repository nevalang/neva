package funcs

import (
	"context"
	"strings"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestTraceFromOrderedMsg_ReconstructsFanInTree(t *testing.T) {
	tracer := runtime.NewTracer()

	ctx := context.Background()
	firstCh := make(chan runtime.OrderedMsg, 1)
	secondCh := make(chan runtime.OrderedMsg, 1)
	resCh := make(chan runtime.OrderedMsg, 1)

	firstOut := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "first/out", Port: "res"}, runtime.NoEffectInterceptor{}, firstCh)
	secondOut := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "second/out", Port: "res"}, runtime.NoEffectInterceptor{}, secondCh)
	firstIn := runtime.NewSingleInport(tracer, firstCh, runtime.PortAddr{Path: "join/in", Port: "first"}, runtime.NoEffectInterceptor{})
	secondIn := runtime.NewSingleInport(tracer, secondCh, runtime.PortAddr{Path: "join/in", Port: "second"}, runtime.NoEffectInterceptor{})
	resOut := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "join/out", Port: "res"}, runtime.NoEffectInterceptor{}, resCh)
	resIn := runtime.NewSingleInport(tracer, resCh, runtime.PortAddr{Path: "prog/out", Port: "stop"}, runtime.NoEffectInterceptor{})

	if !firstOut.Send(ctx, runtime.NewStringMsg("a")) {
		t.Fatalf("first send failed")
	}
	if !secondOut.Send(ctx, runtime.NewStringMsg("b")) {
		t.Fatalf("second send failed")
	}

	firstMsg, ok := firstIn.Receive(ctx)
	if !ok {
		t.Fatalf("first receive failed")
	}
	secondMsg, ok := secondIn.Receive(ctx)
	if !ok {
		t.Fatalf("second receive failed")
	}

	if !resOut.Send(ctx, runtime.NewStringMsg("ab"), firstMsg, secondMsg) {
		t.Fatalf("result send failed")
	}
	out, ok := resIn.Receive(ctx)
	if !ok {
		t.Fatalf("result receive failed")
	}

	tree, ok := traceFromOrderedMsg(tracer, out)
	if !ok {
		t.Fatalf("expected trace tree")
	}
	if len(tree.Parents) != 2 {
		t.Fatalf("expected 2 parents, got %d", len(tree.Parents))
	}
}

func TestFormatTerminationDataflowTrace_MissingTracePanics(t *testing.T) {
	tracer := runtime.NewTracer()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on missing termination trace")
		}
	}()

	_ = formatTerminationDataflowTrace("panic cause dataflow trace", tracer, runtime.OrderedMsg{})
}

func TestFormatTerminationDataflowTrace_ContainsSinkAndComponent(t *testing.T) {
	tracer := runtime.NewTracer()

	ctx := context.Background()
	ch := make(chan runtime.OrderedMsg, 1)
	out := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "producer/out", Port: "res"}, runtime.NoEffectInterceptor{}, ch)
	in := runtime.NewSingleInport(tracer, ch, runtime.PortAddr{Path: "prog/out", Port: "stop"}, runtime.NoEffectInterceptor{})

	if !out.Send(ctx, runtime.NewStringMsg("hello")) {
		t.Fatalf("send failed")
	}
	msg, ok := in.Receive(ctx)
	if !ok {
		t.Fatalf("receive failed")
	}

	got := formatTerminationDataflowTrace("panic cause dataflow trace", tracer, msg)
	if !strings.Contains(got, "sink: prog:stop") {
		t.Fatalf("expected sink in formatted trace, got:\n%s", got)
	}
	if !strings.Contains(got, "component: prog") {
		t.Fatalf("expected component in formatted trace, got:\n%s", got)
	}
}
