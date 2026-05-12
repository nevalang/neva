package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFormatIntReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newNamedRuntimeIO([]string{"data", "base"}, []string{"res"})
	handler, err := (formatInt{}).Create(io, runtime.Msg{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	for _, order := range [][]string{{"data", "base"}, {"base", "data"}} {
		sendInOrder(t, inChans, order, map[string]runtime.Msg{
			"data": runtime.NewIntMsg(42),
			"base": runtime.NewIntMsg(10),
		})
		assertOutputEquals(t, outChans, "res", runtime.NewStringMsg("42"), order)
	}

	cancel()
	<-done
}

func TestTernaryReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newNamedRuntimeIO([]string{"if", "then", "else"}, []string{"res"})
	handler, err := (ternarySelector{}).Create(io, runtime.Msg{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	for _, order := range [][]string{{"if", "then", "else"}, {"else", "then", "if"}} {
		sendInOrder(t, inChans, order, map[string]runtime.Msg{
			"if":   runtime.NewBoolMsg(true),
			"then": runtime.NewStringMsg("then"),
			"else": runtime.NewStringMsg("else"),
		})
		assertOutputEquals(t, outChans, "res", runtime.NewStringMsg("then"), order)
	}

	cancel()
	<-done
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

func newNamedRuntimeIO(inNames []string, outNames []string) (runtime.IO, map[string]chan runtime.OrderedMsg, map[string]chan runtime.OrderedMsg) {
	interceptor := runtime.ProdInterceptor{}
	inports := make(map[string]runtime.Inport, len(inNames))
	outports := make(map[string]runtime.Outport, len(outNames))
	inChans := make(map[string]chan runtime.OrderedMsg, len(inNames))
	outChans := make(map[string]chan runtime.OrderedMsg, len(outNames))

	for _, name := range inNames {
		ch := make(chan runtime.OrderedMsg)
		inChans[name] = ch
		inports[name] = runtime.NewInport(nil, runtime.NewSingleInport(ch, runtime.PortAddr{Path: "test/in", Port: name}, interceptor))
	}

	for _, name := range outNames {
		ch := make(chan runtime.OrderedMsg, 1)
		outChans[name] = ch
		outports[name] = runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: "test/out", Port: name}, interceptor, ch), nil)
	}

	return runtime.IO{In: runtime.NewInports(inports), Out: runtime.NewOutports(outports)}, inChans, outChans
}
