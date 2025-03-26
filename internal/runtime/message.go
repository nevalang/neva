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
	Union() UnionMsg

	Equal(Msg) bool
}

// Internal

type internalMsg struct{}

func (internalMsg) String() string { panic("unexpected String method call on internal message type") }
func (internalMsg) Bool() bool     { panic("unexpected Bool method call on internal message type") }
func (internalMsg) Int() int64     { panic("unexpected Int method call on internal message type") }
func (internalMsg) Float() float64 { panic("unexpected Float method call on internal message type") }
func (internalMsg) Str() string    { panic("unexpected Str method call on internal message type") }
func (internalMsg) List() []Msg    { panic("unexpected List method call on internal message type") }
func (internalMsg) Dict() map[string]Msg {
	panic("unexpected Dict method call on internal message type")
}
func (internalMsg) Struct() StructMsg {
	panic("unexpected Struct method call on internal message type")
}
func (internalMsg) Union() UnionMsg { panic("unexpected Union method call on internal message type") }
func (internalMsg) Equal(other Msg) bool {
	panic("unexpected Equal method call on internal message type")
}

// Bool

type BoolMsg struct {
	internalMsg
	v bool
}

func (msg BoolMsg) Bool() bool                   { return msg.v }
func (msg BoolMsg) String() string               { return strconv.FormatBool(msg.v) }
func (msg BoolMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }
func (msg BoolMsg) Equal(other Msg) bool {
	otherBool, ok := other.(BoolMsg)
	return ok && msg.v == otherBool.v
}

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		internalMsg: internalMsg{},
		v:           b,
	}
}

// Int

type IntMsg struct {
	internalMsg
	v int64
}

func (msg IntMsg) Int() int64                   { return msg.v }
func (msg IntMsg) String() string               { return strconv.Itoa(int(msg.v)) }
func (msg IntMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }
func (msg IntMsg) Equal(other Msg) bool {
	otherInt, ok := other.(IntMsg)
	return ok && msg.v == otherInt.v
}

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		internalMsg: internalMsg{},
		v:           n,
	}
}

// Float

type FloatMsg struct {
	internalMsg
	v float64
}

func (msg FloatMsg) Float() float64               { return msg.v }
func (msg FloatMsg) String() string               { return fmt.Sprint(msg.v) }
func (msg FloatMsg) MarshalJSON() ([]byte, error) { return []byte(msg.String()), nil }
func (msg FloatMsg) Equal(other Msg) bool {
	otherFloat, ok := other.(FloatMsg)
	return ok && msg.v == otherFloat.v
}

func NewFloatMsg(n float64) FloatMsg {
	return FloatMsg{
		internalMsg: internalMsg{},
		v:           n,
	}
}

// --- STRING ---
type StringMsg struct {
	internalMsg
	v string
}

func (msg StringMsg) Str() string { return msg.v }

func (msg StringMsg) String() string { return msg.v }

func (msg StringMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.String())
}

func (msg StringMsg) Equal(other Msg) bool {
	otherString, ok := other.(StringMsg)
	return ok && msg.v == otherString.v
}

func NewStringMsg(s string) StringMsg {
	return StringMsg{
		internalMsg: internalMsg{},
		v:           s,
	}
}

// --- LIST ---
type ListMsg struct {
	internalMsg
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
func (msg ListMsg) Equal(other Msg) bool {
	otherList, ok := other.(ListMsg)
	if !ok {
		return false
	}
	if len(msg.v) != len(otherList.v) {
		return false
	}
	for i, v := range msg.v {
		if !v.Equal(otherList.v[i]) {
			return false
		}
	}
	return true
}

func NewListMsg(v []Msg) ListMsg {
	return ListMsg{
		internalMsg: internalMsg{},
		v:           v,
	}
}

// --- DICT ---
type DictMsg struct {
	internalMsg
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
func (msg DictMsg) Equal(other Msg) bool {
	otherDict, ok := other.(DictMsg)
	if !ok {
		return false
	}
	if len(msg.v) != len(otherDict.v) {
		return false
	}
	for k, v := range msg.v {
		otherV, ok := otherDict.v[k]
		if !ok || !v.Equal(otherV) {
			return false
		}
	}
	return true
}

func NewDictMsg(d map[string]Msg) DictMsg {
	return DictMsg{
		internalMsg: internalMsg{},
		v:           d,
	}
}

// --- STRUCT ---
type StructMsg struct {
	internalMsg
	names  []string // must be sorted for binary search
	fields []Msg    // must be equal length to names
}

func (msg StructMsg) Struct() StructMsg { return msg }

// Get returns the value of a field by name.
// It panics if the field is not found.
// It uses binary search to find the field, assuming the names are sorted.
func (msg StructMsg) Get(name string) Msg {
	if field, ok := msg.get(name); ok {
		return field
	}
	panic(fmt.Sprintf("field %q not found", name))
}

func (msg StructMsg) get(name string) (Msg, bool) {
	for i, n := range msg.names {
		if n == name {
			return msg.fields[i], true
		}
	}
	return nil, false
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

// Equal implements strict equality for StructMsg messages.
// It returns false if the lengths of the names and fields are different.
// It returns false if any of the fields are not equal.
func (msg StructMsg) Equal(other Msg) bool {
	otherStruct, ok := other.(StructMsg)
	if !ok {
		return false
	}
	if len(msg.names) != len(otherStruct.names) {
		return false
	}
	for i, name := range msg.names {
		otherField, ok := otherStruct.get(name)
		if !ok {
			return false
		}
		if !msg.fields[i].Equal(otherField) {
			return false
		}
	}
	return true
}

func NewStructMsg(names []string, fields []Msg) StructMsg {
	if len(names) != len(fields) {
		panic("names and fields must have the same length")
	}
	return StructMsg{
		internalMsg: internalMsg{},
		names:       names,
		fields:      fields,
	}
}

// --- UNION ---
type UnionMsg struct {
	internalMsg
	tag  string
	data Msg
}

func (msg UnionMsg) Union() UnionMsg {
	return msg
}

func (msg UnionMsg) Tag() string {
	return msg.tag
}

func (msg UnionMsg) Data() Msg {
	return msg.data
}

func (msg UnionMsg) String() string {
	return fmt.Sprintf(`{ "tag": "%s", "data": %v }`, msg.tag, msg.data)
}

// Equal implements strict equality for UnionMsg messages.
// If one union has data and another doesn't, it returns false.
// It returns false if tags are different.
// It returns false if data is different.
// Tags are compared as Go strings and data is compared recursevely using Equal method.
func (msg UnionMsg) Equal(other Msg) bool {
	otherUnion, ok := other.(UnionMsg)
	if !ok {
		return false
	}

	if msg.data != nil && otherUnion.data == nil {
		return false
	} else if msg.data == nil && otherUnion.data != nil {
		return false
	}

	if msg.tag != otherUnion.tag {
		return false
	}

	if msg.data == nil {
		return true
	}

	return msg.data.Equal(otherUnion.data)
}

// Match implements pattern matching for UnionMsg messages.
// It handles case when one of the unions doesn't have data
// and compares only tags. Otherwise it uses Equal method.
func (msg UnionMsg) Match(pattern UnionMsg) bool {
	if msg.data != nil && pattern.data == nil ||
		msg.data == nil && pattern.data != nil {
		return msg.tag == pattern.tag
	}
	return msg.data.Equal(pattern.data)
}

func NewUnionMsg(tag string, data Msg) UnionMsg {
	return UnionMsg{
		internalMsg: internalMsg{},
		tag:         tag,
		data:        data,
	}
}

// --- OPERATIONS ---

func Match(msg Msg, pattern Msg) bool {
	msgUnion, ok := msg.(UnionMsg)
	if !ok {
		return msg.Equal(pattern)
	}

	patternUnion, ok := pattern.(UnionMsg)
	if !ok {
		return msg.Equal(pattern)
	}

	if msgUnion.data != nil && patternUnion.data == nil ||
		msgUnion.data == nil && patternUnion.data != nil {
		return msgUnion.tag == patternUnion.tag
	}

	return msgUnion.data.Equal(patternUnion.data)
}
