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

// BoolAnd returns the logical conjunction of left and right.
func BoolAnd(left, right Msg) BoolMsg { return NewBoolMsg(left.Bool() && right.Bool()) }

// BoolOr returns the logical disjunction of left and right.
func BoolOr(left, right Msg) BoolMsg { return NewBoolMsg(left.Bool() || right.Bool()) }

// BoolNot returns the logical negation of value.
func BoolNot(value Msg) BoolMsg { return NewBoolMsg(!value.Bool()) }
