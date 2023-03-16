package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

type Msg interface {
	fmt.Stringer
	Type() Type
	Bool() bool
	Int() int
	Float() float64
	Str() string
	Vec() []Msg
	Map() map[string]Msg
}

type Type uint8

const (
	BoolMsgType Type = iota
	IntMsgType
	FloatMsgType
	StrMsgType
	VecMsgType
	MapMsgType
)

// Empty

type emptyMsg struct{}

func (emptyMsg) Bool() bool          { return false }
func (emptyMsg) Int() int            { return 0 }
func (emptyMsg) Float() float64      { return 0 }
func (emptyMsg) Str() string         { return "" }
func (emptyMsg) Vec() []Msg          { return []Msg{} }
func (emptyMsg) Map() map[string]Msg { return map[string]Msg{} }

// Int

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int       { return msg.v }
func (msg IntMsg) Type() Type     { return IntMsgType }
func (msg IntMsg) String() string { return strconv.Itoa(int(msg.v)) }

func NewIntMsg(n int) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
		v:        n,
	}
}

// Str

type StrMsg struct {
	emptyMsg
	v string
}

func (msg StrMsg) Str() string    { return msg.v }
func (msg StrMsg) Type() Type     { return StrMsgType }
func (msg StrMsg) String() string { return strconv.Quote(msg.v) }

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

func (msg BoolMsg) Bool() bool     { return msg.v }
func (msg BoolMsg) Type() Type     { return BoolMsgType }
func (msg BoolMsg) String() string { return fmt.Sprint(msg.v) }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        b,
	}
}

// Map

type MapMsg struct {
	emptyMsg
	v map[string]Msg
}

func (msg MapMsg) Map() map[string]Msg { return msg.v }
func (msg MapMsg) Type() Type          { return MapMsgType }
func (msg MapMsg) String() string {
	b := &strings.Builder{}
	b.WriteString("{")
	c := 0
	for k, el := range msg.v {
		c++
		if c < len(msg.v) {
			fmt.Fprintf(b, " %s: %s, ", k, el.String())
			continue
		}
		fmt.Fprintf(b, "%s: %s ", k, el.String())
	}
	b.WriteString("}")
	return b.String()
}

func NewMapMsg(v map[string]Msg) MapMsg {
	return MapMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

// Vec

type VecMsg struct {
	emptyMsg
	v []Msg
}

func (msg VecMsg) Vec() []Msg { return msg.v }
func (msg VecMsg) Type() Type { return VecMsgType }
func (msg VecMsg) String() string {
	b := &strings.Builder{}
	b.WriteString("[")
	c := 0
	for _, el := range msg.v {
		c++
		if c < len(msg.v) {
			fmt.Fprintf(b, "%s, ", el.String())
			continue
		}
		fmt.Fprint(b, el.String())
	}
	b.WriteString("]")
	return b.String()
}

func NewVecMsg(v []Msg) VecMsg {
	return VecMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

// Eq (NotEq, Less, Greater, Max, Min?)

func Eq(a, b Msg) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Type() {
	case BoolMsgType:
		return a.Bool() == b.Bool()
	case IntMsgType:
		return a.Int() == b.Int()
	case StrMsgType:
		return a.Str() == b.Str()
	case VecMsgType:
		l1 := a.Vec()
		l2 := b.Vec()
		if len(l1) != len(l2) {
			return false
		}
		for i := range l1 {
			if !Eq(l1[i], l2[i]) {
				return false
			}
		}
	case MapMsgType:
		s1 := a.Map()
		s2 := a.Map()
		if len(s1) != len(s2) {
			return false
		}
		for k := range s1 {
			if !Eq(s1[k], s2[k]) {
				return false
			}
		}
	}
	return false
}
