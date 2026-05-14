package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

func startHandler(ctx context.Context, handler func(context.Context)) (context.CancelFunc, <-chan struct{}) {
	runCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	go func() {
		handler(runCtx)
		close(done)
	}()
	return cancel, done
}

func benchNewBinaryRuntimeIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	leftIn := make(chan runtime.OrderedMsg, 1)
	rightIn := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)

	inports := runtime.NewInports(map[string]runtime.Inport{
		"left":  runtime.NewInport(nil, runtime.NewSingleInport(tracer, leftIn, runtime.PortAddr{Path: "test/in", Port: "left"}, interceptor)),
		"right": runtime.NewInport(nil, runtime.NewSingleInport(tracer, rightIn, runtime.PortAddr{Path: "test/in", Port: "right"}, interceptor)),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, leftIn, rightIn, resultOut
}

func benchNewSelectRuntimeIO(size int) (runtime.IO, []chan runtime.OrderedMsg, []chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	ifInputs := make([]chan runtime.OrderedMsg, size)
	ifRead := make([]<-chan runtime.OrderedMsg, size)
	thenInputs := make([]chan runtime.OrderedMsg, size)
	thenRead := make([]<-chan runtime.OrderedMsg, size)

	for i := range size {
		ifInputs[i] = make(chan runtime.OrderedMsg, 1)
		ifRead[i] = ifInputs[i]
		thenInputs[i] = make(chan runtime.OrderedMsg, 1)
		thenRead[i] = thenInputs[i]
	}

	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	inports := runtime.NewInports(map[string]runtime.Inport{
		"if":   runtime.NewInport(runtime.NewArrayInport(tracer, ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
		"then": runtime.NewInport(runtime.NewArrayInport(tracer, thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
	})

	resultOut := make(chan runtime.OrderedMsg, 1)
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, ifInputs, thenInputs, resultOut
}

func benchNewStructRuntimeIO(names []string) (runtime.IO, map[string]chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	inputs := make(map[string]chan runtime.OrderedMsg, len(names))
	inportsMap := make(map[string]runtime.Inport, len(names))

	for _, name := range names {
		ch := make(chan runtime.OrderedMsg, 1)
		inputs[name] = ch
		inportsMap[name] = runtime.NewInport(nil, runtime.NewSingleInport(tracer, ch, runtime.PortAddr{Path: "test/in", Port: name}, interceptor))
	}

	resultOut := make(chan runtime.OrderedMsg, 1)
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: runtime.NewInports(inportsMap), Out: outports}, inputs, resultOut
}

func benchNewThreeInputsRuntimeIO(firstName, secondName, thirdName string) (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	firstIn := make(chan runtime.OrderedMsg, 1)
	secondIn := make(chan runtime.OrderedMsg, 1)
	thirdIn := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)

	inports := runtime.NewInports(map[string]runtime.Inport{
		firstName:  runtime.NewInport(nil, runtime.NewSingleInport(tracer, firstIn, runtime.PortAddr{Path: "test/in", Port: firstName}, interceptor)),
		secondName: runtime.NewInport(nil, runtime.NewSingleInport(tracer, secondIn, runtime.PortAddr{Path: "test/in", Port: secondName}, interceptor)),
		thirdName:  runtime.NewInport(nil, runtime.NewSingleInport(tracer, thirdIn, runtime.PortAddr{Path: "test/in", Port: thirdName}, interceptor)),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, firstIn, secondIn, thirdIn, resultOut
}

func benchNewStringSliceRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return benchNewThreeInputsRuntimeIO("data", "from", "to")
}

func benchNewTernaryRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return benchNewThreeInputsRuntimeIO("if", "then", "else")
}

func benchNewListSliceRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return benchNewThreeInputsRuntimeIO("data", "from", "to")
}
