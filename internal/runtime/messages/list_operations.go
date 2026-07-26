package messages

// ListAt returns the item at index. Negative indexes count from the end.
// The boolean is false when index is outside the list bounds.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListAt(list ListMsg, index int64) (Msg, bool) {
	switch typed := list.(type) {
	case untypedListMsg:
		value, found := listItem(typed.v, index)
		return value, found
	case boolListMsg:
		value, found := listItem(typed.v, index)
		return NewBoolMsg(value), found
	case intListMsg:
		value, found := listItem(typed.v, index)
		return NewIntMsg(value), found
	case floatListMsg:
		value, found := listItem(typed.v, index)
		return NewFloatMsg(value), found
	case stringListMsg:
		value, found := listItem(typed.v, index)
		return NewStringMsg(value), found
	default:
		panic("unexpected list implementation")
	}
}

// ListSlice returns an immutable copy of the normalized range of list.
// Negative bounds count from the end, bounds are clamped, and reversed ranges
// return an empty list. Typed scalar storage remains typed.
//
//nolint:ireturn,varnamelen // Msg is the value-layer contract; from/to match the public slice component.
func ListSlice(list ListMsg, from, to int64) Msg {
	switch typed := list.(type) {
	case untypedListMsg:
		return NewListMsg(sliceList(typed.v, from, to))
	case boolListMsg:
		return NewListBoolMsg(sliceList(typed.v, from, to))
	case intListMsg:
		return NewListIntMsg(sliceList(typed.v, from, to))
	case floatListMsg:
		return NewListFloatMsg(sliceList(typed.v, from, to))
	case stringListMsg:
		return NewListStringMsg(sliceList(typed.v, from, to))
	default:
		panic("unexpected list implementation")
	}
}

// ListAppend returns a new list with value appended. It preserves typed scalar
// storage when value has the matching scalar representation; otherwise it
// returns an untyped list.
//
//nolint:ireturn // Msg is the value-layer contract.
func ListAppend(list ListMsg, value Msg) Msg {
	switch typed := list.(type) {
	case untypedListMsg:
		return NewListMsg(appendValue(typed.v, value))
	case boolListMsg:
		if scalar, ok := value.(BoolMsg); ok {
			return NewListBoolMsg(appendValue(typed.v, scalar.v))
		}
	case intListMsg:
		if scalar, ok := value.(IntMsg); ok {
			return NewListIntMsg(appendValue(typed.v, scalar.v))
		}
	case floatListMsg:
		if scalar, ok := value.(FloatMsg); ok {
			return NewListFloatMsg(appendValue(typed.v, scalar.v))
		}
	case stringListMsg:
		if scalar, ok := value.(StringMsg); ok {
			return NewListStringMsg(appendValue(typed.v, scalar.v))
		}
	default:
		panic("unexpected list implementation")
	}

	return NewListMsg(appendValue(ListToMessageSlice(list), value))
}

//nolint:ireturn // Generic helper returns a scalar storage value.
func listItem[T any](items []T, index int64) (T, bool) {
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

func sliceList[T any](items []T, from, to int64) []T {
	start, end := normalizeSliceBounds(from, to, int64(len(items)))
	return append([]T(nil), items[start:end]...)
}

func appendValue[T any](items []T, value T) []T {
	result := make([]T, len(items)+1)
	copy(result, items)
	result[len(items)] = value
	return result
}

func normalizeSliceBounds(from, to, length int64) (int64, int64) {
	start := normalizeSliceIndex(from, length)
	end := normalizeSliceIndex(to, length)
	if start > end {
		start = end
	}
	return start, end
}

func normalizeSliceIndex(index, length int64) int64 {
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
