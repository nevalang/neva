package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestListConcatTypedInt(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"left", "right"}, []string{"res"})
	handler, err := (listConcat{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["left"] <- runtime.OrderedMsg{Msg: messages.NewListIntMsg([]int64{1, 2})}
	inChans["right"] <- runtime.OrderedMsg{Msg: messages.NewListIntMsg([]int64{3, 4})}

	select {
	case got := <-outChans["res"]:
		values, ok := messages.ListAsInts(got.List())
		if !ok || len(values) != 4 || values[0] != 1 || values[3] != 4 {
			t.Fatalf("result = %v, want typed [1 2 3 4]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}
