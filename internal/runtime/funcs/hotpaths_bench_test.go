package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkIntIsLesserHotpath measures the steady-state cost of `int_is_lesser`
// in a tight in-memory runtime loop.
func BenchmarkIntIsLesserHotpath(b *testing.B) {
	runtimeIO, leftInput, rightInput, resultOutput := newBinaryRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := intIsLesser{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	left := runtime.NewIntMsg(7)
	right := runtime.NewIntMsg(42)

	b.ResetTimer()
	for range b.N {
		leftInput <- runtime.OrderedMsg{Msg: left}
		rightInput <- runtime.OrderedMsg{Msg: right}
		<-resultOutput
	}
}

// BenchmarkSelectHotpath measures the select/then fan-in path for the common
// two-slot case.
func BenchmarkSelectHotpath(b *testing.B) {
	runtimeIO, ifInputs, thenInputs, resultOutput := newSelectRuntimeIO(2)
	var zeroConfig runtime.Msg
	handler, err := selector{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	ifMsg := runtime.NewBoolMsg(true)
	then0 := runtime.NewIntMsg(1)
	then1 := runtime.NewIntMsg(2)

	b.ResetTimer()
	for range b.N {
		ifInputs[0] <- runtime.OrderedMsg{Msg: ifMsg}
		thenInputs[0] <- runtime.OrderedMsg{Msg: then0}
		thenInputs[1] <- runtime.OrderedMsg{Msg: then1}
		<-resultOutput
	}
}

// BenchmarkStructHotpath measures building a small struct from three single
// inputs in a tight runtime loop.
func BenchmarkStructHotpath(b *testing.B) {
	runtimeIO, inputs, resultOutput := newStructRuntimeIO([]string{"a", "b", "c"})
	var zeroConfig runtime.Msg
	handler, err := structBuilder{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	msgA := runtime.NewIntMsg(1)
	msgB := runtime.NewIntMsg(2)
	msgC := runtime.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		inputs["a"] <- runtime.OrderedMsg{Msg: msgA}
		inputs["b"] <- runtime.OrderedMsg{Msg: msgB}
		inputs["c"] <- runtime.OrderedMsg{Msg: msgC}
		<-resultOutput
	}
}

// BenchmarkNotEqHotpath measures `not_eq` in a tight in-memory runtime loop.
func BenchmarkNotEqHotpath(b *testing.B) {
	runtimeIO, leftInput, rightInput, resultOutput := newBinaryRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := notEq{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	left := runtime.NewIntMsg(7)
	right := runtime.NewIntMsg(42)

	b.ResetTimer()
	for range b.N {
		leftInput <- runtime.OrderedMsg{Msg: left}
		rightInput <- runtime.OrderedMsg{Msg: right}
		<-resultOutput
	}
}

// BenchmarkStringSliceHotpath measures string slice extraction with fixed
// in-range bounds.
func BenchmarkStringSliceHotpath(b *testing.B) {
	runtimeIO, dataIn, fromIn, toIn, resultOutput := newStringSliceRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := stringSlice{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	data := runtime.NewStringMsg("abcd")
	from := runtime.NewIntMsg(1)
	to := runtime.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		dataIn <- runtime.OrderedMsg{Msg: data}
		fromIn <- runtime.OrderedMsg{Msg: from}
		toIn <- runtime.OrderedMsg{Msg: to}
		<-resultOutput
	}
}

// BenchmarkTernaryElseHotpath measures ternary selector cost on the else path.
func BenchmarkTernaryElseHotpath(b *testing.B) {
	runtimeIO, ifIn, thenIn, elseIn, resultOutput := newTernaryRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := ternarySelector{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	ifMsg := runtime.NewBoolMsg(false)
	thenMsg := runtime.NewIntMsg(42)
	elseMsg := runtime.NewIntMsg(0)

	b.ResetTimer()
	for range b.N {
		ifIn <- runtime.OrderedMsg{Msg: ifMsg}
		thenIn <- runtime.OrderedMsg{Msg: thenMsg}
		elseIn <- runtime.OrderedMsg{Msg: elseMsg}
		<-resultOutput
	}
}

// BenchmarkListSliceHotpath measures list slicing with fixed in-range bounds.
func BenchmarkListSliceHotpath(b *testing.B) {
	runtimeIO, dataIn, fromIn, toIn, resultOutput := newListSliceRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := listSlice{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()
	defer func() {
		cancel()
		<-done
	}()

	data := runtime.NewListMsg([]runtime.Msg{
		runtime.NewIntMsg(1),
		runtime.NewIntMsg(2),
		runtime.NewIntMsg(3),
		runtime.NewIntMsg(4),
	})
	from := runtime.NewIntMsg(1)
	to := runtime.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		dataIn <- runtime.OrderedMsg{Msg: data}
		fromIn <- runtime.OrderedMsg{Msg: from}
		toIn <- runtime.OrderedMsg{Msg: to}
		<-resultOutput
	}
}

func newSelectRuntimeIO(size int) (runtime.IO, []chan runtime.OrderedMsg, []chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
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

	interceptor := runtime.ProdInterceptor{}
	inports := runtime.NewInports(map[string]runtime.Inport{
		"if":   runtime.NewInport(runtime.NewArrayInport(ifRead, runtime.PortAddr{Path: "test/in", Port: "if"}, interceptor), nil),
		"then": runtime.NewInport(runtime.NewArrayInport(thenRead, runtime.PortAddr{Path: "test/in", Port: "then"}, interceptor), nil),
	})

	resultOut := make(chan runtime.OrderedMsg, 1)
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, ifInputs, thenInputs, resultOut
}

func newStructRuntimeIO(names []string) (runtime.IO, map[string]chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	interceptor := runtime.ProdInterceptor{}
	inputs := make(map[string]chan runtime.OrderedMsg, len(names))
	inportsMap := make(map[string]runtime.Inport, len(names))

	for _, name := range names {
		ch := make(chan runtime.OrderedMsg, 1)
		inputs[name] = ch
		inportsMap[name] = runtime.NewInport(nil, runtime.NewSingleInport(ch, runtime.PortAddr{Path: "test/in", Port: name}, interceptor))
	}

	resultOut := make(chan runtime.OrderedMsg, 1)
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: runtime.NewInports(inportsMap), Out: outports}, inputs, resultOut
}

func newStringSliceRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return newThreeInputsRuntimeIO("data", "from", "to")
}

func newTernaryRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return newThreeInputsRuntimeIO("if", "then", "else")
}

func newListSliceRuntimeIO() (
	runtime.IO,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
	chan runtime.OrderedMsg,
) {
	return newThreeInputsRuntimeIO("data", "from", "to")
}

func newThreeInputsRuntimeIO(
	firstName, secondName, thirdName string,
) (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	interceptor := runtime.ProdInterceptor{}
	firstIn := make(chan runtime.OrderedMsg, 1)
	secondIn := make(chan runtime.OrderedMsg, 1)
	thirdIn := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)

	inports := runtime.NewInports(map[string]runtime.Inport{
		firstName:  runtime.NewInport(nil, runtime.NewSingleInport(firstIn, runtime.PortAddr{Path: "test/in", Port: firstName}, interceptor)),
		secondName: runtime.NewInport(nil, runtime.NewSingleInport(secondIn, runtime.PortAddr{Path: "test/in", Port: secondName}, interceptor)),
		thirdName:  runtime.NewInport(nil, runtime.NewSingleInport(thirdIn, runtime.PortAddr{Path: "test/in", Port: thirdName}, interceptor)),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, firstIn, secondIn, thirdIn, resultOut
}
