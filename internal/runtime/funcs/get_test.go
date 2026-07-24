package funcs

import (
	"context"
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

func TestDictValueByKeyTypedMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dict runtime.Msg
	}{
		{"bool", runtime.NewDictBoolMsg(map[string]bool{"present": true})},
		{"int", runtime.NewDictIntMsg(map[string]int64{"present": 1})},
		{"float", runtime.NewDictFloatMsg(map[string]float64{"present": 1.5})},
		{"string", runtime.NewDictStringMsg(map[string]string{"present": "x"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, found := dictValueByKey(tt.dict.Dict(), "missing"); found {
				t.Fatal("missing key was found")
			}
		})
	}
}

func TestGetDictValueTypedMissSendsError(t *testing.T) {
	t.Parallel()

	io, inputs, outputs := newIO([]string{"dict", "key"}, []string{"res", "err"})
	handler, err := (getDictValue{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { handler(ctx); close(done) }()
	defer func() { cancel(); <-done }()

	go func() {
		inputs["dict"] <- runtime.OrderedMsg{Msg: runtime.NewDictIntMsg(map[string]int64{"present": 1})}
	}()
	go func() { inputs["key"] <- runtime.OrderedMsg{Msg: runtime.NewStringMsg("missing")} }()

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
