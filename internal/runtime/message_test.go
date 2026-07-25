package runtime

import (
	"testing"
)

func mustMarshal(t *testing.T, msg Msg) []byte {
	t.Helper()
	marshaler, ok := msg.(interface{ MarshalJSON() ([]byte, error) })
	if !ok {
		t.Fatalf("message type %T does not implement MarshalJSON", msg)
	}
	b, err := marshaler.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	return b
}

func TestDictMsgMarshalJSONPreservesStringValues(t *testing.T) {
	msg := NewDictMsg(map[string]Msg{
		"text": NewStringMsg(`a:"b,c\d`),
		"nums": NewListMsg([]Msg{
			NewIntMsg(1),
			NewIntMsg(2),
		}),
	})

	b := mustMarshal(t, msg)
	if got, want := string(b), `{"nums": [1, 2], "text": "a:\"b,c\\d"}`; got != want {
		t.Fatalf("MarshalJSON() = %q, want %q", got, want)
	}
}

func TestStructMsgMarshalJSONPreservesStringValues(t *testing.T) {
	msg := NewStructMsg([]StructField{
		NewStructField("text", NewStringMsg(`a:"b,c\d`)),
		NewStructField("nums", NewListMsg([]Msg{
			NewIntMsg(1),
			NewIntMsg(2),
		})),
	})

	b, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if got, want := string(b), `{"nums": [1, 2], "text": "a:\"b,c\\d"}`; got != want {
		t.Fatalf("MarshalJSON() = %q, want %q", got, want)
	}
}

func TestUnionMsgStringTagOnly(t *testing.T) {
	msg := NewUnionMsg("Friday", nil)
	if got, want := msg.String(), `{ "tag": "Friday" }`; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestUnionMsgStringQuotesStringData(t *testing.T) {
	msg := NewUnionMsg("Name", NewStringMsg(`a:"b,c\d`))
	if got, want := msg.String(), `{ "tag": "Name", "data": "a:\"b,c\\d" }`; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestUnionMsgStringUsesNestedJSONFormatting(t *testing.T) {
	msg := NewUnionMsg("Payload", NewDictMsg(map[string]Msg{
		"text": NewStringMsg("a:b,c"),
		"nums": NewListMsg([]Msg{
			NewIntMsg(1),
			NewIntMsg(2),
		}),
	}))

	if got, want := msg.String(), `{ "tag": "Payload", "data": {"nums": [1, 2], "text": "a:b,c"} }`; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestBytesMsgMarshalJSON(t *testing.T) {
	msg := NewBytesMsg([]byte("hello"))

	b, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
	if got, want := string(b), `"aGVsbG8="`; got != want {
		t.Fatalf("MarshalJSON() = %q, want %q", got, want)
	}
}

func TestBytesMsgEqual(t *testing.T) {
	a := NewBytesMsg([]byte{1, 2, 3})
	b := NewBytesMsg([]byte{1, 2, 3})
	c := NewBytesMsg([]byte{1, 2, 4})

	if !Equal(a, b) {
		t.Fatal("Equal() = false, want true")
	}
	if Equal(a, c) {
		t.Fatal("Equal() = true, want false")
	}
}

func TestEqualContainerRepresentations(t *testing.T) {
	testCases := []struct {
		left  Msg
		right Msg
		name  string
		want  bool
	}{
		{
			name:  "typed and untyped lists with equal integers",
			left:  NewListIntMsg([]int64{1, 2}),
			right: NewListMsg([]Msg{NewIntMsg(1), NewIntMsg(2)}),
			want:  true,
		},
		{
			name:  "typed lists with different scalar kinds",
			left:  NewListIntMsg([]int64{1}),
			right: NewListFloatMsg([]float64{1}),
			want:  false,
		},
		{
			name:  "typed and untyped dictionaries with equal integers",
			left:  NewDictIntMsg(map[string]int64{"one": 1}),
			right: NewDictMsg(map[string]Msg{"one": NewIntMsg(1)}),
			want:  true,
		},
		{
			name:  "dictionaries with different values",
			left:  NewDictStringMsg(map[string]string{"one": "first"}),
			right: NewDictStringMsg(map[string]string{"one": "second"}),
			want:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := Equal(testCase.left, testCase.right); got != testCase.want {
				t.Fatalf("Equal() = %t, want %t", got, testCase.want)
			}
		})
	}
}

func TestEqualStructAndUnion(t *testing.T) {
	leftStruct := NewStructMsg([]StructField{
		NewStructField("value", NewListIntMsg([]int64{1, 2})),
	})
	rightStruct := NewStructMsg([]StructField{
		NewStructField("value", NewListMsg([]Msg{NewIntMsg(1), NewIntMsg(2)})),
	})
	if !Equal(leftStruct, rightStruct) {
		t.Fatal("Equal() = false, want equal structs with equivalent list storage")
	}

	if !Equal(NewUnionMsg("Value", leftStruct), NewUnionMsg("Value", rightStruct)) {
		t.Fatal("Equal() = false, want equal unions with equivalent data")
	}
}

func TestListMsgAccessors(t *testing.T) {
	t.Run("untyped", func(t *testing.T) {
		list := NewListMsg([]Msg{NewIntMsg(1)}).List()
		if got := list.Untyped(); len(got) != 1 || !Equal(got[0], NewIntMsg(1)) {
			t.Fatalf("Untyped() = %v, want one integer message", got)
		}
		mustPanic(t, func() { list.Ints() })
	})

	t.Run("typed", func(t *testing.T) {
		list := NewListIntMsg([]int64{1, 2}).List()
		if got, want := list.Ints(), []int64{1, 2}; !equalInt64s(got, want) {
			t.Fatalf("Ints() = %v, want %v", got, want)
		}
		mustPanic(t, func() { list.Untyped() })
	})
}

func TestDictMsgAccessors(t *testing.T) {
	t.Run("untyped", func(t *testing.T) {
		dict := NewDictMsg(map[string]Msg{"one": NewIntMsg(1)}).Dict()
		if got := dict.Untyped()["one"]; !Equal(got, NewIntMsg(1)) {
			t.Fatalf("Untyped()[one] = %v, want integer message", got)
		}
		mustPanic(t, func() { dict.Ints() })
	})

	t.Run("typed", func(t *testing.T) {
		dict := NewDictIntMsg(map[string]int64{"one": 1}).Dict()
		if got, want := dict.Ints()["one"], int64(1); got != want {
			t.Fatalf("Ints()[one] = %d, want %d", got, want)
		}
		mustPanic(t, func() { dict.Untyped() })
	})
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("function did not panic")
		}
	}()
	fn()
}

func equalInt64s(left, right []int64) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
