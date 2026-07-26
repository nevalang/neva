package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestStreamToListMaterializesTypedInts(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data"}, []string{"res"})
	handler, err := (streamToList{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamOpenMsg()}
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamDataMsg(messages.NewIntMsg(1))}
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamDataMsg(messages.NewIntMsg(2))}
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamCloseMsg()}

	select {
	case got := <-outChans["res"]:
		values, ok := messages.ListAsInts(got.List())
		if !ok || len(values) != 2 || values[0] != 1 || values[1] != 2 {
			t.Fatalf("result = %v, want typed [1 2]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}

func TestStreamToDictMaterializesTypedStrings(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data"}, []string{"res"})
	handler, err := (streamToDict{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	entry := messages.NewStructMsg([]messages.StructField{
		messages.NewStructField("key", messages.NewStringMsg("one")),
		messages.NewStructField("value", messages.NewStringMsg("value")),
	})
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamOpenMsg()}
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamDataMsg(entry)}
	inChans["data"] <- runtime.OrderedMsg{Msg: messages.NewStreamCloseMsg()}

	select {
	case got := <-outChans["res"]:
		values, ok := messages.DictAsStrings(got.Dict())
		if !ok || values["one"] != "value" {
			t.Fatalf("result = %v, want typed {one: value}", got)
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}
