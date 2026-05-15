package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

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
	go func() {
		handler(ctx)
		close(done)
	}()

	src1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "data1"},
		interceptor,
		dataInputs[1],
	)
	if !src1.Send(ctx, runtime.NewStringMsg("v")) {
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

func TestSelectorSendsIfCause(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	ifInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	thenInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	ifRead := []<-chan runtime.OrderedMsg{ifInputs[0], ifInputs[1]}
	thenRead := []<-chan runtime.OrderedMsg{thenInputs[0], thenInputs[1]}
	resOutCh := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"if":   runtime.NewInport(runtime.NewArrayInport(tracer, ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
			"then": runtime.NewInport(runtime.NewArrayInport(tracer, thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resOutCh), nil),
		}),
	}

	handler, err := (selector{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	srcIf1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "if1"},
		interceptor,
		ifInputs[1],
	)
	srcThen0 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "then0"},
		interceptor,
		thenInputs[0],
	)
	srcThen1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "then1"},
		interceptor,
		thenInputs[1],
	)
	if !srcIf1.Send(ctx, runtime.NewBoolMsg(true)) {
		t.Fatal("failed to send if[1]")
	}
	if !srcThen0.Send(ctx, runtime.NewStringMsg("zero")) {
		t.Fatal("failed to send then[0]")
	}
	if !srcThen1.Send(ctx, runtime.NewStringMsg("one")) {
		t.Fatal("failed to send then[1]")
	}

	select {
	case out := <-resOutCh:
		if !out.Equal(runtime.NewStringMsg("one")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("one"))
		}
		assertHopCauseCount(t, tracer, out, 1)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}

func TestRaceSendsDataAndCaseCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	caseInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
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
	go func() {
		handler(ctx)
		close(done)
	}()

	srcData := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "data"},
		interceptor,
		dataIn,
	)
	srcCase1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "case1"},
		interceptor,
		caseInputs[1],
	)
	if !srcData.Send(ctx, runtime.NewIntMsg(42)) {
		t.Fatal("failed to send data")
	}
	if !srcCase1.Send(ctx, runtime.NewStringMsg("pick-1")) {
		t.Fatal("failed to send case[1]")
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

func TestSwitchMatchedCaseSendsTwoCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	caseInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	caseRead := []<-chan runtime.OrderedMsg{caseInputs[0], caseInputs[1]}
	caseOut0 := make(chan runtime.OrderedMsg, 1)
	caseOut1 := make(chan runtime.OrderedMsg, 1)
	caseOutWrite := []chan<- runtime.OrderedMsg{caseOut0, caseOut1}
	elseOut := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(nil, runtime.NewSingleInport(tracer, dataIn, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
			"case": runtime.NewInport(runtime.NewArrayInport(tracer, caseRead, runtime.PortAddr{Path: "test/in", Port: "case"}, interceptor), nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"case": runtime.NewOutport(nil, runtime.NewArrayOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "case"}, interceptor, caseOutWrite)),
			"else": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "else"}, interceptor, elseOut), nil),
		}),
	}

	handler, err := (switchRouter{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	srcData := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "data"},
		interceptor,
		dataIn,
	)
	srcCase0 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "case0"},
		interceptor,
		caseInputs[0],
	)
	srcCase1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "case1"},
		interceptor,
		caseInputs[1],
	)
	if !srcData.Send(ctx, runtime.NewStringMsg("match")) {
		t.Fatal("failed to send data")
	}
	if !srcCase0.Send(ctx, runtime.NewStringMsg("nope")) {
		t.Fatal("failed to send case[0]")
	}
	if !srcCase1.Send(ctx, runtime.NewStringMsg("match")) {
		t.Fatal("failed to send case[1]")
	}

	select {
	case out := <-caseOut1:
		if !out.Equal(runtime.NewStringMsg("match")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("match"))
		}
		assertHopCauseCount(t, tracer, out, 2)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting case[1] result")
	}

	cancel()
	<-done
}

func TestMatchSendsDataIfThenCauses(t *testing.T) {
	t.Parallel()

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataIn := make(chan runtime.OrderedMsg, 1)
	ifInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	thenInputs := []chan runtime.OrderedMsg{
		make(chan runtime.OrderedMsg, 1),
		make(chan runtime.OrderedMsg, 1),
	}
	ifRead := []<-chan runtime.OrderedMsg{ifInputs[0], ifInputs[1]}
	thenRead := []<-chan runtime.OrderedMsg{thenInputs[0], thenInputs[1]}
	elseIn := make(chan runtime.OrderedMsg, 1)
	resOut := make(chan runtime.OrderedMsg, 1)

	io := runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"data": runtime.NewInport(nil, runtime.NewSingleInport(tracer, dataIn, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
			"if":   runtime.NewInport(runtime.NewArrayInport(tracer, ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
			"then": runtime.NewInport(runtime.NewArrayInport(tracer, thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
			"else": runtime.NewInport(nil, runtime.NewSingleInport(tracer, elseIn, runtime.PortAddr{Path: "test/in", Port: "else"}, interceptor)),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resOut), nil),
		}),
	}

	handler, err := (matchSelector{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	srcData := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "data"},
		interceptor,
		dataIn,
	)
	srcIf0 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "if0"},
		interceptor,
		ifInputs[0],
	)
	srcIf1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "if1"},
		interceptor,
		ifInputs[1],
	)
	srcThen0 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "then0"},
		interceptor,
		thenInputs[0],
	)
	srcThen1 := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "then1"},
		interceptor,
		thenInputs[1],
	)
	srcElse := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "else"},
		interceptor,
		elseIn,
	)
	if !srcData.Send(ctx, runtime.NewStringMsg("k")) {
		t.Fatal("failed to send data")
	}
	if !srcIf0.Send(ctx, runtime.NewStringMsg("x")) {
		t.Fatal("failed to send if[0]")
	}
	if !srcIf1.Send(ctx, runtime.NewStringMsg("k")) {
		t.Fatal("failed to send if[1]")
	}
	if !srcThen0.Send(ctx, runtime.NewStringMsg("zero")) {
		t.Fatal("failed to send then[0]")
	}
	if !srcThen1.Send(ctx, runtime.NewStringMsg("one")) {
		t.Fatal("failed to send then[1]")
	}
	if !srcElse.Send(ctx, runtime.NewStringMsg("fallback")) {
		t.Fatal("failed to send else")
	}

	select {
	case out := <-resOut:
		if !out.Equal(runtime.NewStringMsg("one")) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewStringMsg("one"))
		}
		assertHopCauseCount(t, tracer, out, 3)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}

func TestBinaryOperatorHelperSendsTwoCauses(t *testing.T) {
	t.Parallel()

	runtimeIO, leftInput, rightInput, resultOutput := newBinaryRuntimeIO()
	handler, err := (intAdd{}).Create(runtimeIO, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	tracer := runtime.TracerFromIO(runtimeIO)
	srcLeft := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "left"},
		runtime.NoEffectInterceptor{},
		leftInput,
	)
	srcRight := runtime.NewSingleOutport(
		tracer,
		runtime.PortAddr{Path: "src/out", Port: "right"},
		runtime.NoEffectInterceptor{},
		rightInput,
	)
	if !srcLeft.Send(ctx, runtime.NewIntMsg(20)) {
		t.Fatal("failed to send left")
	}
	if !srcRight.Send(ctx, runtime.NewIntMsg(22)) {
		t.Fatal("failed to send right")
	}

	select {
	case out := <-resultOutput:
		if !out.Equal(runtime.NewIntMsg(42)) {
			t.Fatalf("payload = %v, want %v", out, runtime.NewIntMsg(42))
		}
		assertHopCauseCount(t, tracer, out, 2)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting result")
	}

	cancel()
	<-done
}

func assertHopCauseCount(t *testing.T, tracer *runtime.Tracer, msg runtime.OrderedMsg, want int) {
	t.Helper()

	hop, ok := tracer.HopByOrderedMsg(msg)
	if !ok {
		t.Fatal("hop not found for ordered message")
	}
	if len(hop.CauseIndexes) != want {
		t.Fatalf("cause count = %d, want %d (indexes=%v)", len(hop.CauseIndexes), want, hop.CauseIndexes)
	}
}
