package messages

import "strconv"

// Int

type IntMsg struct {
	internalMsg
	v int64
}

func (msg IntMsg) Int() int64                   { return msg.v }
func (msg IntMsg) String() string               { return strconv.Itoa(int(msg.v)) }
func (msg IntMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		internalMsg: internalMsg{},
		v:           n,
	}
}

func equalInt(left IntMsg, right Msg) bool {
	rightTyped, ok := right.(IntMsg)
	return ok && left.v == rightTyped.v
}
