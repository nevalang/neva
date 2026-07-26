package runtime

import "testing"

func TestOrderedMessageValueCompatibility(t *testing.T) {
	orderedInt := OrderedMsg{Msg: NewIntMsg(1)}

	if !Equal(orderedInt, NewIntMsg(1)) {
		t.Fatal("Equal() = false, want left ordered payload to compare by value")
	}
	if Equal(NewIntMsg(1), orderedInt) {
		t.Fatal("Equal() = true, want right ordered envelope to remain distinct")
	}

	union := NewUnionMsg("tag", NewIntMsg(1))
	if got, ok := AsUnion(OrderedMsg{Msg: union}); !ok || !Equal(got, union) {
		t.Fatalf("AsUnion() = (%v, %t), want ordered union payload", got, ok)
	}

	if Match(OrderedMsg{Msg: union}, NewUnionMsg("tag", nil)) {
		t.Fatal("Match() = true, want ordered envelope to retain strict-equality fallback")
	}
}
