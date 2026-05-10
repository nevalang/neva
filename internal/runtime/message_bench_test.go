package runtime

import (
	"strconv"
	"testing"
)

var (
	intSink  int64
	boolSink bool
)

func makeDictMsg(size int) Msg {
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewIntMsg(int64(i))
	}
	return NewDictMsg(entries)
}

func makeDictMsgAndKeys(size int) (Msg, []string) {
	entries := make(map[string]Msg, size)
	keys := make([]string, size)
	for i := range size {
		key := "k" + strconv.Itoa(i)
		keys[i] = key
		entries[key] = NewIntMsg(int64(i))
	}
	return NewDictMsg(entries), keys
}

func BenchmarkMsgListIter(b *testing.B) {
	for _, size := range []int{8, 64, 512, 1024} {
		b.Run("n="+strconv.Itoa(size), func(b *testing.B) {
			items := make([]Msg, size)
			for i := range items {
				items[i] = NewIntMsg(int64(i))
			}
			listMsg := NewListMsg(items)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				var sum int64
				for _, item := range listMsg.List() {
					sum += item.Int()
				}
				intSink = sum
			}
		})
	}
}

func BenchmarkMsgDictLookup(b *testing.B) {
	for _, size := range []int{16, 128, 1024} {
		b.Run("hot_n="+strconv.Itoa(size), func(b *testing.B) {
			msg := makeDictMsg(size)
			hotKey := "k" + strconv.Itoa(size-1)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				intSink = msg.Dict()[hotKey].Int()
			}
		})

		b.Run("mixed_n="+strconv.Itoa(size), func(b *testing.B) {
			msg, keys := makeDictMsgAndKeys(size)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				var sum int64
				data := msg.Dict()
				for _, key := range keys {
					sum += data[key].Int()
				}
				intSink = sum
			}
		})
	}
}

func BenchmarkMsgEqualList(b *testing.B) {
	for _, size := range []int{16, 128, 512} {
		b.Run("equal_n="+strconv.Itoa(size), func(b *testing.B) {
			itemsLeft := make([]Msg, size)
			itemsRight := make([]Msg, size)
			for i := range itemsLeft {
				val := strconv.Itoa(i)
				itemsLeft[i] = NewStringMsg(val)
				itemsRight[i] = NewStringMsg(val)
			}
			left := NewListMsg(itemsLeft)
			right := NewListMsg(itemsRight)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				boolSink = left.Equal(right)
			}
		})

		b.Run("unequal_early_n="+strconv.Itoa(size), func(b *testing.B) {
			itemsLeft := make([]Msg, size)
			itemsRight := make([]Msg, size)
			for i := range itemsLeft {
				val := strconv.Itoa(i)
				itemsLeft[i] = NewStringMsg(val)
				itemsRight[i] = NewStringMsg(val)
			}
			itemsRight[0] = NewStringMsg("x")
			left := NewListMsg(itemsLeft)
			right := NewListMsg(itemsRight)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				boolSink = left.Equal(right)
			}
		})
	}
}

func BenchmarkMsgStructGet(b *testing.B) {
	fields := make([]StructField, 0, 32)
	for i := range 32 {
		fields = append(fields, NewStructField("f"+strconv.Itoa(i), NewIntMsg(int64(i))))
	}
	msg := NewStructMsg(fields)

	b.ReportAllocs()
	b.ResetTimer()
	//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
	for i := 0; i < b.N; i++ {
		intSink = msg.Struct().Get("f31").Int()
	}
}

func BenchmarkMsgUnionUnbox(b *testing.B) {
	withData := NewUnionMsg("Some", NewIntMsg(42))
	noData := NewUnionMsgNoData("None")

	b.Run("with_data", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
		for i := 0; i < b.N; i++ {
			u := withData.Union()
			if !u.HasData() {
				b.Fatalf("expected payload for with_data")
			}
			intSink = u.Data().Int()
		}
	})

	b.Run("no_data", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
		for i := 0; i < b.N; i++ {
			u := noData.Union()
			boolSink = u.HasData()
		}
	})
}
