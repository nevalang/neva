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
	return equalIntValue(left.v, right)
}

func equalIntValue(left int64, right Msg) bool {
	rightTyped, ok := right.(IntMsg)
	return ok && left == rightTyped.v
}

// IntAdd returns the sum of left and right.
func IntAdd(left, right Msg) IntMsg { return NewIntMsg(left.Int() + right.Int()) }

// IntSubtract returns left minus right.
func IntSubtract(left, right Msg) IntMsg { return NewIntMsg(left.Int() - right.Int()) }

// IntMultiply returns the product of left and right.
func IntMultiply(left, right Msg) IntMsg { return NewIntMsg(left.Int() * right.Int()) }

// IntDivide returns left divided by right.
func IntDivide(left, right Msg) IntMsg { return NewIntMsg(left.Int() / right.Int()) }

// IntModulo returns the remainder of left divided by right.
func IntModulo(left, right Msg) IntMsg { return NewIntMsg(left.Int() % right.Int()) }

// IntNegate returns the arithmetic negation of value.
func IntNegate(value Msg) IntMsg { return NewIntMsg(-value.Int()) }

// IntIncrement returns value plus one.
func IntIncrement(value Msg) IntMsg { return NewIntMsg(value.Int() + 1) }

// IntDecrement returns value minus one.
func IntDecrement(value Msg) IntMsg { return NewIntMsg(value.Int() - 1) }

// IntPower returns base raised to exponent.
func IntPower(base, exponent Msg) IntMsg {
	return NewIntMsg(intPower(base.Int(), exponent.Int()))
}

func intPower(base, exponent int64) int64 {
	if exponent < 0 {
		return 1
	}

	result := int64(1)
	for exponent > 0 {
		if exponent&1 == 1 {
			result *= base
		}
		base *= base
		exponent >>= 1
	}
	return result
}

// IntBitwiseAnd returns the bitwise conjunction of left and right.
func IntBitwiseAnd(left, right Msg) IntMsg { return NewIntMsg(left.Int() & right.Int()) }

// IntBitwiseOr returns the bitwise disjunction of left and right.
func IntBitwiseOr(left, right Msg) IntMsg { return NewIntMsg(left.Int() | right.Int()) }

// IntBitwiseXor returns the bitwise exclusive disjunction of left and right.
func IntBitwiseXor(left, right Msg) IntMsg { return NewIntMsg(left.Int() ^ right.Int()) }

// IntShiftLeft returns left shifted left by right bits.
func IntShiftLeft(left, right Msg) IntMsg { return NewIntMsg(left.Int() << right.Int()) }

// IntShiftRight returns left shifted right by right bits.
func IntShiftRight(left, right Msg) IntMsg { return NewIntMsg(left.Int() >> right.Int()) }

// IntIsGreater reports whether left is greater than right.
func IntIsGreater(left, right Msg) BoolMsg { return NewBoolMsg(left.Int() > right.Int()) }

// IntIsGreaterOrEqual reports whether left is greater than or equal to right.
func IntIsGreaterOrEqual(left, right Msg) BoolMsg { return NewBoolMsg(left.Int() >= right.Int()) }

// IntIsLesser reports whether left is less than right.
func IntIsLesser(left, right Msg) BoolMsg { return NewBoolMsg(left.Int() < right.Int()) }

// IntIsLesserOrEqual reports whether left is less than or equal to right.
func IntIsLesserOrEqual(left, right Msg) BoolMsg { return NewBoolMsg(left.Int() <= right.Int()) }

// IntFromFloat converts value to an integer.
func IntFromFloat(value Msg) IntMsg { return NewIntMsg(int64(value.Float())) }
