package runtime

type Msg interface {
	Str() string
	Int() int
	Bool() bool
	Struct() map[string]Msg
}

type emptyMsg struct{}

func (msg emptyMsg) Str() string {
	return ""
}

func (msg emptyMsg) Int() int {
	return 0
}

func (msg emptyMsg) Bool() bool {
	return false
}

func (msg emptyMsg) Struct() map[string]Msg {
	return nil
}

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int {
	return msg.v
}

func NewIntMsg(n int) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
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
		emptyMsg: emptyMsg{},
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
		emptyMsg: emptyMsg{},
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
		emptyMsg: emptyMsg{},
		v:        v,
	}
}
