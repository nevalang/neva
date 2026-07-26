package messages

import "testing"

func TestDictLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dict Msg
		name string
		want int
	}{
		{name: "untyped", dict: NewDictMsg(map[string]Msg{"one": NewIntMsg(1)}), want: 1},
		{name: "bool", dict: NewDictBoolMsg(map[string]bool{"one": true, "two": false}), want: 2},
		{name: "int", dict: NewDictIntMsg(map[string]int64{}), want: 0},
		{name: "float", dict: NewDictFloatMsg(map[string]float64{"one": 1}), want: 1},
		{name: "string", dict: NewDictStringMsg(map[string]string{"one": "one", "two": "two"}), want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DictLen(tt.dict.Dict()); got != tt.want {
				t.Fatalf("DictLen() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDictToMessageMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dict Msg
		want map[string]Msg
		name string
	}{
		{
			name: "untyped",
			dict: NewDictMsg(map[string]Msg{"one": NewStringMsg("one"), "two": NewIntMsg(2)}),
			want: map[string]Msg{"one": NewStringMsg("one"), "two": NewIntMsg(2)},
		},
		{
			name: "bool",
			dict: NewDictBoolMsg(map[string]bool{"one": true}),
			want: map[string]Msg{"one": NewBoolMsg(true)},
		},
		{
			name: "int",
			dict: NewDictIntMsg(map[string]int64{"one": 1}),
			want: map[string]Msg{"one": NewIntMsg(1)},
		},
		{
			name: "float",
			dict: NewDictFloatMsg(map[string]float64{"one": 1.5}),
			want: map[string]Msg{"one": NewFloatMsg(1.5)},
		},
		{
			name: "string",
			dict: NewDictStringMsg(map[string]string{"one": "one"}),
			want: map[string]Msg{"one": NewStringMsg("one")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DictToMessageMap(tt.dict.Dict())
			if len(got) != len(tt.want) {
				t.Fatalf("DictToMessageMap length = %d, want %d", len(got), len(tt.want))
			}
			for key, want := range tt.want {
				if !Equal(got[key], want) {
					t.Fatalf("DictToMessageMap[%q] = %v, want %v", key, got[key], want)
				}
			}
		})
	}
}

func TestDictToMessageMapTypedResultDoesNotShareStorage(t *testing.T) {
	t.Parallel()

	values := map[string]int64{"one": 1}
	boxed := DictToMessageMap(NewDictIntMsg(values).Dict())
	values["one"] = 99

	if got := boxed["one"].Int(); got != 1 {
		t.Fatalf("boxed value = %d, want 1", got)
	}
}

func TestDictToMessageMapUntypedReturnsExistingStorage(t *testing.T) {
	t.Parallel()

	values := map[string]Msg{"one": NewIntMsg(1)}
	boxed := DictToMessageMap(NewDictMsg(values).Dict())
	boxed["one"] = NewIntMsg(2)

	if got := values["one"].Int(); got != 2 {
		t.Fatalf("untyped storage value = %d, want 2", got)
	}
}

func TestDictFromMessagesUsesScalarStorage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		values map[string]Msg
		want   Msg
		name   string
	}{
		{
			name:   "bool",
			values: map[string]Msg{"a": NewBoolMsg(true), "b": NewBoolMsg(false)},
			want:   NewDictBoolMsg(map[string]bool{"a": true, "b": false}),
		},
		{
			name:   "int",
			values: map[string]Msg{"a": NewIntMsg(1), "b": NewIntMsg(2)},
			want:   NewDictIntMsg(map[string]int64{"a": 1, "b": 2}),
		},
		{
			name:   "float",
			values: map[string]Msg{"a": NewFloatMsg(1.5), "b": NewFloatMsg(2.5)},
			want:   NewDictFloatMsg(map[string]float64{"a": 1.5, "b": 2.5}),
		},
		{
			name:   "string",
			values: map[string]Msg{"a": NewStringMsg("one"), "b": NewStringMsg("two")},
			want:   NewDictStringMsg(map[string]string{"a": "one", "b": "two"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DictFromMessages(tt.values)
			if _, ok := DictAsUntyped(got.Dict()); ok {
				t.Fatal("DictFromMessages result did not use scalar storage")
			}
			if !Equal(got, tt.want) {
				t.Fatalf("DictFromMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictFromMessagesFallsBackToUntypedStorage(t *testing.T) {
	t.Parallel()

	for _, values := range []map[string]Msg{
		nil,
		{"a": NewIntMsg(1), "b": NewStringMsg("two")},
		{"a": NewListIntMsg([]int64{1})},
	} {
		got := DictFromMessages(values)
		if _, ok := DictAsUntyped(got.Dict()); !ok {
			t.Fatalf("DictFromMessages(%v) did not use untyped storage", values)
		}
	}
}

func TestDictFromMessagesDoesNotShareScalarStorage(t *testing.T) {
	t.Parallel()

	values := map[string]Msg{"one": NewIntMsg(1)}
	result := DictFromMessages(values)
	values["one"] = NewIntMsg(99)

	ints, ok := DictAsInts(result.Dict())
	if !ok || ints["one"] != 1 {
		t.Fatalf("DictFromMessages result = %v, want typed {one: 1}", result)
	}
}

func TestDictGetValueByKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dict  Msg
		want  Msg
		name  string
		key   string
		found bool
	}{
		{
			name:  "untyped value",
			dict:  NewDictMsg(map[string]Msg{"one": NewStringMsg("one")}),
			key:   "one",
			want:  NewStringMsg("one"),
			found: true,
		},
		{
			name:  "bool false value",
			dict:  NewDictBoolMsg(map[string]bool{"one": false}),
			key:   "one",
			want:  NewBoolMsg(false),
			found: true,
		},
		{
			name:  "int value",
			dict:  NewDictIntMsg(map[string]int64{"one": 1}),
			key:   "one",
			want:  NewIntMsg(1),
			found: true,
		},
		{
			name:  "float value",
			dict:  NewDictFloatMsg(map[string]float64{"one": 1.5}),
			key:   "one",
			want:  NewFloatMsg(1.5),
			found: true,
		},
		{
			name:  "string value",
			dict:  NewDictStringMsg(map[string]string{"one": "one"}),
			key:   "one",
			want:  NewStringMsg("one"),
			found: true,
		},
		{
			name:  "typed miss",
			dict:  NewDictIntMsg(map[string]int64{"one": 1}),
			key:   "missing",
			found: false,
		},
		{
			name:  "untyped miss",
			dict:  NewDictMsg(map[string]Msg{"one": NewIntMsg(1)}),
			key:   "missing",
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := DictGetValueByKey(tt.dict.Dict(), tt.key)
			if found != tt.found {
				t.Fatalf("DictGetValueByKey found = %t, want %t", found, tt.found)
			}
			if found && !Equal(got, tt.want) {
				t.Fatalf("DictGetValueByKey value = %v, want %v", got, tt.want)
			}
		})
	}
}
