package runtime

import (
	"fmt"
	"strconv"
)

// Msg methods don't return errors because they can be used not only at startup.
// If runtime functions need to validate message at startup, they must do it by themselves.
type Msg interface {
	fmt.Stringer
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	Type() MsgType
}

type MsgType uint8

const (
	UnknownMsgType MsgType = 0
	BoolMsgType    MsgType = 1
	IntMsgType     MsgType = 2
	FloatMsgType   MsgType = 3
	StrMsgType     MsgType = 4
)

// Empty

type emptyMsg struct{}

func (emptyMsg) Bool() bool     { return false }
func (emptyMsg) Int() int64     { return 0 }
func (emptyMsg) Float() float64 { return 0 }
func (emptyMsg) Str() string    { return "" }
func (emptyMsg) String() string { return "empty" } // stringer
func (emptyMsg) Type() MsgType  { return UnknownMsgType }

// Int

type IntMsg struct {
	emptyMsg
	v int64
}

func (msg IntMsg) Type() MsgType  { return IntMsgType }
func (msg IntMsg) Int() int64     { return msg.v }
func (msg IntMsg) String() string { return strconv.Itoa(int(msg.v)) }

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
		v:        n,
	}
}

// Float

type FloatMsg struct {
	emptyMsg
	v float64
}

func (msg FloatMsg) Type() MsgType     { return FloatMsgType }
func (msg FloatMsg) FloatMsg() float64 { return msg.v }

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		emptyMsg: emptyMsg{},
		v:        n,
	}
}

// Str

type StrMsg struct {
	emptyMsg
	v string
}

func (msg StrMsg) Type() MsgType  { return StrMsgType }
func (msg StrMsg) Str() string    { return msg.v }
func (msg StrMsg) String() string { return msg.v }

func NewStrMsg(s string) StrMsg {
	return StrMsg{
		emptyMsg: emptyMsg{},
		v:        s,
	}
}

// Bool

type BoolMsg struct {
	emptyMsg
	v bool
}

func (msg BoolMsg) Type() MsgType  { return BoolMsgType }
func (msg BoolMsg) Bool() bool     { return msg.v }
func (msg BoolMsg) String() string { return fmt.Sprint(msg.v) }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        b,
	}
}
