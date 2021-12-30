package runtime

import "fmt"

type Msg interface {
	Sig() struct{}
	Bool() bool
	Int() int
	Str() string
}

type emptyMsg struct{}

func (msg emptyMsg) Str() (_ string)   { return }
func (msg emptyMsg) Int() (_ int)      { return }
func (msg emptyMsg) Bool() (_ bool)    { return }
func (msg emptyMsg) Sig() (_ struct{}) { return }

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int {
	return msg.v
}

func (msg IntMsg) String() string {
	return fmt.Sprintf("'%d'", msg.v)
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

func (msg BoolMsg) String() string {
	return fmt.Sprintf("%v", msg.v)
}

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        b,
	}
}

type SigMsg struct {
	emptyMsg
}

func (msg SigMsg) Sig() struct{} {
	return struct{}{}
}

func NewSigMsg() SigMsg {
	return SigMsg{emptyMsg{}}
}
