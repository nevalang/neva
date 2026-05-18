package runtime

import (
	"strconv"
	"testing"
)

var (
	intSink    int64
	boolSink   bool
	floatSink  float64
	stringSink string
)

//nolint:ireturn // benchmark helper returns runtime.Msg by design.
func makeDictMsg(size int) Msg {
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewIntMsg(int64(i))
	}
	return NewDictMsg(entries)
}

// BenchmarkMsgListIter measures raw list traversal and integer extraction cost.
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
				for _, item := range listMsg.List().Msgs() {
					sum += item.Int()
				}
				intSink = sum
			}
		})
	}
}

// BenchmarkMsgDictLookup measures dictionary lookup in hot-key and mixed-keys modes.
func BenchmarkMsgDictLookup(b *testing.B) {
	for _, size := range []int{16, 128, 1024} {
		b.Run("hot_n="+strconv.Itoa(size), func(b *testing.B) {
			msg := makeDictMsg(size)
			hotKey := "k" + strconv.Itoa(size-1)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				intSink = msg.Dict().Msgs()[hotKey].Int()
			}
		})

		b.Run("mixed_n="+strconv.Itoa(size), func(b *testing.B) {
			entries := make(map[string]Msg, size)
			keys := make([]string, size)
			for i := range size {
				key := "k" + strconv.Itoa(i)
				keys[i] = key
				entries[key] = NewIntMsg(int64(i))
			}
			msg := NewDictMsg(entries)

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				var sum int64
				data := msg.Dict().Msgs()
				for _, key := range keys {
					sum += data[key].Int()
				}
				intSink = sum
			}
		})
	}
}

// BenchmarkMsgEqualList measures list equality for equal and early-unequal inputs.
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

// BenchmarkMsgStructGet measures repeated field lookup in a medium struct.
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

// BenchmarkMsgListIterScalars compares list traversal cost across scalar payload kinds.
func BenchmarkMsgListIterScalars(b *testing.B) {
	for _, size := range []int{8, 64, 512, 1024} {
		b.Run("int_n="+strconv.Itoa(size), func(b *testing.B) { benchListIterInt(b, size) })
		b.Run("float_n="+strconv.Itoa(size), func(b *testing.B) { benchListIterFloat(b, size) })
		b.Run("bool_n="+strconv.Itoa(size), func(b *testing.B) { benchListIterBool(b, size) })
		b.Run("string_n="+strconv.Itoa(size), func(b *testing.B) { benchListIterString(b, size) })
	}
}

// BenchmarkMsgDictLookupScalars compares hot-key lookup cost across scalar payload kinds.
func BenchmarkMsgDictLookupScalars(b *testing.B) {
	for _, size := range []int{16, 128, 1024} {
		hotKey := "k" + strconv.Itoa(size-1)
		b.Run("int_hot_n="+strconv.Itoa(size), func(b *testing.B) { benchDictLookupInt(b, size, hotKey) })
		b.Run("float_hot_n="+strconv.Itoa(size), func(b *testing.B) { benchDictLookupFloat(b, size, hotKey) })
		b.Run("bool_hot_n="+strconv.Itoa(size), func(b *testing.B) { benchDictLookupBool(b, size, hotKey) })
		b.Run("string_hot_n="+strconv.Itoa(size), func(b *testing.B) { benchDictLookupString(b, size, hotKey) })
	}
}

func benchListIterInt(b *testing.B, size int) {
	b.Helper()
	items := make([]Msg, size)
	for i := range items {
		items[i] = NewIntMsg(int64(i))
	}
	msg := NewListMsg(items)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var sum int64
		for _, item := range msg.List() {
			sum += item.Int()
		}
		intSink = sum
	}
}

func benchListIterFloat(b *testing.B, size int) {
	b.Helper()
	items := make([]Msg, size)
	for i := range items {
		items[i] = NewFloatMsg(float64(i) + 0.25)
	}
	msg := NewListMsg(items)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var sum float64
		for _, item := range msg.List() {
			sum += item.Float()
		}
		floatSink = sum
	}
}

func benchListIterBool(b *testing.B, size int) {
	b.Helper()
	items := make([]Msg, size)
	for i := range items {
		items[i] = NewBoolMsg(i%2 == 0)
	}
	msg := NewListMsg(items)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var count int64
		for _, item := range msg.List() {
			if item.Bool() {
				count++
			}
		}
		intSink = count
	}
}

func benchListIterString(b *testing.B, size int) {
	b.Helper()
	items := make([]Msg, size)
	for i := range items {
		items[i] = NewStringMsg("v" + strconv.Itoa(i))
	}
	msg := NewListMsg(items)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var total int64
		for _, item := range msg.List() {
			total += int64(len(item.Str()))
		}
		intSink = total
	}
}

func benchDictLookupInt(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewIntMsg(int64(i))
	}
	msg := NewDictMsg(entries)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		intSink = msg.Dict()[hotKey].Int()
	}
}

func benchDictLookupFloat(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewFloatMsg(float64(i) + 0.25)
	}
	msg := NewDictMsg(entries)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		floatSink = msg.Dict()[hotKey].Float()
	}
}

func benchDictLookupBool(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewBoolMsg(i%2 == 0)
	}
	msg := NewDictMsg(entries)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		boolSink = msg.Dict()[hotKey].Bool()
	}
}

func benchDictLookupString(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]Msg, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = NewStringMsg("v" + strconv.Itoa(i))
	}
	msg := NewDictMsg(entries)
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		stringSink = msg.Dict()[hotKey].Str()
	}
}
