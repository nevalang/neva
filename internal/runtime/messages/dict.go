package messages

// DictToMsgs returns the boxed entries of dict.
//
// For an untyped dictionary, it returns the existing boxed representation. For
// a typed scalar dictionary, it allocates a boxed map and converts every value.
func DictToMsgs(dict DictMsg) map[string]Msg {
	if values, ok := AsDictInts(dict); ok {
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewIntMsg(value)
		}
		return msgs
	}
	if values, ok := AsDictStrings(dict); ok {
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewStringMsg(value)
		}
		return msgs
	}
	if values, ok := AsDictBools(dict); ok {
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewBoolMsg(value)
		}
		return msgs
	}
	if values, ok := AsDictFloats(dict); ok {
		msgs := make(map[string]Msg, len(values))
		for key, value := range values {
			msgs[key] = NewFloatMsg(value)
		}
		return msgs
	}

	return dict.Untyped()
}

// DictGet returns the value stored for key without materializing typed storage.
// The boolean is false when key is not present.
//
//nolint:ireturn // Msg is the value-layer contract.
func DictGet(dict DictMsg, key string) (Msg, bool) {
	if values, ok := AsDictInts(dict); ok {
		value, found := values[key]
		return NewIntMsg(value), found
	}
	if values, ok := AsDictStrings(dict); ok {
		value, found := values[key]
		return NewStringMsg(value), found
	}
	if values, ok := AsDictBools(dict); ok {
		value, found := values[key]
		return NewBoolMsg(value), found
	}
	if values, ok := AsDictFloats(dict); ok {
		value, found := values[key]
		return NewFloatMsg(value), found
	}

	value, found := dict.Untyped()[key]
	return value, found
}
