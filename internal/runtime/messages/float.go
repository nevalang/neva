package messages

import "fmt"

// Float

type FloatMsg struct {
	internalMsg
	v float64
}

func (msg FloatMsg) Float() float64               { return msg.v }
func (msg FloatMsg) String() string               { return fmt.Sprint(msg.v) }
func (msg FloatMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		internalMsg: internalMsg{},
		v:           n,
	}
}

func equalFloat(left FloatMsg, right Msg) bool {
	return equalFloatValue(left.v, right)
}

func equalFloatValue(left float64, right Msg) bool {
	rightTyped, ok := right.(FloatMsg)
	return ok && left == rightTyped.v
}

// FloatAdd returns the sum of left and right.
func FloatAdd(left, right Msg) FloatMsg { return NewFloatMsg(left.Float() + right.Float()) }

// FloatSubtract returns left minus right.
func FloatSubtract(left, right Msg) FloatMsg { return NewFloatMsg(left.Float() - right.Float()) }

// FloatMultiply returns the product of left and right.
func FloatMultiply(left, right Msg) FloatMsg { return NewFloatMsg(left.Float() * right.Float()) }

// FloatDivide returns left divided by right.
func FloatDivide(left, right Msg) FloatMsg { return NewFloatMsg(left.Float() / right.Float()) }

// FloatNegate returns the arithmetic negation of value.
func FloatNegate(value Msg) FloatMsg { return NewFloatMsg(-value.Float()) }

// FloatIsGreater reports whether left is greater than right.
func FloatIsGreater(left, right Msg) BoolMsg { return NewBoolMsg(left.Float() > right.Float()) }

// FloatIsLesser reports whether left is less than right.
func FloatIsLesser(left, right Msg) BoolMsg { return NewBoolMsg(left.Float() < right.Float()) }

// FloatFromInt converts value to a float.
func FloatFromInt(value Msg) FloatMsg { return NewFloatMsg(float64(value.Int())) }
