package runtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// OrderedMsg is a message with a chronological index.
// Ordered messages can be compared and sorted by their index.
type OrderedMsg struct {
	Msg
	index uint64
}

func (o OrderedMsg) String() string {
	return fmt.Sprint(o.Msg)
}

type Msg interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	List() []Msg
	Dict() map[string]Msg
	Struct() StructMsg
}

// Base

type baseMsg struct{}

func (baseMsg) String() string       { return "<base>" }
func (baseMsg) Bool() bool           { return false }
func (baseMsg) Int() int64           { return 0 }
func (baseMsg) Float() float64       { return 0 }
func (baseMsg) Str() string          { return "" }
func (baseMsg) List() []Msg          { return nil }
func (baseMsg) Dict() map[string]Msg { return nil }
func (baseMsg) Struct() StructMsg    { return StructMsg{} }

// Bool

type BoolMsg struct {
	baseMsg
	v bool
}

func (msg BoolMsg) Bool() bool                   { return msg.v }
func (msg BoolMsg) String() string               { return strconv.FormatBool(msg.v) }
func (msg BoolMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		baseMsg: baseMsg{},
		v:       b,
	}
}

// Int

type IntMsg struct {
	baseMsg
	v int64
}

func (msg IntMsg) Int() int64                   { return msg.v }
func (msg IntMsg) String() string               { return strconv.Itoa(int(msg.v)) }
func (msg IntMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		baseMsg: baseMsg{},
		v:       n,
	}
}

// Float

type FloatMsg struct {
	baseMsg
	v float64
}

func (msg FloatMsg) Float() float64               { return msg.v }
func (msg FloatMsg) String() string               { return fmt.Sprint(msg.v) }
func (msg FloatMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		baseMsg: baseMsg{},
		v:       n,
	}
}

// Str

type StringMsg struct {
	baseMsg
	v string
}

func (msg StringMsg) Str() string    { return msg.v }
func (msg StringMsg) String() string { return msg.v }
func (msg StringMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.String())
}

func NewStringMsg(s string) StringMsg {
	return StringMsg{
		baseMsg: baseMsg{},
		v:       s,
	}
}

// List
type ListMsg struct {
	baseMsg
	v []Msg
}

func (msg ListMsg) List() []Msg { return msg.v }
func (msg ListMsg) String() string {
	bb, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bb)
}
func (msg ListMsg) MarshalJSON() ([]byte, error) { return json.Marshal(msg.v) }

func NewListMsg(v []Msg) ListMsg {
	return ListMsg{
		baseMsg: baseMsg{},
		v:       v,
	}
}

// Dictionary
type DictMsg struct {
	baseMsg
	v map[string]Msg
}

func (msg DictMsg) Dict() map[string]Msg { return msg.v }
func (msg DictMsg) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(msg.v)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, ":", ": ")
	jsonString = strings.ReplaceAll(jsonString, ",", ", ")

	return []byte(jsonString), nil
}
func (msg DictMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewDictMsg(d map[string]Msg) DictMsg {
	return DictMsg{
		baseMsg: baseMsg{},
		v:       d,
	}
}

// Structure
type StructMsg struct {
	baseMsg
	names  []string // must be sorted for binary search
	fields []Msg    // must be equal length to names
}

func (msg StructMsg) Struct() StructMsg { return msg }

// Get returns the value of a field by name.
// It panics if the field is not found.
// It uses binary search to find the field, assuming the names are sorted.
func (msg StructMsg) Get(name string) Msg {
	for i, n := range msg.names {
		if n == name {
			return msg.fields[i]
		}
	}
	panic(fmt.Sprintf("field %q not found", name))
}

func (msg StructMsg) MarshalJSON() ([]byte, error) {
	m := make(map[string]Msg, len(msg.names))
	for i, name := range msg.names {
		m[name] = msg.fields[i]
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, ":", ": ")
	jsonString = strings.ReplaceAll(jsonString, ",", ", ")

	return []byte(jsonString), nil
}

func (msg StructMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewStructMsg(names []string, fields []Msg) StructMsg {
	if len(names) != len(fields) {
		panic("names and fields must have the same length")
	}
	return StructMsg{
		baseMsg: baseMsg{},
		names:   names,
		fields:  fields,
	}
}
