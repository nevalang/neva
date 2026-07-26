package runtime

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// OrderedMsg is a transport envelope with payload and runtime ordering metadata.
type OrderedMsg struct {
	Msg
	index uint64
}

// String is just a simple stringer that ignores index while formatting.
func (o OrderedMsg) String() string { return fmt.Sprint(o.Msg) }

// Equal reports whether two runtime messages have the same value.
//
// Ordering metadata is ignored when the left-hand value is an OrderedMsg.
func Equal(left, right Msg) bool {
	if ordered, ok := left.(OrderedMsg); ok {
		return Equal(ordered.Msg, right)
	}
	if _, ok := right.(OrderedMsg); ok {
		return false
	}
	return messages.Equal(left, right)
}

// AsUnion returns the union payload of a message, ignoring transport ordering metadata.
func AsUnion(msg Msg) (UnionMsg, bool) {
	if ordered, ok := msg.(OrderedMsg); ok {
		return messages.AsUnion(ordered.Msg)
	}
	return messages.AsUnion(msg)
}

// Match compares two runtime messages while preserving transport compatibility.
func Match(msg, pattern Msg) bool {
	if _, ok := msg.(OrderedMsg); ok {
		return Equal(msg, pattern)
	}
	return messages.Match(msg, pattern)
}
