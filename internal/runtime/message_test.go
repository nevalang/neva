package runtime

import (
	"testing"
)

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
