package core

import (
	"fmt"
	"strconv"
)

type Msg interface {
	Type() Type

	Bool() bool
	Int() int
	Str() string
	List() []Msg
	Struct() map[string]Msg

	String() string // for logging
}

type Type uint8

const (
	Bool Type = iota
	Int
	Str
	List
	Struct
)

type emptyMsg struct{}

func (msg emptyMsg) Int() int               { return 0 }
func (msg emptyMsg) Bool() bool             { return false }
func (msg emptyMsg) Str() string            { return "" }
func (msg emptyMsg) List() []Msg            { return []Msg{} }
func (msg emptyMsg) Struct() map[string]Msg { return map[string]Msg{} }

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int       { return msg.v }
func (msg IntMsg) Type() Type     { return Int }
func (msg IntMsg) String() string { return strconv.Itoa(msg.v) }

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

func (msg StrMsg) Str() string    { return msg.v }
func (msg StrMsg) Type() Type     { return Str }
func (msg StrMsg) String() string { return msg.v }

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

func (msg BoolMsg) Bool() bool     { return msg.v }
func (msg BoolMsg) Type() Type     { return Bool }
func (msg BoolMsg) String() string { return fmt.Sprint(msg.v) }

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

func (msg StructMsg) Struct() map[string]Msg { return msg.v }
func (msg StructMsg) Type() Type             { return Struct }
func (msg StructMsg) String() string         { return fmt.Sprint(msg.v) }

func NewStructMsg(v map[string]Msg) StructMsg {
	return StructMsg{
		emptyMsg: emptyMsg{},
		v:        map[string]Msg{},
	}
}

type ListMsg struct {
	emptyMsg
	v []Msg
}

func (msg ListMsg) List() []Msg { return msg.v }
func (msg ListMsg) Type() Type  { return List }
func (msg ListMsg) String() string {
	// b := &strings.Builder{}
	// b.WriteString("[")
	// for _, el := range msg.v {
	// 	fmt.Fprintf(b, "%s,", el.String())
	// }
	// b.WriteString("]")
	// return b.String()
	return fmt.Sprintf("[%s]", msg.v)
}

func NewListMsg(b bool) ListMsg {
	return ListMsg{
		emptyMsg: emptyMsg{},
		v:        []Msg{},
	}
}

func Eq(a, b Msg) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Type() {
	case Bool:
		return a.Bool() == b.Bool()
	case Int:
		return a.Int() == b.Int()
	case Str:
		return a.Str() == b.Str()
	case List:
		l1 := a.List()
		l2 := b.List()
		if len(l1) != len(l2) {
			return false
		}
		for i := range l1 {
			if !Eq(l1[i], l2[i]) {
				return false
			}
		}
	case Struct:
		s1 := a.Struct()
		s2 := a.Struct()
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
