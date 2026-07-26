package messages

import (
	"encoding/json"
	"fmt"
)

type UnionMsg struct {
	internalMsg
	data Msg
	tag  string
}

func (msg UnionMsg) Union() UnionMsg {
	return msg
}

func (msg UnionMsg) Tag() string {
	return msg.tag
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg UnionMsg) Data() Msg {
	return msg.data
}

func (msg UnionMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg UnionMsg) MarshalJSON() ([]byte, error) {
	if msg.data == nil {
		return fmt.Appendf(nil, `{ "tag": %q }`, msg.tag), nil
	}

	dataJSON, err := json.Marshal(msg.data)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}
	dataJSON = addJSONSpaces(dataJSON)

	return fmt.Appendf(nil, `{ "tag": %q, "data": %s }`, msg.tag, dataJSON), nil
}

func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := msg.(UnionMsg)
	return unionMsg, ok
}

// Uint8Index validates idx and returns it as uint8 or panics.
func Uint8Index(idx int) uint8 {
	if idx < 0 {
		panic(fmt.Sprintf("runtime: negative index %d", idx))
	}
	if idx > int(^uint8(0)) {
		panic(fmt.Sprintf("runtime: index %d overflows uint8", idx))
	}
	// #nosec G115 -- bounds checked above
	return uint8(idx)
}

func NewUnionMsg(tag string, data Msg) UnionMsg {
	return UnionMsg{
		internalMsg: internalMsg{},
		tag:         tag,
		data:        data,
	}
}

// Match compares two messages and return true if they matches and false otherwise.
// Unlike Equal it compares only some aspects of the messages.
func Match(msg Msg, pattern Msg) bool {
	// at the moment we only match unions
	// maybe in the future we'll add support for more types e.g. structs
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	msgUnion, ok := msg.(UnionMsg)
	if !ok {
		return Equal(msg, pattern)
	}

	// both msg and pattern must be unions to perform pattern matching
	// if at least one of them is not, strict equality will be applied instead
	patternUnion, ok := pattern.(UnionMsg)
	if !ok {
		return Equal(msg, pattern)
	}

	// if tags are not equal data does not matter, there's no match
	if msgUnion.tag != patternUnion.tag {
		return false
	}

	// if pattern doesn't have data we match by tag
	// and by this time we know tags are equal
	if patternUnion.data == nil {
		return true
	}

	// if we here we know that pattern has data
	// so if msg doesn't, they don't match
	if msgUnion.data == nil {
		return false
	}

	// by this time we know
	// both msg and pattern are union messages
	// they both have the same tags and some data inside
	// so we apply strict equality to the data they wrap
	// maybe in the future we'll consider recursive matching, we'll see
	return Equal(msgUnion.data, patternUnion.data)
}

func equalUnionValue(left UnionMsg, right Msg) bool {
	rightTyped, ok := right.(UnionMsg)
	return ok && equalUnions(left, rightTyped)
}

func equalUnions(left, right UnionMsg) bool {
	if left.tag != right.tag || (left.data == nil) != (right.data == nil) {
		return false
	}
	return left.data == nil || Equal(left.data, right.data)
}
