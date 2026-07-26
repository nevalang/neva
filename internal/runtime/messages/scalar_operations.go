package messages

// BoolAnd returns the logical conjunction of left and right.
func BoolAnd(left, right Msg) BoolMsg { return NewBoolMsg(left.Bool() && right.Bool()) }

// BoolOr returns the logical disjunction of left and right.
func BoolOr(left, right Msg) BoolMsg { return NewBoolMsg(left.Bool() || right.Bool()) }

// BoolNot returns the logical negation of value.
func BoolNot(value Msg) BoolMsg { return NewBoolMsg(!value.Bool()) }

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
	result := int64(1)
	for range exponent.Int() {
		result *= base.Int()
	}
	return NewIntMsg(result)
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

// StringConcat returns the concatenation of left and right.
func StringConcat(left, right Msg) StringMsg { return NewStringMsg(left.Str() + right.Str()) }

// StringIsGreater reports whether left is lexicographically greater than right.
func StringIsGreater(left, right Msg) BoolMsg { return NewBoolMsg(left.Str() > right.Str()) }

// StringIsLesser reports whether left is lexicographically less than right.
func StringIsLesser(left, right Msg) BoolMsg { return NewBoolMsg(left.Str() < right.Str()) }

// IntFromFloat converts value to an integer.
func IntFromFloat(value Msg) IntMsg { return NewIntMsg(int64(value.Float())) }

// FloatFromInt converts value to a float.
func FloatFromInt(value Msg) FloatMsg { return NewFloatMsg(float64(value.Int())) }
