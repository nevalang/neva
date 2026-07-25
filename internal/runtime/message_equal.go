package runtime

import "bytes"

// Equal reports whether two runtime messages have the same value.
//
// Equality is defined across typed and untyped list/dict storage. Transport
// ordering metadata is ignored when the left-hand value is an OrderedMsg.
func Equal(left, right Msg) bool {
	switch leftTyped := left.(type) {
	case OrderedMsg:
		return Equal(leftTyped.Msg, right)
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
	case listValueMsg:
		return equalListValues(leftTyped, right)
	case dictValueMsg:
		return equalDictValues(leftTyped, right)
	case StructMsg:
		return equalStructValue(leftTyped, right)
	case UnionMsg:
		return equalUnionValue(leftTyped, right)
	default:
		panic("unexpected runtime message implementation")
	}
}

func equalBool(left BoolMsg, right Msg) bool {
	rightTyped, ok := right.(BoolMsg)
	return ok && left.v == rightTyped.v
}

func equalInt(left IntMsg, right Msg) bool {
	rightTyped, ok := right.(IntMsg)
	return ok && left.v == rightTyped.v
}

func equalFloat(left FloatMsg, right Msg) bool {
	rightTyped, ok := right.(FloatMsg)
	return ok && left.v == rightTyped.v
}

func equalString(left StringMsg, right Msg) bool {
	rightTyped, ok := right.(StringMsg)
	return ok && left.v == rightTyped.v
}

func equalBytes(left BytesMsg, right Msg) bool {
	rightTyped, ok := right.(BytesMsg)
	return ok && bytes.Equal(left.v, rightTyped.v)
}

func equalListValues(left listValueMsg, right Msg) bool {
	rightTyped, ok := right.(listValueMsg)
	return ok && equalLists(left.v, rightTyped.v)
}

func equalDictValues(left dictValueMsg, right Msg) bool {
	rightTyped, ok := right.(dictValueMsg)
	return ok && dictEqual(left.v, rightTyped.v)
}

func equalStructValue(left StructMsg, right Msg) bool {
	rightTyped, ok := right.(StructMsg)
	return ok && equalStructs(left, rightTyped)
}

func equalUnionValue(left UnionMsg, right Msg) bool {
	rightTyped, ok := right.(UnionMsg)
	return ok && equalUnions(left, rightTyped)
}

func equalLists(left, right ListMsg) bool {
	switch leftTyped := left.(type) {
	case untypedListMsg:
		return listEqualUntyped(leftTyped.v, right)
	case boolListMsg:
		return listEqualBool(leftTyped.v, right)
	case intListMsg:
		return listEqualInt(leftTyped.v, right)
	case floatListMsg:
		return listEqualFloat(leftTyped.v, right)
	case stringListMsg:
		return listEqualString(leftTyped.v, right)
	default:
		panic("unexpected list implementation")
	}
}

func equalStructs(left, right StructMsg) bool {
	if len(left.fields) != len(right.fields) {
		return false
	}
	for i := range left.fields {
		rightValue, ok := right.get(left.fields[i].name)
		if !ok || !Equal(left.fields[i].value, rightValue) {
			return false
		}
	}
	return true
}

func equalUnions(left, right UnionMsg) bool {
	if left.tag != right.tag || (left.data == nil) != (right.data == nil) {
		return false
	}
	return left.data == nil || Equal(left.data, right.data)
}
