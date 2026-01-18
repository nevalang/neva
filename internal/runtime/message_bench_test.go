package runtime_test

import (
	"strconv"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

var (
	intSink  int64
	boolSink bool
)

func BenchmarkMsgListIter(b *testing.B) {
	items := make([]runtime.Msg, 1024)
	for i := range items {
		items[i] = runtime.NewIntMsg(int64(i))
	}
	listMsg := runtime.NewListMsg(items)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum int64
		for _, item := range listMsg.List() {
			sum += item.Int()
		}
		intSink = sum
	}
}

func BenchmarkMsgDictLookup(b *testing.B) {
	entries := make(map[string]runtime.Msg, 1024)
	keys := make([]string, 1024)
	for i := range keys {
		key := "k" + strconv.Itoa(i)
		keys[i] = key
		entries[key] = runtime.NewIntMsg(int64(i))
	}
	msg := runtime.NewDictMsg(entries)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum int64
		data := msg.Dict()
		for _, key := range keys {
			sum += data[key].Int()
		}
		intSink = sum
	}
}

func BenchmarkMsgEqualList(b *testing.B) {
	itemsLeft := make([]runtime.Msg, 512)
	itemsRight := make([]runtime.Msg, 512)
	for i := range itemsLeft {
		itemsLeft[i] = runtime.NewStringMsg(strconv.Itoa(i))
		itemsRight[i] = runtime.NewStringMsg(strconv.Itoa(i))
	}
	left := runtime.NewListMsg(itemsLeft)
	right := runtime.NewListMsg(itemsRight)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		boolSink = left.Equal(right)
	}
}

func BenchmarkMsgStructGet(b *testing.B) {
	fields := make([]runtime.StructField, 0, 32)
	for i := 0; i < 32; i++ {
		fields = append(fields, runtime.NewStructField("f"+strconv.Itoa(i), runtime.NewIntMsg(int64(i))))
	}
	msg := runtime.NewStructMsg(fields)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intSink = msg.Struct().Get("f31").Int()
	}
}
