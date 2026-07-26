package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestStreamZipReceivesBothInputsConcurrently(t *testing.T) {
	io, inChans, outChans := newIO([]string{"left", "right"}, []string{"res"})
	handler, err := (streamZip{}).Create(io, nil)
	if err != nil {
		t.Fatalf("create stream zip handler: %v", err)
	}

	cancel, handlerDone := runHandler(handler)
	t.Cleanup(func() {
		cancel()
		<-handlerDone
	})

	assertSendAcceptedBeforeOtherInput(t, inChans["right"], messages.NewStreamOpenMsg())
	inChans["left"] <- runtime.OrderedMsg{Msg: messages.NewStreamOpenMsg()}
	assertOutputEquals(t, outChans, "res", messages.NewStreamOpenMsg(), []string{"left", "right"})

	assertSendAcceptedBeforeOtherInput(t, inChans["right"], messages.NewStreamDataMsg(messages.NewIntMsg(2)))
	inChans["left"] <- runtime.OrderedMsg{Msg: messages.NewStreamDataMsg(messages.NewIntMsg(1))}
	assertOutputEquals(t, outChans, "res", messages.NewStreamDataMsg(messages.NewStructMsg([]messages.StructField{
		messages.NewStructField("left", messages.NewIntMsg(1)),
		messages.NewStructField("right", messages.NewIntMsg(2)),
	})), []string{"left", "right"})
}

func assertSendAcceptedBeforeOtherInput(t *testing.T, input chan<- runtime.OrderedMsg, msg messages.Msg) {
	t.Helper()

	sent := make(chan struct{})
	go func() {
		input <- runtime.OrderedMsg{Msg: msg}
		close(sent)
	}()

	select {
	case <-sent:
	case <-time.After(time.Second):
		t.Fatal("runtime function did not receive independently available input")
	}
}

func TestArrayPortToStreamReceivesSlotsConcurrently(t *testing.T) {
	io, portInputs, resultOutput := newArrayPortToStreamIO(2)
	handler, err := (arrayPortToStream{}).Create(io, nil)
	if err != nil {
		t.Fatalf("create array port to stream handler: %v", err)
	}

	cancel, handlerDone := runHandler(handler)
	t.Cleanup(func() {
		cancel()
		<-handlerDone
	})

	assertOutputEquals(t, map[string]chan runtime.OrderedMsg{"res": resultOutput}, "res", messages.NewStreamOpenMsg(), []string{"port"})
	assertSendAcceptedBeforeOtherInput(t, portInputs[1], messages.NewIntMsg(2))
	portInputs[0] <- runtime.OrderedMsg{Msg: messages.NewIntMsg(1)}
	assertOutputEquals(t, map[string]chan runtime.OrderedMsg{"res": resultOutput}, "res", messages.NewStreamDataMsg(messages.NewIntMsg(1)), []string{"port[0]"})
	assertOutputEquals(t, map[string]chan runtime.OrderedMsg{"res": resultOutput}, "res", messages.NewStreamDataMsg(messages.NewIntMsg(2)), []string{"port[1]"})
}

func newArrayPortToStreamIO(size int) (runtime.IO, []chan runtime.OrderedMsg, chan runtime.OrderedMsg) {
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}
	inputs := make([]chan runtime.OrderedMsg, size)
	readers := make([]<-chan runtime.OrderedMsg, size)
	for idx := range size {
		inputs[idx] = make(chan runtime.OrderedMsg)
		readers[idx] = inputs[idx]
	}
	resultOutput := make(chan runtime.OrderedMsg, 3)

	input := runtime.NewArrayInport(tracer, readers, runtime.PortAddr{Path: "test/in", Port: "port"}, interceptor)
	output := runtime.NewSingleOutport(tracer, runtime.PortAddr{Path: "test/out", Port: "res"}, interceptor, resultOutput)
	return runtime.IO{
		In: runtime.NewInports(map[string]runtime.Inport{
			"port": runtime.NewInport(input, nil),
		}),
		Out: runtime.NewOutports(map[string]runtime.Outport{
			"res": runtime.NewOutport(output, nil),
		}),
	}, inputs, resultOutput
}
