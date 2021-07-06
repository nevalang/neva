package runtime

// Msg represents different possible values for message.
type Msg struct {
	Str  string // empty string if Message.Type != "str"
	Int  int    // 0 if Message.Type != "int"
	Bool bool   // false if Message.Type != "bool"
}
