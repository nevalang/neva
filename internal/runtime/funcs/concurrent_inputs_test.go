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
	handler, err := (formatInt{}).Create(io, nil)
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
		sendDone := make(chan struct{})
		go func(order []string) {
			for _, name := range order {
				switch name {
				case "data":
					inChans[name] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(42)}
				case "base":
					inChans[name] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(10)}
				}
			}
			close(sendDone)
		}(order)

		select {
		case <-sendDone:
		case <-time.After(time.Second):
			t.Fatalf("sending blocked for order %v", order)
		}

		select {
		case got := <-outChans["res"]:
			if !got.Msg.Equal(runtime.NewStringMsg("42")) {
				t.Fatalf("result = %v, want 42", got.Msg)
			}
		case <-time.After(time.Second):
			t.Fatalf("no result for order %v", order)
		}
	}

	cancel()
	<-done
}

func TestTernaryReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newNamedRuntimeIO([]string{"if", "then", "else"}, []string{"res"})
	handler, err := (ternarySelector{}).Create(io, nil)
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
		sendDone := make(chan struct{})
		go func(order []string) {
			for _, name := range order {
				switch name {
				case "if":
					inChans[name] <- runtime.OrderedMsg{Msg: runtime.NewBoolMsg(true)}
				case "then":
					inChans[name] <- runtime.OrderedMsg{Msg: runtime.NewStringMsg("then")}
				case "else":
					inChans[name] <- runtime.OrderedMsg{Msg: runtime.NewStringMsg("else")}
				}
			}
			close(sendDone)
		}(order)

		select {
		case <-sendDone:
		case <-time.After(time.Second):
			t.Fatalf("sending blocked for order %v", order)
		}

		select {
		case got := <-outChans["res"]:
			if !got.Msg.Equal(runtime.NewStringMsg("then")) {
				t.Fatalf("result = %v, want then", got.Msg)
			}
		case <-time.After(time.Second):
			t.Fatalf("no result for order %v", order)
		}
	}

	cancel()
	<-done
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
