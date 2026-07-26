package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestListPrependTypedInt(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data", "lst"}, []string{"res"})
	handler, err := (listPrepend{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewIntMsg(1)}
	inChans["lst"] <- runtime.OrderedMsg{Msg: messages.NewListIntMsg([]int64{2, 3})}

	select {
	case got := <-outChans["res"]:
		values, ok := messages.ListAsInts(got.List())
		if !ok || len(values) != 3 || values[0] != 1 || values[1] != 2 || values[2] != 3 {
			t.Fatalf("result = %v, want typed [1 2 3]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}
