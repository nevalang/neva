package runtime

import "fbp/internal/types"

// Msg that goes from outport to inports.
type Msg struct {
	Type  types.Type
	Value Values
}

// Values represents different possible values for message.
type Values struct {
	Str  string // empty string if Message.Type != "str"
	Int  int    // 0 if Message.Type != "int"
	Bool bool   // false if Message.Type != "bool"
}
