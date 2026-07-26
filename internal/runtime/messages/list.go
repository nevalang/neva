package messages

// ListToMsgs returns the boxed elements of list.
//
// For an untyped list, it returns the existing boxed representation. For a
// typed scalar list, it allocates a boxed slice and converts every element.
func ListToMsgs(list ListMsg) []Msg {
	if values, ok := AsListInts(list); ok {
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewIntMsg(values[i])
		}
		return msgs
	}
	if values, ok := AsListStrings(list); ok {
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewStringMsg(values[i])
		}
		return msgs
	}
	if values, ok := AsListBools(list); ok {
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewBoolMsg(values[i])
		}
		return msgs
	}
	if values, ok := AsListFloats(list); ok {
		msgs := make([]Msg, len(values))
		for i := range values {
			msgs[i] = NewFloatMsg(values[i])
		}
		return msgs
	}

	return list.Untyped()
}
