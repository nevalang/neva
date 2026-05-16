package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// OrderedMsg is a transport envelope with payload and runtime ordering metadata.
type OrderedMsg struct {
	Msg
	index uint64
}

// String is just a simple stringer that ignores index while formatting.
func (o OrderedMsg) String() string { return fmt.Sprint(o.Msg) }

//nolint:interfacebloat // Msg is runtime contract and intentionally broad.
type Msg interface {
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	Bytes() []byte
	List() ListMsg
	Dict() DictMsg
	Struct() StructMsg
	Union() UnionMsg

	MarshalJSON() ([]byte, error)
	Equal(Msg) bool
}

// Internal

type internalMsg struct {
}

func (internalMsg) String() string { panic("unexpected String method call on internal message type") }
func (internalMsg) Bool() bool     { panic("unexpected Bool method call on internal message type") }
func (internalMsg) Int() int64     { panic("unexpected Int method call on internal message type") }
func (internalMsg) Float() float64 { panic("unexpected Float method call on internal message type") }
func (internalMsg) Str() string    { panic("unexpected Str method call on internal message type") }
func (internalMsg) Bytes() []byte  { panic("unexpected Bytes method call on internal message type") }

//nolint:ireturn // Msg contract uses interfaces.
func (internalMsg) List() ListMsg { panic("unexpected List method call on internal message type") }

//nolint:ireturn // Msg contract uses interfaces.
func (internalMsg) Dict() DictMsg {
	panic("unexpected Dict method call on internal message type")
}
func (internalMsg) Struct() StructMsg {
	panic("unexpected Struct method call on internal message type")
}
func (internalMsg) Union() UnionMsg { panic("unexpected Union method call on internal message type") }
func (internalMsg) MarshalJSON() ([]byte, error) {
	panic("unexpected MarshalJSON method call on internal message type")
}
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
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type StringMsg struct {
	internalMsg
	v string
}

func (msg StringMsg) Str() string { return msg.v }

func (msg StringMsg) String() string { return msg.v }

func (msg StringMsg) MarshalJSON() ([]byte, error) {
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

// --- BYTES ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type BytesMsg struct {
	internalMsg
	v []byte
}

func (msg BytesMsg) Bytes() []byte { return msg.v }

func (msg BytesMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg BytesMsg) MarshalJSON() ([]byte, error) {
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return json.Marshal(msg.v)
}

func (msg BytesMsg) Equal(other Msg) bool {
	otherBytes, ok := other.(BytesMsg)
	return ok && bytes.Equal(msg.v, otherBytes.v)
}

func NewBytesMsg(v []byte) BytesMsg {
	return BytesMsg{
		internalMsg: internalMsg{},
		v:           v,
	}
}

// --- LIST ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type ListMsg interface {
	Msgs() []Msg
	Bools() []bool
	Ints() []int64
	Floats() []float64
	Strings() []string
	Len() int
	String() string
	MarshalJSON() ([]byte, error)
	Equal(ListMsg) bool
}

//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type listValueMsg struct {
	internalMsg
	v ListMsg
}

//nolint:ireturn // Msg contract uses interfaces.
func (msg listValueMsg) List() ListMsg { return msg.v }

func (msg listValueMsg) String() string {
	return msg.v.String()
}

//nolint:wrapcheck // Delegates to list implementation.
func (msg listValueMsg) MarshalJSON() ([]byte, error) {
	return msg.v.MarshalJSON()
}

func (msg listValueMsg) Equal(other Msg) bool {
	otherList, ok := other.(listValueMsg)
	return ok && msg.v.Equal(otherList.v)
}

type genericListMsg struct{ v []Msg }
type boolListMsg struct{ v []bool }
type intListMsg struct{ v []int64 }
type floatListMsg struct{ v []float64 }
type stringListMsg struct{ v []string }

func (msg genericListMsg) Msgs() []Msg { return msg.v }
func (genericListMsg) Bools() []bool   { panic("unexpected Bools method call on generic list message") }
func (genericListMsg) Ints() []int64   { panic("unexpected Ints method call on generic list message") }
func (genericListMsg) Floats() []float64 {
	panic("unexpected Floats method call on generic list message")
}
func (genericListMsg) Strings() []string {
	panic("unexpected Strings method call on generic list message")
}
func (msg genericListMsg) Len() int { return len(msg.v) }
func (msg genericListMsg) String() string {
	bb, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bb)
}

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg genericListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg genericListMsg) Equal(other ListMsg) bool {
	return listEqualGeneric(msg.v, other)
}

func (boolListMsg) Msgs() []Msg        { panic("unexpected Msgs method call on bool list message") }
func (msg boolListMsg) Bools() []bool  { return msg.v }
func (boolListMsg) Ints() []int64      { panic("unexpected Ints method call on bool list message") }
func (boolListMsg) Floats() []float64  { panic("unexpected Floats method call on bool list message") }
func (boolListMsg) Strings() []string  { panic("unexpected Strings method call on bool list message") }
func (msg boolListMsg) Len() int       { return len(msg.v) }
func (msg boolListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg boolListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg boolListMsg) Equal(other ListMsg) bool { return listEqualBool(msg.v, other) }

func (intListMsg) Msgs() []Msg        { panic("unexpected Msgs method call on int list message") }
func (intListMsg) Bools() []bool      { panic("unexpected Bools method call on int list message") }
func (msg intListMsg) Ints() []int64  { return msg.v }
func (intListMsg) Floats() []float64  { panic("unexpected Floats method call on int list message") }
func (intListMsg) Strings() []string  { panic("unexpected Strings method call on int list message") }
func (msg intListMsg) Len() int       { return len(msg.v) }
func (msg intListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg intListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg intListMsg) Equal(other ListMsg) bool { return listEqualInt(msg.v, other) }

func (floatListMsg) Msgs() []Msg   { panic("unexpected Msgs method call on float list message") }
func (floatListMsg) Bools() []bool { panic("unexpected Bools method call on float list message") }
func (floatListMsg) Ints() []int64 { panic("unexpected Ints method call on float list message") }
func (msg floatListMsg) Floats() []float64 {
	return msg.v
}
func (floatListMsg) Strings() []string  { panic("unexpected Strings method call on float list message") }
func (msg floatListMsg) Len() int       { return len(msg.v) }
func (msg floatListMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg floatListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg floatListMsg) Equal(other ListMsg) bool { return listEqualFloat(msg.v, other) }

func (stringListMsg) Msgs() []Msg   { panic("unexpected Msgs method call on string list message") }
func (stringListMsg) Bools() []bool { panic("unexpected Bools method call on string list message") }
func (stringListMsg) Ints() []int64 { panic("unexpected Ints method call on string list message") }
func (stringListMsg) Floats() []float64 {
	panic("unexpected Floats method call on string list message")
}
func (msg stringListMsg) Strings() []string { return msg.v }
func (msg stringListMsg) Len() int          { return len(msg.v) }
func (msg stringListMsg) String() string    { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg stringListMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg stringListMsg) Equal(other ListMsg) bool { return listEqualString(msg.v, other) }

//nolint:ireturn // Msg contract type.
func NewListMsg(v []Msg) Msg {
	return listValueMsg{internalMsg: internalMsg{}, v: genericListMsg{v: v}}
}

//nolint:ireturn // Msg contract type.
func NewListBoolMsg(v []bool) Msg {
	return listValueMsg{internalMsg: internalMsg{}, v: boolListMsg{v: v}}
}

//nolint:ireturn // Msg contract type.
func NewListIntMsg(v []int64) Msg {
	return listValueMsg{internalMsg: internalMsg{}, v: intListMsg{v: v}}
}

//nolint:ireturn // Msg contract type.
func NewListFloatMsg(v []float64) Msg {
	return listValueMsg{internalMsg: internalMsg{}, v: floatListMsg{v: v}}
}

//nolint:ireturn // Msg contract type.
func NewListStringMsg(v []string) Msg {
	return listValueMsg{internalMsg: internalMsg{}, v: stringListMsg{v: v}}
}

// --- DICT ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type DictMsg interface {
	Msgs() map[string]Msg
	Bools() map[string]bool
	Ints() map[string]int64
	Floats() map[string]float64
	Strings() map[string]string
	Len() int
	String() string
	MarshalJSON() ([]byte, error)
	Equal(DictMsg) bool
}

//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type dictValueMsg struct {
	internalMsg
	v DictMsg
}

//nolint:ireturn // Msg contract uses interfaces.
func (msg dictValueMsg) Dict() DictMsg  { return msg.v }
func (msg dictValueMsg) String() string { return msg.v.String() }

//nolint:wrapcheck // Delegates to dict implementation.
func (msg dictValueMsg) MarshalJSON() ([]byte, error) { return msg.v.MarshalJSON() }
func (msg dictValueMsg) Equal(other Msg) bool {
	otherDict, ok := other.(dictValueMsg)
	return ok && msg.v.Equal(otherDict.v)
}

type genericDictMsg struct{ v map[string]Msg }
type boolDictMsg struct{ v map[string]bool }
type intDictMsg struct{ v map[string]int64 }
type floatDictMsg struct{ v map[string]float64 }
type stringDictMsg struct{ v map[string]string }

func (msg genericDictMsg) Msgs() map[string]Msg { return msg.v }
func (genericDictMsg) Bools() map[string]bool {
	panic("unexpected Bools method call on generic dict message")
}
func (genericDictMsg) Ints() map[string]int64 {
	panic("unexpected Ints method call on generic dict message")
}
func (genericDictMsg) Floats() map[string]float64 {
	panic("unexpected Floats method call on generic dict message")
}
func (genericDictMsg) Strings() map[string]string {
	panic("unexpected Strings method call on generic dict message")
}
func (msg genericDictMsg) Len() int       { return len(msg.v) }
func (msg genericDictMsg) String() string { return mustJSON(msg) }
func (msg genericDictMsg) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(msg.v)
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	}
	return addJSONSpaces(jsonData), nil
}
func (msg genericDictMsg) Equal(other DictMsg) bool { return dictEqual(msg, other) }

func (boolDictMsg) Msgs() map[string]Msg       { panic("unexpected Msgs method call on bool dict message") }
func (msg boolDictMsg) Bools() map[string]bool { return msg.v }
func (boolDictMsg) Ints() map[string]int64     { panic("unexpected Ints method call on bool dict message") }
func (boolDictMsg) Floats() map[string]float64 {
	panic("unexpected Floats method call on bool dict message")
}
func (boolDictMsg) Strings() map[string]string {
	panic("unexpected Strings method call on bool dict message")
}
func (msg boolDictMsg) Len() int       { return len(msg.v) }
func (msg boolDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg boolDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg boolDictMsg) Equal(other DictMsg) bool { return dictEqual(msg, other) }

func (intDictMsg) Msgs() map[string]Msg       { panic("unexpected Msgs method call on int dict message") }
func (intDictMsg) Bools() map[string]bool     { panic("unexpected Bools method call on int dict message") }
func (msg intDictMsg) Ints() map[string]int64 { return msg.v }
func (intDictMsg) Floats() map[string]float64 {
	panic("unexpected Floats method call on int dict message")
}
func (intDictMsg) Strings() map[string]string {
	panic("unexpected Strings method call on int dict message")
}
func (msg intDictMsg) Len() int       { return len(msg.v) }
func (msg intDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg intDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg intDictMsg) Equal(other DictMsg) bool { return dictEqual(msg, other) }

func (floatDictMsg) Msgs() map[string]Msg { panic("unexpected Msgs method call on float dict message") }
func (floatDictMsg) Bools() map[string]bool {
	panic("unexpected Bools method call on float dict message")
}
func (floatDictMsg) Ints() map[string]int64 {
	panic("unexpected Ints method call on float dict message")
}
func (msg floatDictMsg) Floats() map[string]float64 { return msg.v }
func (floatDictMsg) Strings() map[string]string {
	panic("unexpected Strings method call on float dict message")
}
func (msg floatDictMsg) Len() int       { return len(msg.v) }
func (msg floatDictMsg) String() string { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg floatDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg floatDictMsg) Equal(other DictMsg) bool { return dictEqual(msg, other) }

func (stringDictMsg) Msgs() map[string]Msg {
	panic("unexpected Msgs method call on string dict message")
}
func (stringDictMsg) Bools() map[string]bool {
	panic("unexpected Bools method call on string dict message")
}
func (stringDictMsg) Ints() map[string]int64 {
	panic("unexpected Ints method call on string dict message")
}
func (stringDictMsg) Floats() map[string]float64 {
	panic("unexpected Floats method call on string dict message")
}
func (msg stringDictMsg) Strings() map[string]string { return msg.v }
func (msg stringDictMsg) Len() int                   { return len(msg.v) }
func (msg stringDictMsg) String() string             { return mustJSON(msg) }

//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg stringDictMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.v)
}
func (msg stringDictMsg) Equal(other DictMsg) bool { return dictEqual(msg, other) }

//nolint:ireturn // Msg contract type.
func NewDictMsg(d map[string]Msg) Msg {
	return dictValueMsg{internalMsg: internalMsg{}, v: genericDictMsg{v: d}}
}

//nolint:ireturn // Msg contract type.
func NewDictBoolMsg(d map[string]bool) Msg {
	return dictValueMsg{internalMsg: internalMsg{}, v: boolDictMsg{v: d}}
}

//nolint:ireturn // Msg contract type.
func NewDictIntMsg(d map[string]int64) Msg {
	return dictValueMsg{internalMsg: internalMsg{}, v: intDictMsg{v: d}}
}

//nolint:ireturn // Msg contract type.
func NewDictFloatMsg(d map[string]float64) Msg {
	return dictValueMsg{internalMsg: internalMsg{}, v: floatDictMsg{v: d}}
}

//nolint:ireturn // Msg contract type.
func NewDictStringMsg(d map[string]string) Msg {
	return dictValueMsg{internalMsg: internalMsg{}, v: stringDictMsg{v: d}}
}

func dictEqual(left DictMsg, right DictMsg) bool {
	if left.Len() != right.Len() {
		return false
	}

	leftMsgs := asGenericDict(left)
	rightMsgs := asGenericDict(right)
	for key, leftVal := range leftMsgs {
		rightVal, ok := rightMsgs[key]
		if !ok || !leftVal.Equal(rightVal) {
			return false
		}
	}
	return true
}

func listEqualGeneric(left []Msg, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case genericListMsg:
		for i := range left {
			if !left[i].Equal(rightTyped.v[i]) {
				return false
			}
		}
	case boolListMsg:
		for i := range left {
			if !left[i].Equal(NewBoolMsg(rightTyped.v[i])) {
				return false
			}
		}
	case intListMsg:
		for i := range left {
			if !left[i].Equal(NewIntMsg(rightTyped.v[i])) {
				return false
			}
		}
	case floatListMsg:
		for i := range left {
			if !left[i].Equal(NewFloatMsg(rightTyped.v[i])) {
				return false
			}
		}
	case stringListMsg:
		for i := range left {
			if !left[i].Equal(NewStringMsg(rightTyped.v[i])) {
				return false
			}
		}
	default:
		panic("unexpected list implementation")
	}
	return true
}

func listEqualBool(left []bool, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case boolListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case genericListMsg:
		for i := range left {
			if !NewBoolMsg(left[i]).Equal(rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualInt(left []int64, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case intListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case genericListMsg:
		for i := range left {
			if !NewIntMsg(left[i]).Equal(rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualFloat(left []float64, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case floatListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case genericListMsg:
		for i := range left {
			if !NewFloatMsg(left[i]).Equal(rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func listEqualString(left []string, right ListMsg) bool {
	switch rightTyped := right.(type) {
	case stringListMsg:
		for i := range left {
			if left[i] != rightTyped.v[i] {
				return false
			}
		}
	case genericListMsg:
		for i := range left {
			if !NewStringMsg(left[i]).Equal(rightTyped.v[i]) {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func asGenericDict(dict DictMsg) map[string]Msg {
	switch typed := dict.(type) {
	case genericDictMsg:
		return typed.v
	case boolDictMsg:
		out := make(map[string]Msg, len(typed.v))
		for key, value := range typed.v {
			out[key] = NewBoolMsg(value)
		}
		return out
	case intDictMsg:
		out := make(map[string]Msg, len(typed.v))
		for key, value := range typed.v {
			out[key] = NewIntMsg(value)
		}
		return out
	case floatDictMsg:
		out := make(map[string]Msg, len(typed.v))
		for key, value := range typed.v {
			out[key] = NewFloatMsg(value)
		}
		return out
	case stringDictMsg:
		out := make(map[string]Msg, len(typed.v))
		for key, value := range typed.v {
			out[key] = NewStringMsg(value)
		}
		return out
	default:
		panic("unexpected dict implementation")
	}
}

func mustJSON(msg interface{ MarshalJSON() ([]byte, error) }) string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func AsListMsgs(list ListMsg) ([]Msg, bool) {
	typed, ok := list.(genericListMsg)
	return typed.v, ok
}

func AsListBools(list ListMsg) ([]bool, bool) {
	typed, ok := list.(boolListMsg)
	return typed.v, ok
}

func AsListInts(list ListMsg) ([]int64, bool) {
	typed, ok := list.(intListMsg)
	return typed.v, ok
}

func AsListFloats(list ListMsg) ([]float64, bool) {
	typed, ok := list.(floatListMsg)
	return typed.v, ok
}

func AsListStrings(list ListMsg) ([]string, bool) {
	typed, ok := list.(stringListMsg)
	return typed.v, ok
}

func AsDictMsgs(dict DictMsg) (map[string]Msg, bool) {
	typed, ok := dict.(genericDictMsg)
	return typed.v, ok
}

func AsDictBools(dict DictMsg) (map[string]bool, bool) {
	typed, ok := dict.(boolDictMsg)
	return typed.v, ok
}

func AsDictInts(dict DictMsg) (map[string]int64, bool) {
	typed, ok := dict.(intDictMsg)
	return typed.v, ok
}

func AsDictFloats(dict DictMsg) (map[string]float64, bool) {
	typed, ok := dict.(floatDictMsg)
	return typed.v, ok
}

func AsDictStrings(dict DictMsg) (map[string]string, bool) {
	typed, ok := dict.(stringDictMsg)
	return typed.v, ok
}

// --- STRUCT ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type StructMsg struct {
	internalMsg
	fields []StructField
}

func (msg StructMsg) Struct() StructMsg { return msg }

// get returns the value of a field by name.
// it panics if the field is not found.
// it uses linear scan to find the field.
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg StructMsg) Get(name string) Msg { //nolint:ireturn,lll // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	if field, ok := msg.get(name); ok {
		return field
	}
	panic(fmt.Sprintf("field %q not found", name))
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return addJSONSpaces(jsonData), nil
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
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type StructField struct {
	value Msg
	name  string
}

// newstructfield constructs a structfield with provided name and value.
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func NewStructField(name string, value Msg) StructField {
	return StructField{name: name, value: value}
}

// newstruct builds a struct message from a slice of structfield.
// underlying struct representation remains unchanged for now.
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func NewStructMsg(fields []StructField) StructMsg { return newStructMsg(fields) }

// --- UNION ---
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (msg UnionMsg) Data() Msg {
	return msg.data
}

func (msg UnionMsg) String() string {
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (msg UnionMsg) MarshalJSON() ([]byte, error) {
	if msg.data == nil {
		return fmt.Appendf(nil, `{ "tag": %q }`, msg.tag), nil
	}

	dataJSON, err := json.Marshal(msg.data)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}
	dataJSON = addJSONSpaces(dataJSON)

	return fmt.Appendf(nil, `{ "tag": %q, "data": %s }`, msg.tag, dataJSON), nil
}

func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := orderedPayload(msg).(UnionMsg)
	return unionMsg, ok
}

//nolint:ireturn // Msg is runtime contract type.
func orderedPayload(msg Msg) Msg {
	if ordered, ok := msg.(OrderedMsg); ok {
		return ordered.Msg
	}
	return msg
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
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

func addJSONSpaces(jsonData []byte) []byte {
	spaced := make([]byte, 0, len(jsonData))
	inString := false
	isEscaped := false

	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	for _, b := range jsonData {
		if inString {
			spaced = append(spaced, b)
			if isEscaped {
				isEscaped = false
				continue
			}
			if b == '\\' {
				isEscaped = true
				continue
			}
			if b == '"' {
				inString = false
			}
			continue
		}

		switch b {
		case '"':
			inString = true
			spaced = append(spaced, b)
		case ':':
			spaced = append(spaced, ':', ' ')
		case ',':
			spaced = append(spaced, ',', ' ')
		default:
			spaced = append(spaced, b)
		}
	}

	return spaced
}
