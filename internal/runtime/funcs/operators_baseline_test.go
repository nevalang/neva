package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestBinaryOperatorsBehavior(t *testing.T) {
	t.Parallel()

	t.Run("int_add", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, intAdd{}, runtime.NewIntMsg(7), runtime.NewIntMsg(5), runtime.NewIntMsg(12))
	})

	t.Run("float_div", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, floatDiv{}, runtime.NewFloatMsg(9.0), runtime.NewFloatMsg(2.0), runtime.NewFloatMsg(4.5))
	})

	t.Run("string_add", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, stringAdd{}, runtime.NewStringMsg("ne"), runtime.NewStringMsg("va"), runtime.NewStringMsg("neva"))
	})

	t.Run("int_bitwise_xor", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, intBitwiseXor{}, runtime.NewIntMsg(14), runtime.NewIntMsg(11), runtime.NewIntMsg(5))
	})

	t.Run("eq", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, eq{}, runtime.NewStringMsg("same"), runtime.NewStringMsg("same"), runtime.NewBoolMsg(true))
	})

	t.Run("not_eq", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, notEq{}, runtime.NewIntMsg(1), runtime.NewIntMsg(2), runtime.NewBoolMsg(true))
	})

	t.Run("and", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, and{}, runtime.NewBoolMsg(true), runtime.NewBoolMsg(false), runtime.NewBoolMsg(false))
	})

	t.Run("or", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, or{}, runtime.NewBoolMsg(false), runtime.NewBoolMsg(true), runtime.NewBoolMsg(true))
	})

	t.Run("int_is_greater_or_equal", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, intIsGreaterOrEqual{}, runtime.NewIntMsg(10), runtime.NewIntMsg(10), runtime.NewBoolMsg(true))
	})

	t.Run("float_is_lesser", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, floatIsLesser{}, runtime.NewFloatMsg(2.5), runtime.NewFloatMsg(3.5), runtime.NewBoolMsg(true))
	})

	t.Run("string_is_greater", func(t *testing.T) {
		t.Parallel()
		assertBinaryOperatorResult(t, strIsGreater{}, runtime.NewStringMsg("z"), runtime.NewStringMsg("a"), runtime.NewBoolMsg(true))
	})
}

func TestUnaryOperatorsBehavior(t *testing.T) {
	t.Parallel()

	t.Run("int_inc", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intInc{}, runtime.NewIntMsg(41), runtime.NewIntMsg(42))
	})

	t.Run("int_dec", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intDec{}, runtime.NewIntMsg(41), runtime.NewIntMsg(40))
	})

	t.Run("int_neg", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, intNeg{}, runtime.NewIntMsg(8), runtime.NewIntMsg(-8))
	})

	t.Run("float_neg", func(t *testing.T) {
		t.Parallel()
		assertUnaryOperatorResult(t, floatNeg{}, runtime.NewFloatMsg(3.14), runtime.NewFloatMsg(-3.14))
	})
}

func assertBinaryOperatorResult(
	t *testing.T,
	creator runtime.FuncCreator,
	left runtime.Msg,
	right runtime.Msg,
	expected runtime.Msg,
) {
	t.Helper()

	runtimeIO, leftInput, rightInput, resultOutput := newBinaryRuntimeIO()
	handler, err := creator.Create(runtimeIO, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	leftInput <- runtime.OrderedMsg{Msg: left}
	rightInput <- runtime.OrderedMsg{Msg: right}

	result := <-resultOutput
	if !result.Equal(expected) {
		t.Fatalf("result = %v, want %v", result, expected)
	}

	cancel()
	<-done
}

func assertUnaryOperatorResult(
	t *testing.T,
	creator runtime.FuncCreator,
	input runtime.Msg,
	expected runtime.Msg,
) {
	t.Helper()

	runtimeIO, dataInput, resultOutput := newUnaryRuntimeIO()
	handler, err := creator.Create(runtimeIO, nil)
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

func newBinaryRuntimeIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	leftIn := make(chan runtime.OrderedMsg, 1)
	rightIn := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)

	interceptor := runtime.ProdInterceptor{}
	inports := runtime.NewInports(map[string]runtime.Inport{
		"left":  runtime.NewInport(nil, runtime.NewSingleInport(leftIn, runtime.PortAddr{Path: "test/in", Port: "left"}, interceptor)),
		"right": runtime.NewInport(nil, runtime.NewSingleInport(rightIn, runtime.PortAddr{Path: "test/in", Port: "right"}, interceptor)),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, leftIn, rightIn, resultOut
}

func newUnaryRuntimeIO() (runtime.IO, chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	input := make(chan runtime.OrderedMsg, 1)
	resultOut := make(chan runtime.OrderedMsg, 1)

	interceptor := runtime.ProdInterceptor{}
	inports := runtime.NewInports(map[string]runtime.Inport{
		"data": runtime.NewInport(nil, runtime.NewSingleInport(input, runtime.PortAddr{Path: "test/in", Port: "data"}, interceptor)),
	})
	outports := runtime.NewOutports(map[string]runtime.Outport{
		"res": runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOut), nil),
	})

	return runtime.IO{In: inports, Out: outports}, input, resultOut
}
