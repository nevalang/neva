package messages

import "encoding/json"

// ListMsg provides access to the storage of a list runtime message.
//
// Exactly one value accessor is valid for an implementation. Untyped exposes
// the boxed representation; scalar accessors expose the corresponding typed
// storage. Calling any other accessor is a runtime invariant violation.
type ListMsg interface {
	Msg
	Untyped() []Msg
	Bools() []bool
	Ints() []int64
	Floats() []float64
	Strings() []string
}

// internalListMsg supplies invariant-violation methods shared by every list
// value. Concrete implementations override List, their one valid accessor,
// String, and MarshalJSON.
type internalListMsg struct{ internalMsg }

func (internalListMsg) Untyped() []Msg    { panic("unexpected Untyped method call on typed list message") }
func (internalListMsg) Bools() []bool     { panic("unexpected Bools method call on list message") }
func (internalListMsg) Ints() []int64     { panic("unexpected Ints method call on list message") }
func (internalListMsg) Floats() []float64 { panic("unexpected Floats method call on list message") }
func (internalListMsg) Strings() []string { panic("unexpected Strings method call on list message") }
func (internalListMsg) String() string {
	panic("unexpected String method call on internal list message")
}
func (internalListMsg) MarshalJSON() ([]byte, error) {
	panic("unexpected MarshalJSON method call on internal list message")
}

// untypedListMsg stores a list whose elements are already boxed messages.
type untypedListMsg struct {
	internalListMsg
	v []Msg
}

func (msg untypedListMsg) Untyped() []Msg { return msg.v }

//nolint:ireturn // ListMsg is the runtime list contract.
func (msg untypedListMsg) List() ListMsg { return msg }
func (msg untypedListMsg) String() string {
	bb, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bb)
}

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg untypedListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// boolListMsg stores unboxed boolean list elements.
type boolListMsg struct {
	internalListMsg
	v []bool
}

func (msg boolListMsg) Bools() []bool { return msg.v }

//nolint:ireturn // ListMsg is the runtime list contract.
func (msg boolListMsg) List() ListMsg  { return msg }
func (msg boolListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg boolListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// intListMsg stores unboxed integer list elements.
type intListMsg struct {
	internalListMsg
	v []int64
}

func (msg intListMsg) Ints() []int64 { return msg.v }

//nolint:ireturn // ListMsg is the runtime list contract.
func (msg intListMsg) List() ListMsg  { return msg }
func (msg intListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg intListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// floatListMsg stores unboxed float list elements.
type floatListMsg struct {
	internalListMsg
	v []float64
}

func (msg floatListMsg) Floats() []float64 {
	return msg.v
}

//nolint:ireturn // ListMsg is the runtime list contract.
func (msg floatListMsg) List() ListMsg  { return msg }
func (msg floatListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg floatListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// stringListMsg stores unboxed string list elements.
type stringListMsg struct {
	internalListMsg
	v []string
}

func (msg stringListMsg) Strings() []string { return msg.v }

//nolint:ireturn // ListMsg is the runtime list contract.
func (msg stringListMsg) List() ListMsg  { return msg }
func (msg stringListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg stringListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}

// NewListMsg creates a list with an untyped boxed representation.
//
//nolint:ireturn // Msg contract type.
func NewListMsg(v []Msg) Msg {
	return untypedListMsg{v: v}
}

// NewListBoolMsg creates a list with unboxed boolean storage.
//
//nolint:ireturn // Msg contract type.
func NewListBoolMsg(v []bool) Msg {
	return boolListMsg{v: v}
}

// NewListIntMsg creates a list with unboxed integer storage.
//
//nolint:ireturn // Msg contract type.
func NewListIntMsg(v []int64) Msg {
	return intListMsg{v: v}
}

// NewListFloatMsg creates a list with unboxed float storage.
//
//nolint:ireturn // Msg contract type.
func NewListFloatMsg(v []float64) Msg {
	return floatListMsg{v: v}
}

// NewListStringMsg creates a list with unboxed string storage.
//
//nolint:ireturn // Msg contract type.
func NewListStringMsg(v []string) Msg {
	return stringListMsg{v: v}
}

func listEqualUntyped(left []Msg, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case untypedListMsg:
		return listEqualUntypedToUntyped(left, rightTyped.v)
	case boolListMsg:
		return listEqualUntypedToBools(left, rightTyped.v)
	case intListMsg:
		return listEqualUntypedToInts(left, rightTyped.v)
	case floatListMsg:
		return listEqualUntypedToFloats(left, rightTyped.v)
	case stringListMsg:
		return listEqualUntypedToStrings(left, rightTyped.v)
	default:
		panic("unexpected list implementation")
	}
}

func listEqualUntypedToUntyped(left []Msg, right []Msg) bool {
	for i := range left {
		if !Equal(left[i], right[i]) {
			return false
		}
	}
	return true
}

func listEqualUntypedToBools(left []Msg, right []bool) bool {
	for i := range left {
		if !Equal(left[i], NewBoolMsg(right[i])) {
			return false
		}
	}
	return true
}

func listEqualUntypedToInts(left []Msg, right []int64) bool {
	for i := range left {
		if !Equal(left[i], NewIntMsg(right[i])) {
			return false
		}
	}
	return true
}

func listEqualUntypedToFloats(left []Msg, right []float64) bool {
	for i := range left {
		if !Equal(left[i], NewFloatMsg(right[i])) {
			return false
		}
	}
	return true
}

func listEqualUntypedToStrings(left []Msg, right []string) bool {
	for i := range left {
		if !Equal(left[i], NewStringMsg(right[i])) {
			return false
		}
	}
	return true
}

func listEqualBool(left []bool, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case boolListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case untypedListMsg:
		for i := range left {
			if !Equal(NewBoolMsg(left[i]), rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualInt(left []int64, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case intListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case untypedListMsg:
		for i := range left {
			if !Equal(NewIntMsg(left[i]), rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualFloat(left []float64, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case floatListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case untypedListMsg:
		for i := range left {
			if !Equal(NewFloatMsg(left[i]), rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualString(left []string, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case stringListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case untypedListMsg:
		for i := range left {
			if !Equal(NewStringMsg(left[i]), rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

// ListAsUntyped returns the boxed list representation when present.
func ListAsUntyped(list ListMsg) ([]Msg, bool) {
	typed, ok := list.(untypedListMsg)
	return typed.v, ok
}

// ListAsBools returns the boolean list representation when present.
func ListAsBools(list ListMsg) ([]bool, bool) {
	typed, ok := list.(boolListMsg)
	return typed.v, ok
}

// ListAsInts returns the integer list representation when present.
func ListAsInts(list ListMsg) ([]int64, bool) {
	typed, ok := list.(intListMsg)
	return typed.v, ok
}

// ListAsFloats returns the float list representation when present.
func ListAsFloats(list ListMsg) ([]float64, bool) {
	typed, ok := list.(floatListMsg)
	return typed.v, ok
}

// ListAsStrings returns the string list representation when present.
func ListAsStrings(list ListMsg) ([]string, bool) {
	typed, ok := list.(stringListMsg)
	return typed.v, ok
}

// ListLen returns the number of elements in list without boxing its storage.
func ListLen(list ListMsg) int {
	switch typed := list.(type) {
	case untypedListMsg:
		return len(typed.v)
	case boolListMsg:
		return len(typed.v)
	case intListMsg:
		return len(typed.v)
	case floatListMsg:
		return len(typed.v)
	case stringListMsg:
		return len(typed.v)
	default:
		panic("unexpected list implementation")
	}
}

// ListToMessageSlice returns the elements of list as runtime messages.
//
// For an untyped list, it returns the existing message slice. For a typed
// scalar list, it allocates an untyped message slice and converts every
// element.
func ListToMessageSlice(list ListMsg) []Msg {
	switch typed := list.(type) {
	case intListMsg:
		values := typed.v
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewIntMsg(values[i])
		}
		return msgs
	case stringListMsg:
		values := typed.v
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewStringMsg(values[i])
		}
		return msgs
	case boolListMsg:
		values := typed.v
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewBoolMsg(values[i])
		}
		return msgs
	case floatListMsg:
		values := typed.v
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewFloatMsg(values[i])
		}
		return msgs
	case untypedListMsg:
		return typed.v
	default:
		panic("unexpected list implementation")
	}
}

func equalListValue(left ListMsg, right Msg) bool {
	rightList, ok := right.(ListMsg)
	return ok && equalLists(left, rightList)
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
