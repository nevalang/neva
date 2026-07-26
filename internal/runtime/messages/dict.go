package messages

import "encoding/json"

// DictMsg provides access to the storage of a dictionary runtime message.
//
// Exactly one value accessor is valid for an implementation. Untyped exposes
// the boxed representation; scalar accessors expose the corresponding typed
// storage. Calling any other accessor is a runtime invariant violation.
type DictMsg interface {
	Msg
	Untyped() map[string]Msg
	Bools() map[string]bool
	Ints() map[string]int64
	Floats() map[string]float64
	Strings() map[string]string
}

// internalDictMsg supplies invariant-violation methods shared by every dict
// value. Concrete implementations override Dict, their one valid accessor,
// String, and MarshalJSON.
type internalDictMsg struct{ internalMsg }

func (internalDictMsg) Untyped() map[string]Msg {
	panic("unexpected Untyped method call on typed dict message")
}
func (internalDictMsg) Bools() map[string]bool { panic("unexpected Bools method call on dict message") }
func (internalDictMsg) Ints() map[string]int64 { panic("unexpected Ints method call on dict message") }
func (internalDictMsg) Floats() map[string]float64 {
	panic("unexpected Floats method call on dict message")
}
func (internalDictMsg) Strings() map[string]string {
	panic("unexpected Strings method call on dict message")
}
func (internalDictMsg) String() string {
	panic("unexpected String method call on internal dict message")
}
func (internalDictMsg) MarshalJSON() ([]byte, error) {
	panic("unexpected MarshalJSON method call on internal dict message")
}

// untypedDictMsg stores a dictionary whose values are already boxed messages.
type untypedDictMsg struct {
	internalDictMsg
	v map[string]Msg
}

func (msg untypedDictMsg) Untyped() map[string]Msg { return msg.v }

//nolint:ireturn // DictMsg is the runtime dictionary contract.
func (msg untypedDictMsg) Dict() DictMsg  { return msg }
func (msg untypedDictMsg) String() string { return mustJSON(msg) }
func (msg untypedDictMsg) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(msg.v)
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	}
	return addJSONSpaces(jsonData), nil
}

// boolDictMsg stores unboxed boolean dictionary values.
type boolDictMsg struct {
	internalDictMsg
	v map[string]bool
}

func (msg boolDictMsg) Bools() map[string]bool { return msg.v }

//nolint:ireturn // DictMsg is the runtime dictionary contract.
func (msg boolDictMsg) Dict() DictMsg  { return msg }
func (msg boolDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg boolDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// intDictMsg stores unboxed integer dictionary values.
type intDictMsg struct {
	internalDictMsg
	v map[string]int64
}

func (msg intDictMsg) Ints() map[string]int64 { return msg.v }

//nolint:ireturn // DictMsg is the runtime dictionary contract.
func (msg intDictMsg) Dict() DictMsg  { return msg }
func (msg intDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg intDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// floatDictMsg stores unboxed float dictionary values.
type floatDictMsg struct {
	internalDictMsg
	v map[string]float64
}

func (msg floatDictMsg) Floats() map[string]float64 { return msg.v }

//nolint:ireturn // DictMsg is the runtime dictionary contract.
func (msg floatDictMsg) Dict() DictMsg  { return msg }
func (msg floatDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg floatDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// stringDictMsg stores unboxed string dictionary values.
type stringDictMsg struct {
	internalDictMsg
	v map[string]string
}

func (msg stringDictMsg) Strings() map[string]string { return msg.v }

//nolint:ireturn // DictMsg is the runtime dictionary contract.
func (msg stringDictMsg) Dict() DictMsg  { return msg }
func (msg stringDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg stringDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// NewDictMsg creates a dictionary with an untyped boxed representation.
//
//nolint:ireturn // Msg contract type.
func NewDictMsg(d map[string]Msg) Msg {
	return untypedDictMsg{v: d}
}

// NewDictBoolMsg creates a dictionary with unboxed boolean storage.
//
//nolint:ireturn // Msg contract type.
func NewDictBoolMsg(d map[string]bool) Msg {
	return boolDictMsg{v: d}
}

// NewDictIntMsg creates a dictionary with unboxed integer storage.
//
//nolint:ireturn // Msg contract type.
func NewDictIntMsg(d map[string]int64) Msg {
	return intDictMsg{v: d}
}

// NewDictFloatMsg creates a dictionary with unboxed float storage.
//
//nolint:ireturn // Msg contract type.
func NewDictFloatMsg(d map[string]float64) Msg {
	return floatDictMsg{v: d}
}

// NewDictStringMsg creates a dictionary with unboxed string storage.
//
//nolint:ireturn // Msg contract type.
func NewDictStringMsg(d map[string]string) Msg {
	return stringDictMsg{v: d}
}

func dictEqual(left DictMsg, right DictMsg) bool {
	if DictLen(left) != DictLen(right) {
		return false
	}

	leftMsgs := DictToMessageMap(left)
	rightMsgs := DictToMessageMap(right)
	for key, leftVal := range leftMsgs {
		rightVal, ok := rightMsgs[key]
		if !ok || !Equal(leftVal, rightVal) {
			return false
		}
	}
	return true
}

// DictAsUntyped returns the boxed dictionary representation when present.
func DictAsUntyped(dict DictMsg) (map[string]Msg, bool) {
	typed, ok := dict.(untypedDictMsg)
	return typed.v, ok
}

// DictAsBools returns the boolean dictionary representation when present.
func DictAsBools(dict DictMsg) (map[string]bool, bool) {
	typed, ok := dict.(boolDictMsg)
	return typed.v, ok
}

// DictAsInts returns the integer dictionary representation when present.
func DictAsInts(dict DictMsg) (map[string]int64, bool) {
	typed, ok := dict.(intDictMsg)
	return typed.v, ok
}

// DictAsFloats returns the float dictionary representation when present.
func DictAsFloats(dict DictMsg) (map[string]float64, bool) {
	typed, ok := dict.(floatDictMsg)
	return typed.v, ok
}

// DictAsStrings returns the string dictionary representation when present.
func DictAsStrings(dict DictMsg) (map[string]string, bool) {
	typed, ok := dict.(stringDictMsg)
	return typed.v, ok
}

// DictLen returns the number of entries in dict without boxing its storage.
func DictLen(dict DictMsg) int {
	switch typed := dict.(type) {
	case untypedDictMsg:
		return len(typed.v)
	case boolDictMsg:
		return len(typed.v)
	case intDictMsg:
		return len(typed.v)
	case floatDictMsg:
		return len(typed.v)
	case stringDictMsg:
		return len(typed.v)
	default:
		panic("unexpected dict implementation")
	}
}

// DictToMessageMap returns the values of dict as runtime messages.
//
// For an untyped dictionary, it returns the existing message map. For a typed
// scalar dictionary, it allocates an untyped message map and converts every
// value.
func DictToMessageMap(dict DictMsg) map[string]Msg {
	switch typed := dict.(type) {
	case intDictMsg:
		values := typed.v
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewIntMsg(value)
		}
		return msgs
	case stringDictMsg:
		values := typed.v
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewStringMsg(value)
		}
		return msgs
	case boolDictMsg:
		values := typed.v
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewBoolMsg(value)
		}
		return msgs
	case floatDictMsg:
		values := typed.v
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewFloatMsg(value)
		}
		return msgs
	case untypedDictMsg:
		return typed.v
	default:
		panic("unexpected dict implementation")
	}
}

// DictGetValueByKey returns the value stored for key without materializing typed storage.
// The boolean is false when key is not present.
//
//nolint:ireturn // Msg is the value-layer contract.
func DictGetValueByKey(dict DictMsg, key string) (Msg, bool) {
	switch typed := dict.(type) {
	case intDictMsg:
		values := typed.v
		value, found := values[key]
		return NewIntMsg(value), found
	case stringDictMsg:
		values := typed.v
		value, found := values[key]
		return NewStringMsg(value), found
	case boolDictMsg:
		values := typed.v
		value, found := values[key]
		return NewBoolMsg(value), found
	case floatDictMsg:
		values := typed.v
		value, found := values[key]
		return NewFloatMsg(value), found
	case untypedDictMsg:
		value, found := typed.v[key]
		return value, found
	default:
		panic("unexpected dict implementation")
	}
}

func equalDictValue(left DictMsg, right Msg) bool {
	rightDict, ok := right.(DictMsg)
	return ok && dictEqual(left, rightDict)
}
