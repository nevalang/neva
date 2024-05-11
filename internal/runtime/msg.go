package runtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Trace interface {
	Trace() []PortAddr
	AppendTrace(PortAddr)
}

type Msg interface {
	Trace
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	List() []Msg
	Map() map[string]Msg
}

// Base (internal helper)

type baseMsg struct {
	trace []PortAddr
}

func (baseMsg) String() string        { return "<empty>" }
func (baseMsg) Bool() bool            { return false }
func (baseMsg) Int() int64            { return 0 }
func (baseMsg) Float() float64        { return 0 }
func (baseMsg) Str() string           { return "" }
func (baseMsg) List() []Msg           { return nil }
func (baseMsg) Map() map[string]Msg   { return nil }
func (msg baseMsg) Trace() []PortAddr { return msg.trace }
func (msg *baseMsg) AppendTrace(pa PortAddr) {
	msg.trace = append(msg.trace, pa)
}

// Int

type IntMsg struct {
	*baseMsg
	v int64
}

func (msg IntMsg) Int() int64                   { return msg.v }
func (msg IntMsg) String() string               { return strconv.Itoa(int(msg.v)) }
func (msg IntMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }
func (msg IntMsg) Trace() []PortAddr            { return msg.baseMsg.trace }
func (msg IntMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
}

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
func (msg FloatMsg) Trace() []PortAddr            { return msg.baseMsg.trace }
func (msg FloatMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
}

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
func (msg StrMsg) Trace() []PortAddr { return msg.baseMsg.trace }
func (msg StrMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
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
func (msg BoolMsg) Trace() []PortAddr            { return msg.baseMsg.trace }
func (msg BoolMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
}

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

func (msg ListMsg) List() []Msg   { return msg.v }
func (msg ListMsg) String() string {
	bb, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bb)
}
func (msg ListMsg) MarshalJSON() ([]byte, error) { return json.Marshal(msg.v) }
func (msg ListMsg) Trace() []PortAddr            { return msg.baseMsg.trace }
func (msg ListMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
}

func NewListMsg(v ...Msg) ListMsg {
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
func (msg MapMsg) Trace() []PortAddr { return msg.baseMsg.trace }
func (msg MapMsg) AppendTrace(pa PortAddr) {
	msg.baseMsg.trace = append(msg.trace, pa)
}

func NewMapMsg(m map[string]Msg) MapMsg {
	return MapMsg{
		baseMsg: &baseMsg{},
		v:       m,
	}
}
