package core

type Msg interface {
	Sig() struct{}
	Bool() bool
	Int() int
	Str() string
	Type() Type
}

type Type uint8

const (
	Sig Type = iota
	Bool
	Int
	Str
)

type emptyMsg struct{}

func (msg emptyMsg) Str() (_ string)   { return }
func (msg emptyMsg) Int() (_ int)      { return }
func (msg emptyMsg) Bool() (_ bool)    { return }
func (msg emptyMsg) Sig() (_ struct{}) { return }

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int   { return msg.v }
func (msg IntMsg) Type() Type { return Int }

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

func (msg StrMsg) Str() string { return msg.v }
func (msg StrMsg) Type() Type  { return Str }

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

func (msg BoolMsg) Bool() bool { return msg.v }
func (msg BoolMsg) Type() Type { return Bool }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        b,
	}
}

type SigMsg struct {
	emptyMsg
}

func (msg SigMsg) Sig() struct{} { return struct{}{} }
func (msg SigMsg) Type() Type    { return Sig }

func NewSigMsg() SigMsg {
	return SigMsg{emptyMsg{}}
}

func Eq(a, b Msg) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case Sig:
		return a.Sig() == b.Sig()
	case Bool:
		return a.Bool() == b.Bool()
	case Int:
		return a.Int() == b.Int()
	case Str:
		return a.Str() == b.Str()
	default:
		return false
	}
}
