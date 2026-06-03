package funcs

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// helpers_test.go contains shared test helpers for runtime func unit tests.

// newIO creates test IO with named single inports and outports.
func newIO(inNames []string, outNames []string) (runtime.IO, map[string]chan runtime.OrderedMsg, map[string]chan runtime.OrderedMsg) {
	interceptor := runtime.NoEffectInterceptor{}
	tracer := runtime.NewTracer()
	inports := make(map[string]runtime.Inport, len(inNames))
	outports := make(map[string]runtime.Outport, len(outNames))
	inChans := make(map[string]chan runtime.OrderedMsg, len(inNames))
	outChans := make(map[string]chan runtime.OrderedMsg, len(outNames))

	for _, name := range inNames {
		ch := make(chan runtime.OrderedMsg)
		inChans[name] = ch
		port := runtime.NewSingleInport(tracer, ch, runtime.PortAddr{Path: "test/in", Port: name}, interceptor)
		inports[name] = runtime.NewInport(nil, port)
	}

	for _, name := range outNames {
		ch := make(chan runtime.OrderedMsg, 1)
		outChans[name] = ch
		port := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: name}, interceptor, ch)
		outports[name] = runtime.NewOutport(port, nil)
	}

	return runtime.IO{In: runtime.NewInports(inports), Out: runtime.NewOutports(outports)}, inChans, outChans
}

// newBinaryIO creates test IO for binary operators.
func newBinaryIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	leftIn := make(chan runtime.OrderedMsg)
	rightIn := make(chan runtime.OrderedMsg)
	resultOut := make(chan runtime.OrderedMsg, 1)
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	leftPort := runtime.NewSingleInport(tracer, leftIn, runtime.PortAddr{Path: "test/in", Port: "left"}, interceptor)
	rightPort := runtime.NewSingleInport(tracer, rightIn, runtime.PortAddr{Path: "test/in", Port: "right"}, interceptor)
	resPort := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut)
	inports := runtime.NewInports(map[string]runtime.Inport{
		"left":  runtime.NewInport(nil, leftPort),
		"right": runtime.NewInport(nil, rightPort),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(resPort, nil),
	})

	return runtime.IO{In: inports, Out: outports}, leftIn, rightIn, resultOut
}

// newUnaryIO creates test IO for unary operators.
func newUnaryIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	input := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	dataPort := runtime.NewSingleInport(tracer, input, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)
	resPort := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut)
	inports := runtime.NewInports(map[string]runtime.Inport{
		"data": runtime.NewInport(nil, dataPort),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(resPort, nil),
	})

	return runtime.IO{In: inports, Out: outports}, input, resultOut
}

// sendInOrder pushes named payloads in a specified order and fails on blocking send.
func sendInOrder(
	t *testing.T,
	inChans map[string]chan runtime.OrderedMsg,
	order []string,
	payload map[string]runtime.Msg,
) {
	t.Helper()

	sendDone := make(chan struct{})
	go func() {
		for _, name := range order {
			inChans[name] <- runtime.OrderedMsg{Msg: payload[name]}
		}
		close(sendDone)
	}()

	select {
	case <-sendDone:
	case <-time.After(time.Second):
		t.Fatalf("sending blocked for order %v", order)
	}
}

// assertOutputEquals asserts expected message on selected output channel.
func assertOutputEquals(
	t *testing.T,
	outChans map[string]chan runtime.OrderedMsg,
	outName string,
	want runtime.Msg,
	order []string,
) {
	t.Helper()

	select {
	case got := <-outChans[outName]:
		if !got.Equal(want) {
			t.Fatalf("result = %v, want %v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatalf("no result for order %v", order)
	}
}

// assertHopCauseIndexes checks exact parent hop indexes for output message.
func assertHopCauseIndexes(t *testing.T, tracer *runtime.Tracer, msg runtime.OrderedMsg, expected []runtime.OrderedMsg) {
	t.Helper()

	hop, ok := tracer.HopByOrderedMsg(msg)
	if !ok {
		t.Fatal("hop not found for ordered message")
	}

	want := make([]uint64, 0, len(expected))
	for _, expectedCause := range expected {
		causeHop, ok := tracer.HopByOrderedMsg(expectedCause)
		if !ok {
			t.Fatalf("cause hop not found for ordered message: %v", expectedCause)
		}
		want = append(want, causeHop.Index)
	}

	got := append([]uint64(nil), hop.CauseIndexes...)
	slices.Sort(got)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Fatalf("cause indexes = %v, want %v", got, want)
	}
}

// runHandler starts runtime function loop and returns cancel+done hooks.
func runHandler(handler func(context.Context)) (context.CancelFunc, <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	return cancel, done
}

// sendTracked sends a message through a source outport and returns its ordered wrapper.
func sendTracked(
	t *testing.T,
	ctx context.Context,
	tracer *runtime.Tracer,
	addr runtime.PortAddr,
	msg runtime.Msg,
	dst chan runtime.OrderedMsg,
) runtime.OrderedMsg {
	t.Helper()

	srcCh := make(chan runtime.OrderedMsg, 1)
	srcOut := runtime.NewSingleOutport(tracer, addr, runtime.NoEffectInterceptor{}, srcCh)
	forwarded := make(chan runtime.OrderedMsg, 1)

	go func() {
		ordered := <-srcCh
		forwarded <- ordered
		dst <- ordered
	}()

	if !srcOut.Send(ctx, msg) {
		t.Fatalf("send failed from %s:%s", addr.Path, addr.Port)
	}

	select {
	case ordered := <-forwarded:
		return ordered
	case <-time.After(time.Second):
		t.Fatalf("timed out tracking send from %s:%s", addr.Path, addr.Port)
		return runtime.OrderedMsg{}
	}
}

// assertBinaryOperatorResult checks binary runtime func correctness for both send orders.
func assertBinaryOperatorResult(
	t *testing.T,
	creator runtime.FuncCreator,
	left runtime.Msg,
	right runtime.Msg,
	expected runtime.Msg,
) {
	t.Helper()

	io, leftInput, rightInput, resultOutput := newBinaryIO()
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	for _, sendRightFirst := range []bool{false, true} {
		sendDone := make(chan struct{})
		go func(sendRightFirst bool) {
			if sendRightFirst {
				rightInput <- runtime.OrderedMsg{Msg: right}
				leftInput <- runtime.OrderedMsg{Msg: left}
			} else {
				leftInput <- runtime.OrderedMsg{Msg: left}
				rightInput <- runtime.OrderedMsg{Msg: right}
			}
			close(sendDone)
		}(sendRightFirst)

		select {
		case <-sendDone:
		case <-time.After(time.Second):
			t.Fatalf("sending inputs blocked (sendRightFirst=%v)", sendRightFirst)
		}

		select {
		case result := <-resultOutput:
			if !result.Equal(expected) {
				t.Fatalf("result = %v, want %v (sendRightFirst=%v)", result, expected, sendRightFirst)
			}
		case <-time.After(time.Second):
			t.Fatalf("operator did not produce output in time (sendRightFirst=%v)", sendRightFirst)
		}
	}

	cancel()
	<-done
}

// assertUnaryOperatorResult checks unary runtime func correctness.
func assertUnaryOperatorResult(
	t *testing.T,
	creator runtime.FuncCreator,
	input runtime.Msg,
	expected runtime.Msg,
) {
	t.Helper()

	io, dataInput, resultOutput := newUnaryIO()
	handler, err := creator.Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	dataInput <- runtime.OrderedMsg{Msg: input}

	result := <-resultOutput
	if !result.Equal(expected) {
		t.Fatalf("result = %v, want %v", result, expected)
	}

	cancel()
	<-done
}
