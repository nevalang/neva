package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func TestTurnWaitsForDoneBeforeReceivingNextData(t *testing.T) {
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

	first := runtime.NewIntMsg(1)
	sendInOrder(t, inChans, []string{"data"}, map[string]runtime.Msg{"data": first})
	assertOutputEquals(t, outChans, "res", first, []string{"data"})

	secondSent := make(chan struct{})
	go func() {
		inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(2)}
		close(secondSent)
	}()

	select {
	case <-secondSent:
		t.Fatal("turn received the second data before done")
	case <-time.After(50 * time.Millisecond):
	}

	sendInOrder(t, inChans, []string{"done"}, map[string]runtime.Msg{"done": runtime.NewBoolMsg(true)})

	select {
	case <-secondSent:
	case <-time.After(time.Second):
		t.Fatal("turn did not receive the second data after done")
	}
	assertOutputEquals(t, outChans, "res", runtime.NewIntMsg(2), []string{"done", "data"})
}
