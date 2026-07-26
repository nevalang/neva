package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

func TestGetDictValueTypedMissSendsError(t *testing.T) {
	t.Parallel()

	io, inputs, outputs := newIO([]string{"dict", "key"}, []string{"res", "err"})
	handler, err := (dictGetByKey{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()
	defer func() { cancel(); <-done }()

	go func() {
		inputs["dict"] <- runtime.OrderedMsg{Msg: messages.NewDictIntMsg(map[string]int64{"present": 1})}
	}()
	go func() { inputs["key"] <- runtime.OrderedMsg{Msg: messages.NewStringMsg("missing")} }()

	select {
	case got := <-outputs["err"]:
		if text := got.Struct().Get("text").Str(); text != "Key not found in dictionary" {
			t.Fatalf("error text = %q", text)
		}
	case <-time.After(time.Second):
		t.Fatal("expected error output")
	}
	select {
	case <-outputs["res"]:
		t.Fatal("unexpected result")
	case <-time.After(20 * time.Millisecond):
	}
}

func TestGetDictValueSendsTypedValue(t *testing.T) {
	t.Parallel()

	io, inputs, outputs := newIO([]string{"dict", "key"}, []string{"res", "err"})
	handler, err := (dictGetByKey{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()
	defer func() { cancel(); <-done }()

	go func() {
		inputs["dict"] <- runtime.OrderedMsg{Msg: messages.NewDictBoolMsg(map[string]bool{"present": false})}
	}()
	go func() { inputs["key"] <- runtime.OrderedMsg{Msg: messages.NewStringMsg("present")} }()

	select {
	case got := <-outputs["res"]:
		if got.Bool() {
			t.Fatal("result = true, want false")
		}
	case <-time.After(time.Second):
		t.Fatal("expected result output")
	}
	select {
	case <-outputs["err"]:
		t.Fatal("unexpected error output")
	case <-time.After(20 * time.Millisecond):
	}
}
