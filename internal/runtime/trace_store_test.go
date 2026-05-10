package runtime

import (
	"context"
	"strings"
	"testing"
)

func resetRuntimeTraceStateForTests() {
	resetTraceStoreForTests()
	counter.Store(0)
}

//nolint:gocyclo,cyclop // test intentionally validates full hop shape in one place.
func TestTracePath_Linear(t *testing.T) {
	resetRuntimeTraceStateForTests()

	ch := make(chan OrderedMsg, 1)
	out := NewSingleOutport(
		PortAddr{Path: "producer/out", Port: "res"},
		ProdInterceptor{},
		ch,
	)
	in := NewSingleInport(
		ch,
		PortAddr{Path: "consumer/in", Port: "data"},
		ProdInterceptor{},
	)

	ctx := context.Background()
	if !out.Send(ctx, NewStringMsg("hello")) {
		t.Fatalf("send failed")
	}

	got, ok := in.Receive(ctx)
	if !ok {
		t.Fatalf("receive failed")
	}

	traceID, hasTrace := TraceIDFromMsg(got)
	if !hasTrace || traceID == 0 {
		t.Fatalf("expected trace metadata on received message")
	}

	path := TracePath(got)
	if len(path) != 1 {
		t.Fatalf("expected 1 trace hop, got %d", len(path))
	}

	hop := path[0]
	if hop.ParentTraceID != 0 {
		t.Fatalf("expected root hop parent to be 0, got %d", hop.ParentTraceID)
	}
	if hop.Sender == nil || hop.Sender.Path != "producer/out" || hop.Sender.Port != "res" {
		t.Fatalf("unexpected sender hop: %#v", hop.Sender)
	}
	if hop.Receiver == nil || hop.Receiver.Path != "consumer/in" || hop.Receiver.Port != "data" {
		t.Fatalf("unexpected receiver hop: %#v", hop.Receiver)
	}
}

func TestTracePath_ForwardedMessageTracksParent(t *testing.T) {
	resetRuntimeTraceStateForTests()

	ctx := context.Background()
	ch1 := make(chan OrderedMsg, 1)
	ch2 := make(chan OrderedMsg, 1)

	out1 := NewSingleOutport(
		PortAddr{Path: "step1/out", Port: "res"},
		ProdInterceptor{},
		ch1,
	)
	in1 := NewSingleInport(
		ch1,
		PortAddr{Path: "step2/in", Port: "data"},
		ProdInterceptor{},
	)
	out2 := NewSingleOutport(
		PortAddr{Path: "step2/out", Port: "res"},
		ProdInterceptor{},
		ch2,
	)
	in2 := NewSingleInport(
		ch2,
		PortAddr{Path: "step3/in", Port: "data"},
		ProdInterceptor{},
	)

	if !out1.Send(ctx, NewStringMsg("x")) {
		t.Fatalf("first send failed")
	}
	mid, ok := in1.Receive(ctx)
	if !ok {
		t.Fatalf("first receive failed")
	}
	if !out2.Send(ctx, mid) {
		t.Fatalf("second send failed")
	}
	last, ok := in2.Receive(ctx)
	if !ok {
		t.Fatalf("second receive failed")
	}

	path := TracePath(last)
	if len(path) != 2 {
		t.Fatalf("expected 2 trace hops, got %d", len(path))
	}
	if path[0].ParentTraceID != path[1].TraceID {
		t.Fatalf(
			"expected parent link %d -> %d, got parent=%d",
			path[0].TraceID,
			path[1].TraceID,
			path[0].ParentTraceID,
		)
	}

	formatted := FormatDataflowTrace(last)
	if formatted == "" {
		t.Fatalf("expected formatted trace")
	}
}

func TestFormatDataflowTrace_NormalizesInOutPathSuffixes(t *testing.T) {
	resetRuntimeTraceStateForTests()

	ctx := context.Background()
	ch1 := make(chan OrderedMsg, 1)
	ch2 := make(chan OrderedMsg, 1)

	out1 := NewSingleOutport(
		PortAddr{Path: "http/in", Port: "req"},
		ProdInterceptor{},
		ch1,
	)
	in1 := NewSingleInport(
		ch1,
		PortAddr{Path: "parse/in", Port: "data"},
		ProdInterceptor{},
	)
	out2 := NewSingleOutport(
		PortAddr{Path: "parse/out", Port: "req"},
		ProdInterceptor{},
		ch2,
	)
	in2 := NewSingleInport(
		ch2,
		PortAddr{Path: "checkout/finalize/in", Port: "err"},
		ProdInterceptor{},
	)

	if !out1.Send(ctx, NewStringMsg("x")) {
		t.Fatalf("first send failed")
	}
	mid, ok := in1.Receive(ctx)
	if !ok {
		t.Fatalf("first receive failed")
	}
	if !out2.Send(ctx, mid) {
		t.Fatalf("second send failed")
	}
	last, ok := in2.Receive(ctx)
	if !ok {
		t.Fatalf("second receive failed")
	}

	formatted := FormatDataflowTrace(last)
	if !strings.Contains(formatted, "panic sink: checkout/finalize:err") {
		t.Fatalf("expected normalized panic sink, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "http:req -> parse:data") {
		t.Fatalf("expected normalized hop with http:req -> parse:data, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "parse:req -> checkout/finalize:err") {
		t.Fatalf("expected normalized hop with parse:req -> checkout/finalize:err, got:\n%s", formatted)
	}
}
