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
	fields []StructField
}

func (msg StructMsg) Struct() StructMsg { return msg }

// get returns the value of a field by name.
// it panics if the field is not found.
// it uses linear scan to find the field.
func (msg StructMsg) Get(name string) Msg {
	if field, ok := msg.get(name); ok {
		return field
	}
	panic(fmt.Sprintf("field %q not found", name))
}

func (msg StructMsg) get(name string) (Msg, bool) {
	for i := range msg.fields {
		if msg.fields[i].name == name {
			return msg.fields[i].value, true
		}
	}
	return nil, false
}

func (msg StructMsg) MarshalJSON() ([]byte, error) {
	m := make(map[string]Msg, len(msg.fields))
	for i := range msg.fields {
		m[msg.fields[i].name] = msg.fields[i].value
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
	if len(msg.fields) != len(otherStruct.fields) {
		return false
	}
	for i := range msg.fields {
		otherField, ok := otherStruct.get(msg.fields[i].name)
		if !ok {
			return false
		}
		if !msg.fields[i].value.Equal(otherField) {
			return false
		}
	}
	return true
}

func newStructMsg(fields []StructField) StructMsg {
	if len(fields) == 0 {
		return StructMsg{internalMsg: internalMsg{}, fields: nil}
	}
	copied := make([]StructField, len(fields))
	copy(copied, fields)
	return StructMsg{
		internalMsg: internalMsg{},
		fields:      copied,
	}
}

// structfield is a helper to construct structs via runtime.newstruct api without exposing fields.
type StructField struct {
	value Msg
	name  string
}

// newstructfield constructs a structfield with provided name and value.
func NewStructField(name string, value Msg) StructField {
	return StructField{name: name, value: value}
}

// newstruct builds a struct message from a slice of structfield.
// underlying struct representation remains unchanged for now.
func NewStructMsg(fields []StructField) StructMsg { return newStructMsg(fields) }

// --- UNION ---
type UnionMsg struct {
	internalMsg
	data Msg
	tag  string
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
	if msg.data == nil {
		return fmt.Sprintf(`{ "tag": %q }`, msg.tag)
	}
	return fmt.Sprintf(`{ "tag": %q, "data": %v }`, msg.tag, msg.data)
}

// Uint8Index validates idx and returns it as uint8 or panics.
func Uint8Index(idx int) uint8 {
	if idx < 0 {
		panic(fmt.Sprintf("runtime: negative index %d", idx))
	}
	if idx > int(^uint8(0)) {
		panic(fmt.Sprintf("runtime: index %d overflows uint8", idx))
	}
	// #nosec G115 -- bounds checked above
	return uint8(idx)
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

func NewUnionMsg(tag string, data Msg) UnionMsg {
	return UnionMsg{
		internalMsg: internalMsg{},
		tag:         tag,
		data:        data,
	}
}

// --- OPERATIONS ---

// Match compares two messages and return true if they matches and false otherwise.
// Unlike Equal it compares only some aspects of the messages.
func Match(msg Msg, pattern Msg) bool {
	// at the moment we only match unions
	// maybe in the future we'll add support for more types e.g. structs
	msgUnion, ok := msg.(UnionMsg)
	if !ok {
		return msg.Equal(pattern)
	}

	// both msg and pattern must be unions to perform pattern matching
	// if at least one of them is not, strict equality will be applied instead
	patternUnion, ok := pattern.(UnionMsg)
	if !ok {
		return msg.Equal(pattern)
	}

	// if tags are not equal data does not matter, there's no match
	if msgUnion.tag != patternUnion.tag {
		return false
	}

	// if pattern doesn't have data we match by tag
	// and by this time we know tags are equal
	if patternUnion.data == nil {
		return true
	}

	// if we here we know that pattern has data
	// so if msg doesn't, they don't match
	if msgUnion.data == nil {
		return false
	}

	// by this time we know
	// both msg and pattern are union messages
	// they both have the same tags and some data inside
	// so we apply strict equality to the data they wrap
	// maybe in the future we'll consider recursive matching, we'll see
	return msgUnion.data.Equal(patternUnion.data)
}
