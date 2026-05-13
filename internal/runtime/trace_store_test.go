package runtime

import (
	"context"
	"sync"
	"testing"
)

func resetRuntimeTraceStateForTests() { counter.Store(0) }

//nolint:gocyclo,cyclop // test intentionally validates full hop shape in one place.
func TestTraceTree_Linear(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ch := make(chan OrderedMsg, 1)
	out := NewSingleOutport(
		PortAddr{Path: "producer/out", Port: "res"},
		NoEffectInterceptor{},
		ch,
	)
	in := NewSingleInport(
		ch,
		PortAddr{Path: "consumer/in", Port: "data"},
		NoEffectInterceptor{},
	)

	ctx := contextWithTracer(context.Background(), tracer)
	if !out.Send(ctx, NewStringMsg("hello")) {
		t.Fatalf("send failed")
	}

	got, ok := in.Receive(ctx)
	if !ok {
		t.Fatalf("receive failed")
	}

	if got.index == 0 {
		t.Fatalf("expected ordered message index")
	}

	tree, ok := TraceCauseTree(ctx, got)
	if !ok {
		t.Fatalf("expected trace tree")
	}
	hop := tree.Hop
	if len(hop.CauseIndexes) != 0 {
		t.Fatalf("expected root hop parents to be empty, got %v", hop.CauseIndexes)
	}
	if len(tree.Parents) != 0 {
		t.Fatalf("expected no parent nodes, got %d", len(tree.Parents))
	}
	if hop.Sender == nil || hop.Sender.Path != "producer/out" || hop.Sender.Port != "res" {
		t.Fatalf("unexpected sender hop: %#v", hop.Sender)
	}
	if hop.Receiver == nil || hop.Receiver.Path != "consumer/in" || hop.Receiver.Port != "data" {
		t.Fatalf("unexpected receiver hop: %#v", hop.Receiver)
	}
}

func TestTraceTree_ForwardedMessageTracksParent(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ctx := contextWithTracer(context.Background(), tracer)
	ch1 := make(chan OrderedMsg, 1)
	ch2 := make(chan OrderedMsg, 1)

	out1 := NewSingleOutport(
		PortAddr{Path: "step1/out", Port: "res"},
		NoEffectInterceptor{},
		ch1,
	)
	in1 := NewSingleInport(
		ch1,
		PortAddr{Path: "step2/in", Port: "data"},
		NoEffectInterceptor{},
	)
	out2 := NewSingleOutport(
		PortAddr{Path: "step2/out", Port: "res"},
		NoEffectInterceptor{},
		ch2,
	)
	in2 := NewSingleInport(
		ch2,
		PortAddr{Path: "step3/in", Port: "data"},
		NoEffectInterceptor{},
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

	tree, ok := TraceCauseTree(ctx, last)
	if !ok {
		t.Fatalf("expected trace tree")
	}
	if len(tree.Hop.CauseIndexes) != 1 || len(tree.Parents) != 1 {
		t.Fatalf(
			"expected one parent node for hop %d, got ids=%v parents=%d",
			tree.Hop.Index,
			tree.Hop.CauseIndexes,
			len(tree.Parents),
		)
	}
	if tree.Hop.CauseIndexes[0] != tree.Parents[0].Hop.Index {
		t.Fatalf("expected parent link to parent node index, got ids=%v parent=%d", tree.Hop.CauseIndexes, tree.Parents[0].Hop.Index)
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

//nolint:gocyclo,cyclop // test intentionally exercises multi-parent fan-in setup end to end.
func TestTraceTree_FanInTracksAllParents(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	baseCtx := contextWithTracer(context.Background(), tracer)
	handlerCtx := baseCtx
	sendCtx := baseCtx
	recvCtx := baseCtx
	firstCh := make(chan OrderedMsg, 1)
	secondCh := make(chan OrderedMsg, 1)
	thirdCh := make(chan OrderedMsg, 1)
	resCh := make(chan OrderedMsg, 1)

	firstOut := NewSingleOutport(PortAddr{Path: "first/out", Port: "res"}, NoEffectInterceptor{}, firstCh)
	secondOut := NewSingleOutport(PortAddr{Path: "second/out", Port: "res"}, NoEffectInterceptor{}, secondCh)
	thirdOut := NewSingleOutport(PortAddr{Path: "third/out", Port: "res"}, NoEffectInterceptor{}, thirdCh)
	resIn := NewSingleInport(resCh, PortAddr{Path: "prog/out", Port: "stop"}, NoEffectInterceptor{})

	handler, err := testFanInCreator{}.Create(IO{
		In: NewInports(map[string]Inport{
			"first":  NewInport(nil, NewSingleInport(firstCh, PortAddr{Path: "fanin/in", Port: "first"}, NoEffectInterceptor{})),
			"second": NewInport(nil, NewSingleInport(secondCh, PortAddr{Path: "fanin/in", Port: "second"}, NoEffectInterceptor{})),
			"third":  NewInport(nil, NewSingleInport(thirdCh, PortAddr{Path: "fanin/in", Port: "third"}, NoEffectInterceptor{})),
		}),
		Out: NewOutports(map[string]Outport{
			"res": NewOutport(NewSingleOutport(PortAddr{Path: "fanin/out", Port: "res"}, NoEffectInterceptor{}, resCh), nil),
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

	tree, hasTree := TraceCauseTree(baseCtx, last)
	if !hasTree {
		t.Fatalf("expected trace tree")
	}
	if len(tree.Hop.CauseIndexes) != 3 {
		t.Fatalf("expected 3 parent indexes, got %v", tree.Hop.CauseIndexes)
	}
	if len(tree.Parents) != 3 {
		t.Fatalf("expected 3 parent nodes, got %d", len(tree.Parents))
	}
	if tree.Hop.CauseIndexes[0] == 0 || tree.Hop.CauseIndexes[1] == 0 || tree.Hop.CauseIndexes[2] == 0 {
		t.Fatalf("expected non-zero parent indexes, got %v", tree.Hop.CauseIndexes)
	}
	for _, parent := range tree.Parents {
		if parent.Hop.Index == 0 {
			t.Fatalf("expected non-zero parent hop index")
		}
	}
}

func TestTraceTree_ExplicitSendCausesTrackSynthesizedOutput(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ctx := contextWithTracer(context.Background(), tracer)
	firstCh := make(chan OrderedMsg, 1)
	secondCh := make(chan OrderedMsg, 1)
	resCh := make(chan OrderedMsg, 1)

	firstOut := NewSingleOutport(PortAddr{Path: "first/out", Port: "res"}, NoEffectInterceptor{}, firstCh)
	secondOut := NewSingleOutport(PortAddr{Path: "second/out", Port: "res"}, NoEffectInterceptor{}, secondCh)
	firstIn := NewSingleInport(firstCh, PortAddr{Path: "join/in", Port: "first"}, NoEffectInterceptor{})
	secondIn := NewSingleInport(secondCh, PortAddr{Path: "join/in", Port: "second"}, NoEffectInterceptor{})
	resOut := NewSingleOutport(PortAddr{Path: "join/out", Port: "res"}, NoEffectInterceptor{}, resCh)
	resIn := NewSingleInport(resCh, PortAddr{Path: "prog/out", Port: "stop"}, NoEffectInterceptor{})

	if !firstOut.Send(ctx, NewStringMsg("a")) {
		t.Fatalf("first send failed")
	}
	if !secondOut.Send(ctx, NewStringMsg("b")) {
		t.Fatalf("second send failed")
	}

	firstOrdered, ok := firstIn.Receive(ctx)
	if !ok {
		t.Fatalf("first receive failed")
	}
	secondOrdered, ok := secondIn.Receive(ctx)
	if !ok {
		t.Fatalf("second receive failed")
	}

	if !resOut.Send(
		ctx,
		NewStringMsg(firstOrdered.Str()+secondOrdered.Str()),
		firstOrdered,
		secondOrdered,
	) {
		t.Fatalf("result send failed")
	}

	last, ok := resIn.Receive(ctx)
	if !ok {
		t.Fatalf("result receive failed")
	}

	tree, hasTree := TraceCauseTree(ctx, last)
	if !hasTree {
		t.Fatalf("expected trace tree")
	}
	if len(tree.Hop.CauseIndexes) != 2 {
		t.Fatalf("expected 2 explicit parents, got %v", tree.Hop.CauseIndexes)
	}
	if len(tree.Parents) != 2 {
		t.Fatalf("expected 2 parent nodes, got %d", len(tree.Parents))
	}
}
