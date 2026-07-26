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

func TestDictToMsgs(t *testing.T) {
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

			got := DictToMsgs(tt.dict.Dict())
			if len(got) != len(tt.want) {
				t.Fatalf("DictToMsgs length = %d, want %d", len(got), len(tt.want))
			}
			for key, want := range tt.want {
				if !Equal(got[key], want) {
					t.Fatalf("DictToMsgs[%q] = %v, want %v", key, got[key], want)
				}
			}
		})
	}
}

func TestDictToMsgsTypedResultDoesNotShareStorage(t *testing.T) {
	t.Parallel()

	values := map[string]int64{"one": 1}
	boxed := DictToMsgs(NewDictIntMsg(values).Dict())
	values["one"] = 99

	if got := boxed["one"].Int(); got != 1 {
		t.Fatalf("boxed value = %d, want 1", got)
	}
}

func TestDictToMsgsUntypedReturnsExistingStorage(t *testing.T) {
	t.Parallel()

	values := map[string]Msg{"one": NewIntMsg(1)}
	boxed := DictToMsgs(NewDictMsg(values).Dict())
	boxed["one"] = NewIntMsg(2)

	if got := values["one"].Int(); got != 2 {
		t.Fatalf("untyped storage value = %d, want 2", got)
	}
}

func TestGetDictValueByKey(t *testing.T) {
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

			got, found := GetDictValueByKey(tt.dict.Dict(), tt.key)
			if found != tt.found {
				t.Fatalf("GetDictValueByKey found = %t, want %t", found, tt.found)
			}
			if found && !Equal(got, tt.want) {
				t.Fatalf("GetDictValueByKey value = %v, want %v", got, tt.want)
			}
		})
	}
}
