package runtime

import (
	"context"
	"sync"
	"testing"
)

func resetRuntimeTraceStateForTests() { counter.Store(0) }

//nolint:gocyclo,cyclop // test intentionally validates full hop shape in one place.
func TestTraceStore_Linear(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ch := make(chan OrderedMsg, 1)
	out := NewSingleOutport(tracer,
		PortAddr{Path: "producer/out", Port: "res"},
		NoEffectInterceptor{},
		ch,
	)
	in := NewSingleInport(tracer,
		ch,
		PortAddr{Path: "consumer/in", Port: "data"},
		NoEffectInterceptor{},
	)

	ctx := context.Background()
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

	hop, ok := tracer.HopByOrderedMsg(got)
	if !ok {
		t.Fatalf("expected trace hop")
	}
	if len(hop.CauseIndexes) != 0 {
		t.Fatalf("expected root hop parents to be empty, got %v", hop.CauseIndexes)
	}
	if hop.Sender == nil || hop.Sender.Path != "producer/out" || hop.Sender.Port != "res" {
		t.Fatalf("unexpected sender hop: %#v", hop.Sender)
	}
	if hop.Receiver == nil || hop.Receiver.Path != "consumer/in" || hop.Receiver.Port != "data" {
		t.Fatalf("unexpected receiver hop: %#v", hop.Receiver)
	}
}

func TestTraceStore_ForwardedMessageTracksParent(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ctx := context.Background()
	ch1 := make(chan OrderedMsg, 1)
	ch2 := make(chan OrderedMsg, 1)

	out1 := NewSingleOutport(tracer,
		PortAddr{Path: "step1/out", Port: "res"},
		NoEffectInterceptor{},
		ch1,
	)
	in1 := NewSingleInport(tracer,
		ch1,
		PortAddr{Path: "step2/in", Port: "data"},
		NoEffectInterceptor{},
	)
	out2 := NewSingleOutport(tracer,
		PortAddr{Path: "step2/out", Port: "res"},
		NoEffectInterceptor{},
		ch2,
	)
	in2 := NewSingleInport(tracer,
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

	hop, ok := tracer.HopByOrderedMsg(last)
	if !ok {
		t.Fatalf("expected trace hop")
	}
	if len(hop.CauseIndexes) != 1 {
		t.Fatalf(
			"expected one parent index for hop %d, got ids=%v",
			hop.Index,
			hop.CauseIndexes,
		)
	}
	parentHops := tracer.HopsByCauseIndexes(hop.CauseIndexes)
	if len(parentHops) != 1 {
		t.Fatalf("expected one resolved parent hop, got %d", len(parentHops))
	}
	if hop.CauseIndexes[0] != parentHops[0].Index {
		t.Fatalf("expected parent link to parent hop index, got ids=%v parent=%d", hop.CauseIndexes, parentHops[0].Index)
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
		var firstMsg, secondMsg, thirdMsg OrderedMsg
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
		_ = resOut.Send(ctx, outMsg, firstMsg, secondMsg, thirdMsg)
	}, nil
}

//nolint:gocyclo,cyclop // test intentionally exercises multi-parent fan-in setup end to end.
func TestTraceStore_FanInTracksAllParents(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	baseCtx := context.Background()
	handlerCtx := baseCtx
	sendCtx := baseCtx
	recvCtx := baseCtx
	firstCh := make(chan OrderedMsg, 1)
	secondCh := make(chan OrderedMsg, 1)
	thirdCh := make(chan OrderedMsg, 1)
	resCh := make(chan OrderedMsg, 1)

	firstOut := NewSingleOutport(tracer, PortAddr{Path: "first/out", Port: "res"}, NoEffectInterceptor{}, firstCh)
	secondOut := NewSingleOutport(tracer, PortAddr{Path: "second/out", Port: "res"}, NoEffectInterceptor{}, secondCh)
	thirdOut := NewSingleOutport(tracer, PortAddr{Path: "third/out", Port: "res"}, NoEffectInterceptor{}, thirdCh)
	resIn := NewSingleInport(tracer, resCh, PortAddr{Path: "prog/out", Port: "stop"}, NoEffectInterceptor{})

	firstIn := NewSingleInport(tracer, firstCh, PortAddr{Path: "fanin/in", Port: "first"}, NoEffectInterceptor{})
	secondIn := NewSingleInport(tracer, secondCh, PortAddr{Path: "fanin/in", Port: "second"}, NoEffectInterceptor{})
	thirdIn := NewSingleInport(tracer, thirdCh, PortAddr{Path: "fanin/in", Port: "third"}, NoEffectInterceptor{})
	resOut := NewSingleOutport(tracer, PortAddr{Path: "fanin/out", Port: "res"}, NoEffectInterceptor{}, resCh)

	handler, err := testFanInCreator{}.Create(IO{
		In: NewInports(map[string]Inport{
			"first":  NewInport(nil, firstIn),
			"second": NewInport(nil, secondIn),
			"third":  NewInport(nil, thirdIn),
		}),
		Out: NewOutports(map[string]Outport{
			"res": NewOutport(resOut, nil),
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

	hop, ok := tracer.HopByOrderedMsg(last)
	if !ok {
		t.Fatalf("expected trace hop")
	}
	if len(hop.CauseIndexes) != 3 {
		t.Fatalf("expected 3 parent indexes, got %v", hop.CauseIndexes)
	}
	parentHops := tracer.HopsByCauseIndexes(hop.CauseIndexes)
	if len(parentHops) != 3 {
		t.Fatalf("expected 3 resolved parent hops, got %d", len(parentHops))
	}
	if hop.CauseIndexes[0] == 0 || hop.CauseIndexes[1] == 0 || hop.CauseIndexes[2] == 0 {
		t.Fatalf("expected non-zero parent indexes, got %v", hop.CauseIndexes)
	}
	for _, parentHop := range parentHops {
		if parentHop.Index == 0 {
			t.Fatalf("expected non-zero parent hop index")
		}
	}
}

func TestTraceStore_ExplicitSendCausesTrackSynthesizedOutput(t *testing.T) {
	resetRuntimeTraceStateForTests()
	tracer := NewTracer()

	ctx := context.Background()
	firstCh := make(chan OrderedMsg, 1)
	secondCh := make(chan OrderedMsg, 1)
	resCh := make(chan OrderedMsg, 1)

	firstOut := NewSingleOutport(tracer, PortAddr{Path: "first/out", Port: "res"}, NoEffectInterceptor{}, firstCh)
	secondOut := NewSingleOutport(tracer, PortAddr{Path: "second/out", Port: "res"}, NoEffectInterceptor{}, secondCh)
	firstIn := NewSingleInport(tracer, firstCh, PortAddr{Path: "join/in", Port: "first"}, NoEffectInterceptor{})
	secondIn := NewSingleInport(tracer, secondCh, PortAddr{Path: "join/in", Port: "second"}, NoEffectInterceptor{})
	resOut := NewSingleOutport(tracer, PortAddr{Path: "join/out", Port: "res"}, NoEffectInterceptor{}, resCh)
	resIn := NewSingleInport(tracer, resCh, PortAddr{Path: "prog/out", Port: "stop"}, NoEffectInterceptor{})

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

	hop, ok := tracer.HopByOrderedMsg(last)
	if !ok {
		t.Fatalf("expected trace hop")
	}
	if len(hop.CauseIndexes) != 2 {
		t.Fatalf("expected 2 explicit parents, got %v", hop.CauseIndexes)
	}
	parentHops := tracer.HopsByCauseIndexes(hop.CauseIndexes)
	if len(parentHops) != 2 {
		t.Fatalf("expected 2 resolved parent hops, got %d", len(parentHops))
	}
}
