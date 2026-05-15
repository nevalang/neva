package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// newNamedRuntimeIO creates test runtime IO with single in/out ports.
func newNamedRuntimeIO(inNames []string, outNames []string) (runtime.IO, map[string]chan runtime.OrderedMsg, map[string]chan runtime.OrderedMsg) {
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

func newBinaryRuntimeIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
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

func newUnaryRuntimeIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
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
