package runtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type IndexedMsg struct {
	data  Msg
	index uint64 // to keep order of messages
}

type Msg interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	List() []Msg
	Map() map[string]Msg // TODO rename maps to dicts
	// IDEA use reflect for structures (instead of maps)
}

// Base (internal helper)

type baseMsg struct{}

func (baseMsg) String() string      { return "<empty>" }
func (baseMsg) Bool() bool          { return false }
func (baseMsg) Int() int64          { return 0 }
func (baseMsg) Float() float64      { return 0 }
func (baseMsg) Str() string         { return "" }
func (baseMsg) List() []Msg         { return nil }
func (baseMsg) Map() map[string]Msg { return nil }

// Int

type IntMsg struct {
	*baseMsg
	v int64
}

func (msg IntMsg) Int() int64                   { return msg.v }
func (msg IntMsg) String() string               { return strconv.Itoa(int(msg.v)) }
func (msg IntMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		baseMsg: &baseMsg{},
		v:       n,
	}
}

// Float

type FloatMsg struct {
	*baseMsg
	v float64
}

func (msg FloatMsg) Float() float64               { return msg.v }
func (msg FloatMsg) String() string               { return fmt.Sprint(msg.v) }
func (msg FloatMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		baseMsg: &baseMsg{},
		v:       n,
	}
}

// Str

type StrMsg struct {
	*baseMsg
	v string
}

func (msg StrMsg) Str() string    { return msg.v }
func (msg StrMsg) String() string { return msg.v }
func (msg StrMsg) MarshalJSON() ([]byte, error) {
	// Escape special characters in the string before marshaling
	jsonString, err := json.Marshal(msg.String())
	if err != nil {
		return nil, err
	}
	return jsonString, nil
}

func NewStrMsg(s string) StrMsg {
	return StrMsg{
		baseMsg: &baseMsg{},
		v:       s,
	}
}

// Bool

type BoolMsg struct {
	*baseMsg
	v bool
}

func (msg BoolMsg) Bool() bool                   { return msg.v }
func (msg BoolMsg) String() string               { return strconv.FormatBool(msg.v) }
func (msg BoolMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		baseMsg: &baseMsg{},
		v:       b,
	}
}

// List
type ListMsg struct {
	*baseMsg
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
		baseMsg: &baseMsg{},
		v:       v,
	}
}

// Map
type MapMsg struct {
	*baseMsg
	v map[string]Msg
}

func (msg MapMsg) Map() map[string]Msg { return msg.v }
func (msg MapMsg) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(msg.v)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, ":", ": ")
	jsonString = strings.ReplaceAll(jsonString, ",", ", ")

	return []byte(jsonString), nil
}
func (msg MapMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewMapMsg(m map[string]Msg) MapMsg {
	return MapMsg{
		baseMsg: &baseMsg{},
		v:       m,
	}
}
