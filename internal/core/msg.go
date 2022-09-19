package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Msg interface {
	Type() Type

	Bool() bool
	Int() int
	Str() string
	List() []Msg
	Dict() map[string]Msg

	String() string // for logging (move to interceptor?)
}

type Type uint8

const (
	BoolMsgType Type = iota
	IntMsgType
	StrMsgType
	ListMsgType
	DictMsgType
)

/* --- EMPTY --- */

type emptyMsg struct{}

func (msg emptyMsg) Int() int             { return 0 }
func (msg emptyMsg) Bool() bool           { return false }
func (msg emptyMsg) Str() string          { return "" }
func (msg emptyMsg) List() []Msg          { return []Msg{} }
func (msg emptyMsg) Dict() map[string]Msg { return map[string]Msg{} }

/* --- INT --- */

type IntMsg struct {
	emptyMsg
	v int
}

func (msg IntMsg) Int() int       { return msg.v }
func (msg IntMsg) Type() Type     { return IntMsgType }
func (msg IntMsg) String() string { return strconv.Itoa(msg.v) } // FIXME this broke printing from operators (this is for logging)

func NewIntMsg(n int) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
		v:        n,
	}
}

/* --- STR --- */

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

/* --- BOOL --- */

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

/* --- DICT --- */

type DictMsg struct {
	emptyMsg
	v map[string]Msg
}

func (msg DictMsg) Dict() map[string]Msg { return msg.v }
func (msg DictMsg) Type() Type           { return DictMsgType }
func (msg DictMsg) String() string { // can't move this to interceptor because print operator needs this
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

func NewDictMsg(v map[string]Msg) DictMsg {
	return DictMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

/* --- LIST --- */

type ListMsg struct {
	emptyMsg
	v []Msg
}

func (msg ListMsg) List() []Msg { return msg.v }
func (msg ListMsg) Type() Type  { return ListMsgType }
func (msg ListMsg) String() string {
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

func NewListMsg(v []Msg) ListMsg {
	return ListMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

/* --- OTHER --- */

func Eq(a, b Msg) bool { // maybe rewrite as a method?
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
	case ListMsgType:
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
	case DictMsgType:
		s1 := a.Dict()
		s2 := a.Dict()
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
