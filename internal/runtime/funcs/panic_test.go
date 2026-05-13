package funcs

import (
	"context"
	"strings"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFormatPanicDataflowTrace_NormalizesPathsAndIncludesComponent(t *testing.T) {
	ctx := runtime.WithTracer(context.Background())
	ch1 := make(chan runtime.OrderedMsg, 1)
	ch2 := make(chan runtime.OrderedMsg, 1)

	out1 := runtime.NewSingleOutport(
		runtime.PortAddr{Path: "http/in", Port: "req"},
		runtime.NoEffectInterceptor{},
		ch1,
	)
	in1 := runtime.NewSingleInport(
		ch1,
		runtime.PortAddr{Path: "parse/in", Port: "data"},
		runtime.NoEffectInterceptor{},
	)
	out2 := runtime.NewSingleOutport(
		runtime.PortAddr{Path: "parse/out", Port: "req"},
		runtime.NoEffectInterceptor{},
		ch2,
	)
	in2 := runtime.NewSingleInport(
		ch2,
		runtime.PortAddr{Path: "checkout/finalize/in", Port: "err"},
		runtime.NoEffectInterceptor{},
	)

	if !out1.Send(ctx, runtime.NewStringMsg("x")) {
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

	formatted := formatPanicDataflowTrace(ctx, last)
	if !strings.Contains(formatted, "panic sink: checkout/finalize:err") {
		t.Fatalf("expected normalized panic sink, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "panic component: finalize") {
		t.Fatalf("expected panic component, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "http:req -> parse:data") {
		t.Fatalf("expected normalized hop with http:req -> parse:data, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "parse:req -> checkout/finalize:err") {
		t.Fatalf("expected normalized hop with parse:req -> checkout/finalize:err, got:\n%s", formatted)
	}
}

func TestFormatPanicDataflowTrace_FanInKeepsAllParentBranches(t *testing.T) {
	ctx := runtime.WithTracer(context.Background())
	firstCh := make(chan runtime.OrderedMsg, 1)
	secondCh := make(chan runtime.OrderedMsg, 1)
	thirdCh := make(chan runtime.OrderedMsg, 1)
	resCh := make(chan runtime.OrderedMsg, 1)

	firstOut := runtime.NewSingleOutport(runtime.PortAddr{Path: "first/out", Port: "res"}, runtime.NoEffectInterceptor{}, firstCh)
	secondOut := runtime.NewSingleOutport(runtime.PortAddr{Path: "second/out", Port: "res"}, runtime.NoEffectInterceptor{}, secondCh)
	thirdOut := runtime.NewSingleOutport(runtime.PortAddr{Path: "third/out", Port: "res"}, runtime.NoEffectInterceptor{}, thirdCh)
	firstIn := runtime.NewSingleInport(firstCh, runtime.PortAddr{Path: "fanin/in", Port: "first"}, runtime.NoEffectInterceptor{})
	secondIn := runtime.NewSingleInport(secondCh, runtime.PortAddr{Path: "fanin/in", Port: "second"}, runtime.NoEffectInterceptor{})
	thirdIn := runtime.NewSingleInport(thirdCh, runtime.PortAddr{Path: "fanin/in", Port: "third"}, runtime.NoEffectInterceptor{})
	resOut := runtime.NewSingleOutport(runtime.PortAddr{Path: "fanin/out", Port: "res"}, runtime.NoEffectInterceptor{}, resCh)
	resIn := runtime.NewSingleInport(resCh, runtime.PortAddr{Path: "prog/out", Port: "stop"}, runtime.NoEffectInterceptor{})

	if !firstOut.Send(ctx, runtime.NewStringMsg("a")) {
		t.Fatalf("first send failed")
	}
	if !secondOut.Send(ctx, runtime.NewStringMsg("b")) {
		t.Fatalf("second send failed")
	}
	if !thirdOut.Send(ctx, runtime.NewStringMsg("c")) {
		t.Fatalf("third send failed")
	}

	firstOrdered, ok := firstIn.Receive(ctx)
	if !ok {
		t.Fatalf("first receive failed")
	}
	secondOrdered, ok := secondIn.Receive(ctx)
	if !ok {
		t.Fatalf("second receive failed")
	}
	thirdOrdered, ok := thirdIn.Receive(ctx)
	if !ok {
		t.Fatalf("third receive failed")
	}

	if !resOut.Send(
		ctx,
		runtime.NewListMsg([]runtime.Msg{firstOrdered.Msg, secondOrdered.Msg, thirdOrdered.Msg}),
		firstOrdered,
		secondOrdered,
		thirdOrdered,
	) {
		t.Fatalf("result send failed")
	}

	last, ok := resIn.Receive(ctx)
	if !ok {
		t.Fatalf("result receive failed")
	}

	formatted := formatPanicDataflowTrace(ctx, last)
	if !strings.Contains(formatted, "fanin:res -> prog:stop") {
		t.Fatalf("expected fan-in output hop in formatted trace, got:\n%s", formatted)
	}
	if strings.Count(formatted, "first:res -> fanin:first")+strings.Count(formatted, "second:res -> fanin:second")+strings.Count(formatted, "third:res -> fanin:third") < 3 {
		t.Fatalf("expected all fan-in parents in formatted trace, got:\n%s", formatted)
	}
}
