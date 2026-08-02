package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestTurnWaitsForDoneBeforeReleasingNextData(t *testing.T) {
	io, inChans, outChans := newIO([]string{"data", "done"}, []string{"res"})
	handler, err := (turn{}).Create(io, nil)
	if err != nil {
		t.Fatalf("create turn handler: %v", err)
	}

	cancel, handlerDone := runHandler(handler)
	t.Cleanup(func() {
		cancel()
		<-handlerDone
	})

	first := messages.NewIntMsg(1)
	sendInOrder(t, inChans, []string{"data"}, map[string]messages.Msg{"data": first})
	assertOutputEquals(t, outChans, "res", first, []string{"data"})

	secondSent := make(chan struct{})
	go func() {
		inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewIntMsg(2)}
		close(secondSent)
	}()

	select {
	case <-secondSent:
	case <-time.After(time.Second):
		t.Fatal("turn did not receive the second data")
	}

	select {
	case output := <-outChans["res"]:
		t.Fatalf("turn released %v before done", output)
	case <-time.After(50 * time.Millisecond):
	}

	sendInOrder(t, inChans, []string{"done"}, map[string]messages.Msg{"done": messages.NewBoolMsg(true)})

	assertOutputEquals(t, outChans, "res", messages.NewIntMsg(2), []string{"done", "data"})
}

func TestTurnReceivesDoneAndNextDataConcurrently(t *testing.T) {
	io, inChans, outChans := newIO([]string{"data", "done"}, []string{"res"})
	handler, err := (turn{}).Create(io, nil)
	if err != nil {
		t.Fatalf("create turn handler: %v", err)
	}

	cancel, handlerDone := runHandler(handler)
	t.Cleanup(func() {
		cancel()
		<-handlerDone
	})

	first := messages.NewIntMsg(1)
	sendInOrder(t, inChans, []string{"data"}, map[string]messages.Msg{"data": first})
	assertOutputEquals(t, outChans, "res", first, []string{"data"})

	stopSender := make(chan struct{})
	t.Cleanup(func() { close(stopSender) })
	nextSent := make(chan struct{})
	go func() {
		select {
		case inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewIntMsg(2)}:
		case <-stopSender:
			return
		}

		select {
		case inChans["done"] <- runtime.OrderedMsg{Msg: messages.NewBoolMsg(true)}:
		case <-stopSender:
			return
		}

		close(nextSent)
	}()

	select {
	case <-nextSent:
	case <-time.After(time.Second):
		t.Fatal("turn blocked data before receiving done")
	}

	assertOutputEquals(t, outChans, "res", messages.NewIntMsg(2), []string{"data", "done"})
}
