package core

type Msg interface {
	Str() string
	Int() int
	Bool() bool
}

type emptyMsg struct{}

func (msg emptyMsg) Str() string            { panic("not implemented") }
func (msg emptyMsg) Int() int               { panic("not implemented") }
func (msg emptyMsg) Bool() bool             { panic("not implemented") }
func (msg emptyMsg) Struct() map[string]Msg { panic("not implemented") }

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int { return msg.v }

func NewIntMsg(v int) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

type StrMsg struct {
	emptyMsg
	v string
}

func (msg StrMsg) Str() string {
	return msg.v
}

func NewStrMsg(v string) StrMsg {
	return StrMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

type BoolMsg struct {
	emptyMsg
	v bool
}

func (msg BoolMsg) Bool() bool {
	return msg.v
}

func NewBoolMsg(v bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}
