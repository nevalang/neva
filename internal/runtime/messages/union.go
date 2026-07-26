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

func NewUnionMsg(tag string, data Msg) UnionMsg {
	return UnionMsg{
		internalMsg: internalMsg{},
		tag:         tag,
		data:        data,
	}
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
