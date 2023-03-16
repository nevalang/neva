package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

type Msg interface {
	fmt.Stringer
	Bool() bool
	Int() int
	Float() float64
	Str() string
	Vec() []Msg
	Map() map[string]Msg
}

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
