package messages

// Equal reports whether two runtime messages have the same value.
//
// Equality is defined across typed and untyped list/dict storage.
func Equal(left, right Msg) bool {
	switch leftTyped := left.(type) {
	case BoolMsg:
		return equalBool(leftTyped, right)
	case IntMsg:
		return equalInt(leftTyped, right)
	case FloatMsg:
		return equalFloat(leftTyped, right)
	case StringMsg:
		return equalString(leftTyped, right)
	case BytesMsg:
		return equalBytes(leftTyped, right)
	case ListMsg:
		return equalListValue(leftTyped, right)
	case DictMsg:
		return equalDictValue(leftTyped, right)
	case StructMsg:
		return equalStructValue(leftTyped, right)
	case UnionMsg:
		return equalUnionValue(leftTyped, right)
	default:
		panic("unexpected runtime message implementation")
	}
}
