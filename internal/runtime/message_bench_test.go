package runtime

import (
	"context"
	"strconv"
	"testing"
)

var (
	intSink  int64
	boolSink bool
)

func BenchmarkMsgListIter(b *testing.B) {
	items := make([]Msg, 1024)
	for i := range items {
		items[i] = NewIntMsg(int64(i))
	}
	listMsg := NewListMsg(items)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		var sum int64
		for _, item := range listMsg.List() {
			sum += item.Int()
		}
		intSink = sum
	}
}

func BenchmarkMsgDictLookup(b *testing.B) {
	entries := make(map[string]Msg, 1024)
	keys := make([]string, 1024)
	for i := range keys {
		key := "k" + strconv.Itoa(i)
		keys[i] = key
		entries[key] = NewIntMsg(int64(i))
	}
	msg := NewDictMsg(entries)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		var sum int64
		data := msg.Dict()
		for _, key := range keys {
			sum += data[key].Int()
		}
		intSink = sum
	}
}

func BenchmarkMsgEqualList(b *testing.B) {
	itemsLeft := make([]Msg, 512)
	itemsRight := make([]Msg, 512)
	for i := range itemsLeft {
		itemsLeft[i] = NewStringMsg(strconv.Itoa(i))
		itemsRight[i] = NewStringMsg(strconv.Itoa(i))
	}
	left := NewListMsg(itemsLeft)
	right := NewListMsg(itemsRight)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		boolSink = left.Equal(right)
	}
}

func BenchmarkMsgStructGet(b *testing.B) {
	fields := make([]StructField, 0, 32)
	for i := 0; i < 32; i++ {
		fields = append(fields, NewStructField("f"+strconv.Itoa(i), NewIntMsg(int64(i))))
	}
	msg := NewStructMsg(fields)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		intSink = msg.Struct().Get("f31").Int()
	}
}

func BenchmarkSinglePortRoundTrip(b *testing.B) {
	ctx := context.Background()
	channel := make(chan OrderedMsg, 1)
	addr := PortAddr{Path: "bench/runtime", Port: "msg"}
	outport := NewSingleOutport(addr, ProdInterceptor{}, channel)
	inport := NewSingleInport(channel, addr, ProdInterceptor{})
	msg := NewIntMsg(42)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if !outport.Send(ctx, msg) {
			b.Fatal("send failed")
		}
		received, ok := inport.Receive(ctx)
		if !ok {
			b.Fatal("receive failed")
		}
		intSink = received.Int()
	}
}
