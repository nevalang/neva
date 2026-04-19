package runtime

import (
	"testing"
)

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

	if !a.Equal(b) {
		t.Fatal("Equal() = false, want true")
	}
	if a.Equal(c) {
		t.Fatal("Equal() = true, want false")
	}
}
