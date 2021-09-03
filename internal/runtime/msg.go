package runtime

// Msg represents set of methods only one of which should return real value.
type Msg interface {
	Str() string
	Int() int
	Bool() bool
	Struct() map[string]Msg
}

type emptyMsg struct{} // emptyMsg exists to allow normal messages define only reasonable methods.

func (msg emptyMsg) Str() string            { return "" }
func (msg emptyMsg) Int() int               { return 0 }
func (msg emptyMsg) Bool() bool             { return false }
func (msg emptyMsg) Struct() map[string]Msg { return nil }

var empty = emptyMsg{} // To avoid initialization of multiple empty messages.

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int {
	return msg.v
}

func NewIntMsg(n int) IntMsg {
	return IntMsg{
		emptyMsg: empty,
		v:        n,
	}
}

type StrMsg struct {
	emptyMsg
	v string
}

func (msg StrMsg) Str() string {
	return msg.v
}

func NewStrMsg(s string) StrMsg {
	return StrMsg{
		emptyMsg: empty,
		v:        s,
	}
}

type BoolMsg struct {
	emptyMsg
	v bool
}

func (msg BoolMsg) Bool() bool {
	return msg.v
}

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: empty,
		v:        b,
	}
}

type StructMsg struct {
	emptyMsg
	v map[string]Msg
}

func (msg StructMsg) Struct() map[string]Msg {
	return msg.v
}

func NewMsgStruct(v map[string]Msg) StructMsg {
	return StructMsg{
		emptyMsg: empty,
		v:        v,
	}
}
