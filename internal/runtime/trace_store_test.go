package runtime

import (
	"context"
	"strings"
	"sync"
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
	if len(hop.ParentTraceIDs) != 0 {
		t.Fatalf("expected root hop parents to be empty, got %v", hop.ParentTraceIDs)
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
	if len(path[0].ParentTraceIDs) != 1 || path[0].ParentTraceIDs[0] != path[1].TraceID {
		t.Fatalf(
			"expected parent link %d -> %d, got parents=%v",
			path[0].TraceID,
			path[1].TraceID,
			path[0].ParentTraceIDs,
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

type testFanInCreator struct{}

func (testFanInCreator) Create(io IO, _ Msg) (func(context.Context), error) {
	firstIn, err := io.In.Single("first")
	if err != nil {
		return nil, err
	}
	secondIn, err := io.In.Single("second")
	if err != nil {
		return nil, err
	}
	thirdIn, err := io.In.Single("third")
	if err != nil {
		return nil, err
	}
	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var firstMsg, secondMsg, thirdMsg Msg
		var firstOK, secondOK, thirdOK bool

		var waitGroup sync.WaitGroup
		waitGroup.Go(func() {
			firstMsg, firstOK = firstIn.Receive(ctx)
		})
		waitGroup.Go(func() {
			secondMsg, secondOK = secondIn.Receive(ctx)
		})
		waitGroup.Go(func() {
			thirdMsg, thirdOK = thirdIn.Receive(ctx)
		})
		waitGroup.Wait()

		if !firstOK || !secondOK || !thirdOK {
			return
		}

		outMsg := NewListMsg([]Msg{firstMsg, secondMsg, thirdMsg})
		_ = resOut.Send(ctx, outMsg)
	}, nil
}

func TestTracePath_FanInTracksAllParents(t *testing.T) {
	resetRuntimeTraceStateForTests()

	baseCtx := context.Background()
	handlerCtx := contextWithTraceActivation(baseCtx)
	sendCtx := context.Background()
	recvCtx := context.Background()
	firstCh := make(chan OrderedMsg, 1)
	secondCh := make(chan OrderedMsg, 1)
	thirdCh := make(chan OrderedMsg, 1)
	resCh := make(chan OrderedMsg, 1)

	firstOut := NewSingleOutport(PortAddr{Path: "first/out", Port: "res"}, ProdInterceptor{}, firstCh)
	secondOut := NewSingleOutport(PortAddr{Path: "second/out", Port: "res"}, ProdInterceptor{}, secondCh)
	thirdOut := NewSingleOutport(PortAddr{Path: "third/out", Port: "res"}, ProdInterceptor{}, thirdCh)
	resIn := NewSingleInport(resCh, PortAddr{Path: "prog/out", Port: "stop"}, ProdInterceptor{})

	handler, err := testFanInCreator{}.Create(IO{
		In: NewInports(map[string]Inport{
			"first":  NewInport(nil, NewSingleInport(firstCh, PortAddr{Path: "fanin/in", Port: "first"}, ProdInterceptor{})),
			"second": NewInport(nil, NewSingleInport(secondCh, PortAddr{Path: "fanin/in", Port: "second"}, ProdInterceptor{})),
			"third":  NewInport(nil, NewSingleInport(thirdCh, PortAddr{Path: "fanin/in", Port: "third"}, ProdInterceptor{})),
		}),
		Out: NewOutports(map[string]Outport{
			"res": NewOutport(NewSingleOutport(PortAddr{Path: "fanin/out", Port: "res"}, ProdInterceptor{}, resCh), nil),
		}),
	}, nil)
	if err != nil {
		t.Fatalf("create handler failed: %v", err)
	}

	go handler(handlerCtx)
	if !firstOut.Send(sendCtx, NewStringMsg("a")) {
		t.Fatalf("first send failed")
	}
	if !secondOut.Send(sendCtx, NewStringMsg("b")) {
		t.Fatalf("second send failed")
	}
	if !thirdOut.Send(sendCtx, NewStringMsg("c")) {
		t.Fatalf("third send failed")
	}

	last, ok := resIn.Receive(recvCtx)
	if !ok {
		t.Fatalf("receive failed")
	}

	path := TracePath(last)
	if len(path) != 2 {
		t.Fatalf("expected output hop plus first parent chain entry, got %d hops", len(path))
	}
	if len(path[0].ParentTraceIDs) != 3 {
		t.Fatalf("expected 3 parent trace ids, got %v", path[0].ParentTraceIDs)
	}
	if path[0].ParentTraceIDs[0] == 0 || path[0].ParentTraceIDs[1] == 0 || path[0].ParentTraceIDs[2] == 0 {
		t.Fatalf("expected non-zero parent trace ids, got %v", path[0].ParentTraceIDs)
	}

	formatted := FormatDataflowTrace(last)
	if !strings.Contains(formatted, "fanin:res -> prog:stop") {
		t.Fatalf("expected fan-in output hop in formatted trace, got:\n%s", formatted)
	}
}
