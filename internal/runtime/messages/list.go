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

// NewListMsg creates the most specific scalar list representation available.
//
// It keeps a homogeneous bool, int, float, or string message slice unboxed.
// Use a scalar constructor such as NewListIntMsg when native scalar storage is
// already available, or NewUntypedListMsg only when boxed storage is required.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewListMsg(values []Msg) ListMsg {
	if len(values) == 0 {
		return NewUntypedListMsg(values)
	}

	switch values[0].(type) {
	case BoolMsg:
		result, ok := listBoolsFromMessages(values)
		if ok {
			return NewListBoolMsg(result)
		}
	case IntMsg:
		result, ok := listIntsFromMessages(values)
		if ok {
			return NewListIntMsg(result)
		}
	case FloatMsg:
		result, ok := listFloatsFromMessages(values)
		if ok {
			return NewListFloatMsg(result)
		}
	case StringMsg:
		result, ok := listStringsFromMessages(values)
		if ok {
			return NewListStringMsg(result)
		}
	default:
		return NewUntypedListMsg(values)
	}
	return NewUntypedListMsg(values)
}

// NewUntypedListMsg creates a list with explicitly boxed message storage.
//
// Prefer NewListMsg, which preserves scalar storage when the values allow it.
// Use this constructor only when boxed storage is required by the caller.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewUntypedListMsg(values []Msg) ListMsg {
	return untypedListMsg{v: values}
}

func listBoolsFromMessages(values []Msg) ([]bool, bool) {
	result := make([]bool, len(values))
	for i, value := range values {
		scalar, ok := value.(BoolMsg)
		if !ok {
			return nil, false
		}
		result[i] = scalar.v
	}
	return result, true
}

func listIntsFromMessages(values []Msg) ([]int64, bool) {
	result := make([]int64, len(values))
	for i, value := range values {
		scalar, ok := value.(IntMsg)
		if !ok {
			return nil, false
		}
		result[i] = scalar.v
	}
	return result, true
}

func listFloatsFromMessages(values []Msg) ([]float64, bool) {
	result := make([]float64, len(values))
	for i, value := range values {
		scalar, ok := value.(FloatMsg)
		if !ok {
			return nil, false
		}
		result[i] = scalar.v
	}
	return result, true
}

func listStringsFromMessages(values []Msg) ([]string, bool) {
	result := make([]string, len(values))
	for i, value := range values {
		scalar, ok := value.(StringMsg)
		if !ok {
			return nil, false
		}
		result[i] = scalar.v
	}
	return result, true
}

// NewListBoolMsg creates a list with unboxed boolean storage.
//
// Prefer this constructor when boolean storage is already available.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewListBoolMsg(v []bool) ListMsg {
	return boolListMsg{v: v}
}

// NewListIntMsg creates a list with unboxed integer storage.
//
// Prefer this constructor when integer storage is already available.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewListIntMsg(v []int64) ListMsg {
	return intListMsg{v: v}
}

// NewListFloatMsg creates a list with unboxed float storage.
//
// Prefer this constructor when float storage is already available.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewListFloatMsg(v []float64) ListMsg {
	return floatListMsg{v: v}
}

// NewListStringMsg creates a list with unboxed string storage.
//
// Prefer this constructor when string storage is already available.
//
//nolint:ireturn // ListMsg is the runtime list contract.
func NewListStringMsg(v []string) ListMsg {
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
		if !equalBoolValue(right[i], left[i]) {
			return false
		}
	}
	return true
}

func listEqualUntypedToInts(left []Msg, right []int64) bool {
	for i := range left {
		if !equalIntValue(right[i], left[i]) {
			return false
		}
	}
	return true
}

func listEqualUntypedToFloats(left []Msg, right []float64) bool {
	for i := range left {
		if !equalFloatValue(right[i], left[i]) {
			return false
		}
	}
	return true
}

func listEqualUntypedToStrings(left []Msg, right []string) bool {
	for i := range left {
		if !equalStringValue(right[i], left[i]) {
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
			if !equalBoolValue(left[i], rightTyped.v[i]) {
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
			if !equalIntValue(left[i], rightTyped.v[i]) {
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
			if !equalFloatValue(left[i], rightTyped.v[i]) {
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
			if !equalStringValue(left[i], rightTyped.v[i]) {
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

// ListAt returns the item at index. Negative indexes count from the end.
// The boolean is false when index is outside the list bounds.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListAt(list ListMsg, index int64) (Msg, bool) {
	switch typed := list.(type) {
	case untypedListMsg:
		value, found := listAt(typed.v, index)
		return value, found
	case boolListMsg:
		value, found := listAt(typed.v, index)
		return NewBoolMsg(value), found
	case intListMsg:
		value, found := listAt(typed.v, index)
		return NewIntMsg(value), found
	case floatListMsg:
		value, found := listAt(typed.v, index)
		return NewFloatMsg(value), found
	case stringListMsg:
		value, found := listAt(typed.v, index)
		return NewStringMsg(value), found
	default:
		panic("unexpected list implementation")
	}
}

// ListSlice returns an immutable copy of the normalized range of list.
// Negative bounds count from the end, bounds are clamped, and reversed ranges
// return an empty list. It retains typed scalar storage and materializes a
// homogeneous scalar result from an untyped source.
//
//nolint:ireturn,varnamelen // Msg is the value-layer contract; from/to match the public slice component.
func ListSlice(list ListMsg, from, to int64) Msg {
	switch typed := list.(type) {
	case untypedListMsg:
		return NewListMsg(listSlice(typed.v, from, to))
	case boolListMsg:
		return NewListBoolMsg(listSlice(typed.v, from, to))
	case intListMsg:
		return NewListIntMsg(listSlice(typed.v, from, to))
	case floatListMsg:
		return NewListFloatMsg(listSlice(typed.v, from, to))
	case stringListMsg:
		return NewListStringMsg(listSlice(typed.v, from, to))
	default:
		panic("unexpected list implementation")
	}
}

// ListAppend returns a new list with value appended. It preserves typed scalar
// storage when value has the matching scalar representation. An untyped list
// is materialized when the resulting values are homogeneous scalars.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListAppend(list ListMsg, value Msg) Msg {
	switch typed := list.(type) {
	case untypedListMsg:
		return NewListMsg(listAppend(typed.v, value))
	case boolListMsg:
		if scalar, ok := value.(BoolMsg); ok {
			return NewListBoolMsg(listAppend(typed.v, scalar.v))
		}
	case intListMsg:
		if scalar, ok := value.(IntMsg); ok {
			return NewListIntMsg(listAppend(typed.v, scalar.v))
		}
	case floatListMsg:
		if scalar, ok := value.(FloatMsg); ok {
			return NewListFloatMsg(listAppend(typed.v, scalar.v))
		}
	case stringListMsg:
		if scalar, ok := value.(StringMsg); ok {
			return NewListStringMsg(listAppend(typed.v, scalar.v))
		}
	default:
		panic("unexpected list implementation")
	}

	return NewListMsg(listAppend(ListToMessageSlice(list), value))
}

// ListPrepend returns a new list with value prepended. It preserves typed
// scalar storage when value has the matching scalar representation. An
// untyped list is materialized when the resulting values are homogeneous
// scalars.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListPrepend(list ListMsg, value Msg) Msg {
	switch typed := list.(type) {
	case untypedListMsg:
		return NewListMsg(listPrepend(typed.v, value))
	case boolListMsg:
		if scalar, ok := value.(BoolMsg); ok {
			return NewListBoolMsg(listPrepend(typed.v, scalar.v))
		}
	case intListMsg:
		if scalar, ok := value.(IntMsg); ok {
			return NewListIntMsg(listPrepend(typed.v, scalar.v))
		}
	case floatListMsg:
		if scalar, ok := value.(FloatMsg); ok {
			return NewListFloatMsg(listPrepend(typed.v, scalar.v))
		}
	case stringListMsg:
		if scalar, ok := value.(StringMsg); ok {
			return NewListStringMsg(listPrepend(typed.v, scalar.v))
		}
	default:
		panic("unexpected list implementation")
	}

	return NewListMsg(listPrepend(ListToMessageSlice(list), value))
}

// ListConcat returns a new list that contains left followed by right. It
// preserves typed scalar storage when both lists use the same representation.
// An untyped input is materialized when the combined values are homogeneous
// scalars; otherwise it returns an untyped list.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListConcat(left, right ListMsg) Msg {
	switch leftTyped := left.(type) {
	case untypedListMsg:
		return NewListMsg(listConcat(leftTyped.v, ListToMessageSlice(right)))
	case boolListMsg:
		if rightTyped, ok := right.(boolListMsg); ok {
			return NewListBoolMsg(listConcat(leftTyped.v, rightTyped.v))
		}
	case intListMsg:
		if rightTyped, ok := right.(intListMsg); ok {
			return NewListIntMsg(listConcat(leftTyped.v, rightTyped.v))
		}
	case floatListMsg:
		if rightTyped, ok := right.(floatListMsg); ok {
			return NewListFloatMsg(listConcat(leftTyped.v, rightTyped.v))
		}
	case stringListMsg:
		if rightTyped, ok := right.(stringListMsg); ok {
			return NewListStringMsg(listConcat(leftTyped.v, rightTyped.v))
		}
	default:
		panic("unexpected list implementation")
	}

	return NewListMsg(listConcat(ListToMessageSlice(left), ListToMessageSlice(right)))
}

//nolint:ireturn // Generic helper returns a scalar storage value.
func listAt[T any](items []T, index int64) (T, bool) {
	length := int64(len(items))
	if index < -length || index >= length {
		var zero T
		return zero, false
	}
	if index < 0 {
		index += length
	}
	return items[index], true
}

func listSlice[T any](items []T, from, to int64) []T {
	start, end := listSliceBounds(from, to, int64(len(items)))
	return append([]T(nil), items[start:end]...)
}

func listAppend[T any](items []T, value T) []T {
	result := make([]T, len(items)+1)
	copy(result, items)
	result[len(items)] = value
	return result
}

func listPrepend[T any](items []T, value T) []T {
	result := make([]T, len(items)+1)
	result[0] = value
	copy(result[1:], items)
	return result
}

func listConcat[T any](left, right []T) []T {
	result := make([]T, len(left)+len(right))
	copy(result, left)
	copy(result[len(left):], right)
	return result
}

func listSliceBounds(from, to, length int64) (int64, int64) {
	start := listSliceIndex(from, length)
	end := listSliceIndex(to, length)
	if start > end {
		start = end
	}
	return start, end
}

func listSliceIndex(index, length int64) int64 {
	if index < 0 {
		index += length
	}
	if index < 0 {
		return 0
	}
	if index > length {
		return length
	}
	return index
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
