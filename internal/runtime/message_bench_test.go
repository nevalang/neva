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
	entries := make(map[string]int64, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = int64(i)
	}
	return NewDictIntMsg(entries)
}

// BenchmarkMsgListIter measures raw list traversal and integer extraction cost.
func BenchmarkMsgListIter(b *testing.B) {
	for _, size := range []int{8, 64, 512, 1024} {
		b.Run("n="+strconv.Itoa(size), func(b *testing.B) {
			items := make([]int64, size)
			for i := range items {
				items[i] = int64(i)
			}
			listMsg := NewListIntMsg(items)
			data, ok := AsListInts(listMsg.List())
			if !ok {
				b.Fatal("expected int list message")
			}

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				var sum int64
				for _, item := range data {
					sum += item
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
			data, ok := AsDictInts(msg.Dict())
			if !ok {
				b.Fatal("expected int dict message")
			}

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				intSink = data[hotKey]
			}
		})

		b.Run("mixed_n="+strconv.Itoa(size), func(b *testing.B) {
			entries := make(map[string]int64, size)
			keys := make([]string, size)
			for i := range size {
				key := "k" + strconv.Itoa(i)
				keys[i] = key
				entries[key] = int64(i)
			}
			msg := NewDictIntMsg(entries)
			data, ok := AsDictInts(msg.Dict())
			if !ok {
				b.Fatal("expected int dict message")
			}

			b.ReportAllocs()
			b.ResetTimer()
			//nolint:intrange // keeps explicit b.N form for older benchmark style consistency.
			for i := 0; i < b.N; i++ {
				var sum int64
				for _, key := range keys {
					sum += data[key]
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
	items := make([]int64, size)
	for i := range items {
		items[i] = int64(i)
	}
	msg := NewListIntMsg(items)
	data, ok := AsListInts(msg.List())
	if !ok {
		b.Fatal("expected int list message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var sum int64
		for _, item := range data {
			sum += item
		}
		intSink = sum
	}
}

func benchListIterFloat(b *testing.B, size int) {
	b.Helper()
	items := make([]float64, size)
	for i := range items {
		items[i] = float64(i) + 0.25
	}
	msg := NewListFloatMsg(items)
	data, ok := AsListFloats(msg.List())
	if !ok {
		b.Fatal("expected float list message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var sum float64
		for _, item := range data {
			sum += item
		}
		floatSink = sum
	}
}

func benchListIterBool(b *testing.B, size int) {
	b.Helper()
	items := make([]bool, size)
	for i := range items {
		items[i] = i%2 == 0
	}
	msg := NewListBoolMsg(items)
	data, ok := AsListBools(msg.List())
	if !ok {
		b.Fatal("expected bool list message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var count int64
		for _, item := range data {
			if item {
				count++
			}
		}
		intSink = count
	}
}

func benchListIterString(b *testing.B, size int) {
	b.Helper()
	items := make([]string, size)
	for i := range items {
		items[i] = "v" + strconv.Itoa(i)
	}
	msg := NewListStringMsg(items)
	data, ok := AsListStrings(msg.List())
	if !ok {
		b.Fatal("expected string list message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var total int64
		for _, item := range data {
			total += int64(len(item))
		}
		intSink = total
	}
}

func benchDictLookupInt(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]int64, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = int64(i)
	}
	msg := NewDictIntMsg(entries)
	data, ok := AsDictInts(msg.Dict())
	if !ok {
		b.Fatal("expected int dict message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		intSink = data[hotKey]
	}
}

func benchDictLookupFloat(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]float64, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = float64(i) + 0.25
	}
	msg := NewDictFloatMsg(entries)
	data, ok := AsDictFloats(msg.Dict())
	if !ok {
		b.Fatal("expected float dict message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		floatSink = data[hotKey]
	}
}

func benchDictLookupBool(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]bool, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = i%2 == 0
	}
	msg := NewDictBoolMsg(entries)
	data, ok := AsDictBools(msg.Dict())
	if !ok {
		b.Fatal("expected bool dict message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		boolSink = data[hotKey]
	}
}

func benchDictLookupString(b *testing.B, size int, hotKey string) {
	b.Helper()
	entries := make(map[string]string, size)
	for i := range size {
		entries["k"+strconv.Itoa(i)] = "v" + strconv.Itoa(i)
	}
	msg := NewDictStringMsg(entries)
	data, ok := AsDictStrings(msg.Dict())
	if !ok {
		b.Fatal("expected string dict message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		stringSink = data[hotKey]
	}
}
