package messages

import "strconv"

// Bool

type BoolMsg struct {
	internalMsg
	v bool
}

func (msg BoolMsg) Bool() bool                   { return msg.v }
func (msg BoolMsg) String() string               { return strconv.FormatBool(msg.v) }
func (msg BoolMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		internalMsg: internalMsg{},
		v:           b,
	}
}

func equalBool(left BoolMsg, right Msg) bool {
	return equalBoolValue(left.v, right)
}

func equalBoolValue(left bool, right Msg) bool {
	rightTyped, ok := right.(BoolMsg)
	return ok && left == rightTyped.v
}
