package runtime

import "testing"

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic: %s", name)
		}
	}()

	fn()
}

func TestStructMsgEqualIgnoresFieldOrder(t *testing.T) {
	left := NewStructMsg([]StructField{
		NewStructField("a", NewIntMsg(1)),
		NewStructField("b", NewStringMsg("x")),
	})
	right := NewStructMsg([]StructField{
		NewStructField("b", NewStringMsg("x")),
		NewStructField("a", NewIntMsg(1)),
	})

	if !left.Equal(right) {
		t.Fatalf("expected struct messages to be equal")
	}
}

func TestMatchUnion(t *testing.T) {
	msgWithData := NewUnionMsg("ok", NewIntMsg(42))
	patternWithData := NewUnionMsg("ok", NewIntMsg(42))
	patternNoData := NewUnionMsgNoData("ok")
	wrongTag := NewUnionMsgNoData("err")

	if !Match(msgWithData, patternWithData) {
		t.Fatalf("expected union with matching tag/data to match")
	}
	if !Match(msgWithData, patternNoData) {
		t.Fatalf("expected union pattern without data to match by tag")
	}
	if Match(msgWithData, wrongTag) {
		t.Fatalf("expected union tags to be required for match")
	}
}

func TestDictMarshalJSONSpacing(t *testing.T) {
	msg := NewDictMsg(map[string]Msg{
		"a": NewIntMsg(1),
	})
	data, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if string(data) != `{"a": 1}` {
		t.Fatalf("unexpected dict json: %s", data)
	}
}

func TestUnionMarshalJSON(t *testing.T) {
	msg := NewUnionMsgNoData("ok")
	data, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	if string(data) != `{ "tag": "ok" }` {
		t.Fatalf("unexpected union json: %s", data)
	}
}

func TestIntMsgNegativeRoundTrip(t *testing.T) {
	msg := NewIntMsg(-42)
	if got := msg.Int(); got != -42 {
		t.Fatalf("expected -42, got %d", got)
	}
}

func TestInvalidMsgPanics(t *testing.T) {
	invalid := Msg{}

	mustPanic(t, "String", func() {
		_ = invalid.String()
	})
	mustPanic(t, "MarshalJSON", func() {
		_, err := invalid.MarshalJSON()
		if err != nil {
			t.Fatalf("unexpected error before panic: %v", err)
		}
	})
	mustPanic(t, "Equal", func() {
		_ = invalid.Equal(Msg{})
	})
}

func TestDictMsgMarshalJSONPreservesStringValues(t *testing.T) {
	msg := NewDictMsg(map[string]Msg{
		"text": NewStringMsg(`a:"b,c\d`),
		"nums": NewListMsg([]Msg{
			NewIntMsg(1),
			NewIntMsg(2),
		}),
	})

	b, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}
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
	msg := NewUnionMsgNoData("Friday")
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

	if !a.Equal(b) {
		t.Fatal("Equal() = false, want true")
	}
	if a.Equal(c) {
		t.Fatal("Equal() = true, want false")
	}
}
