package messages

import (
	"strconv"
	"testing"
)

var (
	intSink    int64
	boolSink   bool
	floatSink  float64
	stringSink string
	msgSink    Msg
	listSink   []Msg
	dictSink   map[string]Msg
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
			data, ok := ListAsInts(listMsg.List())
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
//
//nolint:gocognit // Benchmark variants are intentionally kept adjacent for comparison.
func BenchmarkMsgDictLookup(b *testing.B) {
	for _, size := range []int{16, 128, 1024} {
		b.Run("hot_n="+strconv.Itoa(size), func(b *testing.B) {
			msg := makeDictMsg(size)
			hotKey := "k" + strconv.Itoa(size-1)
			data, ok := DictAsInts(msg.Dict())
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
			data, ok := DictAsInts(msg.Dict())
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

// BenchmarkDictGetValueByKey measures lookup through the representation-independent API.
func BenchmarkDictGetValueByKey(b *testing.B) {
	dict := NewDictIntMsg(map[string]int64{"answer": 42}).Dict()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		msgSink, _ = DictGetValueByKey(dict, "answer")
	}
}

// BenchmarkListToMessageSlice measures the explicit boxing boundary for typed lists.
func BenchmarkListToMessageSlice(b *testing.B) {
	values := make([]int64, 128)
	for i := range values {
		values[i] = int64(i)
	}
	list := NewListIntMsg(values).List()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		listSink = ListToMessageSlice(list)
	}
}

// BenchmarkDictToMessageMap measures the explicit boxing boundary for typed dictionaries.
func BenchmarkDictToMessageMap(b *testing.B) {
	values := make(map[string]int64, 128)
	for i := range 128 {
		values["k"+strconv.Itoa(i)] = int64(i)
	}
	dict := NewDictIntMsg(values).Dict()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		dictSink = DictToMessageMap(dict)
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
				boolSink = Equal(left, right)
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
				boolSink = Equal(left, right)
			}
		})
	}
}

// BenchmarkMsgEqualDict compares equal dictionary representations without
// materializing scalar maps as message maps.
func BenchmarkMsgEqualDict(b *testing.B) {
	for _, size := range []int{16, 128, 512} {
		values := make(map[string]int64, size)
		boxedValues := make(map[string]Msg, size)
		stringValues := make(map[string]string, size)
		for i := range size {
			key := "k" + strconv.Itoa(i)
			values[key] = int64(i)
			boxedValues[key] = NewIntMsg(int64(i))
			stringValues[key] = strconv.Itoa(i)
		}

		b.Run("typed_equal_n="+strconv.Itoa(size), func(b *testing.B) {
			left := NewDictIntMsg(values)
			right := NewDictIntMsg(values)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				boolSink = Equal(left, right)
			}
		})

		b.Run("typed_untyped_equal_n="+strconv.Itoa(size), func(b *testing.B) {
			left := NewDictIntMsg(values)
			right := NewDictMsg(boxedValues)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				boolSink = Equal(left, right)
			}
		})

		b.Run("typed_different_kind_n="+strconv.Itoa(size), func(b *testing.B) {
			left := NewDictIntMsg(values)
			right := NewDictStringMsg(stringValues)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				boolSink = Equal(left, right)
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
		intSink = StructGet(msg.Struct(), "f31").Int()
	}
}

// BenchmarkListSlice measures direct immutable slicing without transport cost.
func BenchmarkListSlice(b *testing.B) {
	values := make([]int64, 512)
	for i := range values {
		values[i] = int64(i)
	}
	list := NewListIntMsg(values).List()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		msgSink = ListSlice(list, 64, 448)
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

// BenchmarkMsgScalarContainerConstruction measures construction of the same
// scalar list/dict payloads that compiler-generated literals produce.
func BenchmarkMsgScalarContainerConstruction(b *testing.B) {
	for _, size := range []int{8, 64, 512} {
		b.Run("list_int_n="+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				items := make([]int64, size)
				for i := range items {
					items[i] = int64(i)
				}
				msgSink = NewListIntMsg(items)
			}
		})

		b.Run("dict_int_n="+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				entries := make(map[string]int64, size)
				for i := range size {
					entries["k"+strconv.Itoa(i)] = int64(i)
				}
				msgSink = NewDictIntMsg(entries)
			}
		})
	}
}

func benchListIterInt(b *testing.B, size int) {
	b.Helper()
	items := make([]int64, size)
	for i := range items {
		items[i] = int64(i)
	}
	msg := NewListIntMsg(items)
	data, ok := ListAsInts(msg.List())
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
	data, ok := ListAsFloats(msg.List())
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
	data, ok := ListAsBools(msg.List())
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
	data, ok := ListAsStrings(msg.List())
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
	data, ok := DictAsInts(msg.Dict())
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
	data, ok := DictAsFloats(msg.Dict())
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
	data, ok := DictAsBools(msg.Dict())
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
	data, ok := DictAsStrings(msg.Dict())
	if !ok {
		b.Fatal("expected string dict message")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		stringSink = data[hotKey]
	}
}
