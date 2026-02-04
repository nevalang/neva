package runtime

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
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

type MsgKind uint8

const (
	MsgKindInvalid MsgKind = iota
	MsgKindBool
	MsgKindInt
	MsgKindFloat
	MsgKindString
	MsgKindList
	MsgKindDict
	MsgKindStruct
	MsgKindUnion
)

// Msg is the runtime representation of a Neva message.
// It is a tagged union optimized for minimal allocations in hot paths.
type Msg struct {
	val  any
	str  string
	bits uint64
	kind MsgKind
}

func (msg Msg) Kind() MsgKind {
	return msg.kind
}

func (msg Msg) IsValid() bool {
	return msg.kind != MsgKindInvalid
}

func (msg Msg) IsBool() bool {
	return msg.kind == MsgKindBool
}

func (msg Msg) IsInt() bool {
	return msg.kind == MsgKindInt
}

func (msg Msg) IsFloat() bool {
	return msg.kind == MsgKindFloat
}

func (msg Msg) IsString() bool {
	return msg.kind == MsgKindString
}

func (msg Msg) IsList() bool {
	return msg.kind == MsgKindList
}

func (msg Msg) IsDict() bool {
	return msg.kind == MsgKindDict
}

func (msg Msg) IsStruct() bool {
	return msg.kind == MsgKindStruct
}

func (msg Msg) IsUnion() bool {
	return msg.kind == MsgKindUnion
}

func (msg Msg) Bool() bool {
	msg.mustKind(MsgKindBool, "Bool")
	return msg.bits == 1
}

func (msg Msg) Int() int64 {
	msg.mustKind(MsgKindInt, "Int")
	// #nosec G115 -- msg.bits stores int64 in two's complement form.
	return int64(msg.bits)
}

func (msg Msg) Float() float64 {
	msg.mustKind(MsgKindFloat, "Float")
	return math.Float64frombits(msg.bits)
}

func (msg Msg) Str() string {
	msg.mustKind(MsgKindString, "Str")
	return msg.str
}

func (msg Msg) List() []Msg {
	msg.mustKind(MsgKindList, "List")
	list, ok := msg.val.([]Msg)
	if !ok {
		panic("unexpected List value type")
	}
	return list
}

func (msg Msg) Dict() map[string]Msg {
	msg.mustKind(MsgKindDict, "Dict")
	dict, ok := msg.val.(map[string]Msg)
	if !ok {
		panic("unexpected Dict value type")
	}
	return dict
}

func (msg Msg) Struct() StructMsg {
	msg.mustKind(MsgKindStruct, "Struct")
	structMsg, ok := msg.val.(StructMsg)
	if !ok {
		panic("unexpected Struct value type")
	}
	return structMsg
}

func (msg Msg) Union() UnionMsg {
	msg.mustKind(MsgKindUnion, "Union")
	unionMsg, ok := msg.val.(UnionMsg)
	if !ok {
		panic("unexpected Union value type")
	}
	return unionMsg
}

func (msg Msg) String() string {
	switch msg.kind {
	case MsgKindInvalid:
		panic("unexpected String call on invalid message")
	case MsgKindBool:
		return strconv.FormatBool(msg.Bool())
	case MsgKindInt:
		return strconv.FormatInt(msg.Int(), 10)
	case MsgKindFloat:
		return strconv.FormatFloat(msg.Float(), 'g', -1, 64)
	case MsgKindString:
		return msg.Str()
	case MsgKindList, MsgKindDict, MsgKindStruct:
		b, err := msg.MarshalJSON()
		if err != nil {
			panic(err)
		}
		return string(b)
	case MsgKindUnion:
		return msg.Union().String()
	default:
		panic(fmt.Sprintf("unexpected String call on unknown message kind: %d", msg.kind))
	}
}

func (msg Msg) MarshalJSON() ([]byte, error) {
	switch msg.kind {
	case MsgKindInvalid:
		panic("unexpected MarshalJSON call on invalid message")
	case MsgKindBool:
		return []byte(strconv.FormatBool(msg.Bool())), nil
	case MsgKindInt:
		return []byte(strconv.FormatInt(msg.Int(), 10)), nil
	case MsgKindFloat:
		return []byte(strconv.FormatFloat(msg.Float(), 'g', -1, 64)), nil
	case MsgKindString:
		return json.Marshal(msg.Str())
	case MsgKindList:
		return json.Marshal(msg.List())
	case MsgKindDict:
		return marshalMapWithSpaces(msg.Dict())
	case MsgKindStruct:
		return msg.Struct().MarshalJSON()
	case MsgKindUnion:
		return msg.Union().MarshalJSON()
	default:
		panic(fmt.Sprintf("unexpected MarshalJSON call on unknown message kind: %d", msg.kind))
	}
}

func (msg Msg) Equal(other Msg) bool {
	if msg.kind != other.kind {
		return false
	}
	switch msg.kind {
	case MsgKindInvalid:
		panic("unexpected Equal call on invalid message")
	case MsgKindBool, MsgKindInt:
		return msg.bits == other.bits
	case MsgKindFloat:
		return msg.Float() == other.Float()
	case MsgKindString:
		return msg.str == other.str
	case MsgKindList:
		return equalMsgLists(msg.List(), other.List())
	case MsgKindDict:
		return equalMsgDicts(msg.Dict(), other.Dict())
	case MsgKindStruct:
		return msg.Struct().Equal(other.Struct())
	case MsgKindUnion:
		return msg.Union().Equal(other.Union())
	default:
		return false
	}
}

func (msg Msg) mustKind(kind MsgKind, method string) {
	if msg.kind != kind {
		panic(fmt.Sprintf("unexpected %s call on %s message", method, msg.kind))
	}
}

func (kind MsgKind) String() string {
	switch kind {
	case MsgKindInvalid:
		return "invalid"
	case MsgKindBool:
		return "bool"
	case MsgKindInt:
		return "int"
	case MsgKindFloat:
		return "float"
	case MsgKindString:
		return "string"
	case MsgKindList:
		return "list"
	case MsgKindDict:
		return "dict"
	case MsgKindStruct:
		return "struct"
	case MsgKindUnion:
		return "union"
	default:
		return "unknown"
	}
}

// Bool

func NewBoolMsg(b bool) Msg {
	var bits uint64
	if b {
		bits = 1
	}
	return Msg{kind: MsgKindBool, bits: bits}
}

// Int

func NewIntMsg(n int64) Msg {
	// #nosec G115 -- store int64 bit pattern in uint64 for compact storage.
	return Msg{kind: MsgKindInt, bits: uint64(n)}
}

// Float

func NewFloatMsg(n float64) Msg {
	return Msg{kind: MsgKindFloat, bits: math.Float64bits(n)}
}

// --- STRING ---
func NewStringMsg(s string) Msg {
	return Msg{kind: MsgKindString, str: s}
}

// --- LIST ---
func NewListMsg(v []Msg) Msg {
	return Msg{kind: MsgKindList, val: v}
}

// --- DICT ---
func NewDictMsg(d map[string]Msg) Msg {
	return Msg{kind: MsgKindDict, val: d}
}

// --- STRUCT ---
type StructMsg struct {
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
	return Msg{}, false
}

func (msg StructMsg) Fields() []StructField {
	return msg.fields
}

func (msg StructMsg) MarshalJSON() ([]byte, error) {
	if len(msg.fields) == 0 {
		return []byte("{}"), nil
	}

	fields := make([]StructField, len(msg.fields))
	copy(fields, msg.fields)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].name < fields[j].name
	})

	var bb strings.Builder
	bb.Grow(2 + len(msg.fields)*8)
	bb.WriteByte('{')

	for i := range fields {
		if i > 0 {
			bb.WriteByte(',')
		}

		nameJSON, err := json.Marshal(fields[i].name)
		if err != nil {
			return nil, err
		}

		valueJSON, err := json.Marshal(fields[i].value)
		if err != nil {
			return nil, err
		}

		bb.Write(nameJSON)
		bb.WriteByte(':')
		bb.Write(valueJSON)
	}

	bb.WriteByte('}')

	jsonString := bb.String()
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
func (msg StructMsg) Equal(other StructMsg) bool {
	if len(msg.fields) != len(other.fields) {
		return false
	}
	for i := range msg.fields {
		otherField, ok := other.get(msg.fields[i].name)
		if !ok {
			return false
		}
		if !msg.fields[i].value.Equal(otherField) {
			return false
		}
	}
	return true
}

func (msg StructMsg) Msg() Msg {
	return Msg{kind: MsgKindStruct, val: msg}
}

func newStructMsg(fields []StructField) StructMsg {
	if len(fields) == 0 {
		return StructMsg{fields: nil}
	}
	copied := make([]StructField, len(fields))
	copy(copied, fields)
	return StructMsg{
		fields: copied,
	}
}

// structfield is a helper to construct structs via runtime.newstruct api without exposing fields.
type StructField struct {
	name  string
	value Msg
}

func (field StructField) Name() string {
	return field.name
}

func (field StructField) Value() Msg {
	return field.value
}

// newstructfield constructs a structfield with provided name and value.
func NewStructField(name string, value Msg) StructField {
	return StructField{name: name, value: value}
}

// newstruct builds a struct message from a slice of structfield.
// underlying struct representation remains unchanged for now.
func NewStructMsg(fields []StructField) Msg { return newStructMsg(fields).Msg() }

func NewStructValue(fields []StructField) StructMsg { return newStructMsg(fields) }

// --- UNION ---
type UnionMsg struct {
	tag     string
	data    Msg
	hasData bool
}

func (msg UnionMsg) Union() UnionMsg {
	return msg
}

func (msg UnionMsg) Tag() string {
	return msg.tag
}

func (msg UnionMsg) HasData() bool {
	return msg.hasData
}

func (msg UnionMsg) Data() Msg {
	return msg.data
}

func (msg UnionMsg) String() string {
	if !msg.hasData {
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

func (msg UnionMsg) MarshalJSON() ([]byte, error) {
	if !msg.hasData {
		return []byte(fmt.Sprintf(`{ "tag": %q }`, msg.tag)), nil
	}
	payload, err := json.Marshal(msg.data)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{ "tag": %q, "data": %s }`, msg.tag, payload)), nil
}

// Equal implements strict equality for UnionMsg messages.
// If one union has data and another doesn't, it returns false.
// It returns false if tags are different.
// It returns false if data is different.
// Tags are compared as Go strings and data is compared recursevely using Equal method.
func (msg UnionMsg) Equal(other UnionMsg) bool {
	if msg.hasData != other.hasData {
		return false
	}
	if msg.tag != other.tag {
		return false
	}
	if !msg.hasData {
		return true
	}
	return msg.data.Equal(other.data)
}

func (msg UnionMsg) Msg() Msg {
	return Msg{kind: MsgKindUnion, val: msg}
}

func newUnionMsg(tag string, data Msg) UnionMsg {
	return UnionMsg{
		tag:     tag,
		data:    data,
		hasData: data.IsValid(),
	}
}

func NewUnionMsg(tag string, data Msg) Msg {
	return newUnionMsg(tag, data).Msg()
}

func NewUnionMsgNoData(tag string) Msg {
	return UnionMsg{tag: tag}.Msg()
}

func NewUnionValue(tag string, data Msg) UnionMsg {
	return newUnionMsg(tag, data)
}

func NewUnionValueNoData(tag string) UnionMsg {
	return UnionMsg{tag: tag}
}

// --- OPERATIONS ---

// Match compares two messages and return true if they matches and false otherwise.
// Unlike Equal it compares only some aspects of the messages.
func Match(msg Msg, pattern Msg) bool {
	if !msg.IsUnion() || !pattern.IsUnion() {
		return msg.Equal(pattern)
	}

	msgUnion := msg.Union()
	patternUnion := pattern.Union()

	if msgUnion.tag != patternUnion.tag {
		return false
	}

	if !patternUnion.hasData {
		return true
	}

	if !msgUnion.hasData {
		return false
	}

	// by this time we know
	// both msg and pattern are union messages
	// they both have the same tags and some data inside
	// so we apply strict equality to the data they wrap
	// maybe in the future we'll consider recursive matching, we'll see
	return msgUnion.data.Equal(patternUnion.data)
}

func equalMsgLists(left []Msg, right []Msg) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if !left[i].Equal(right[i]) {
			return false
		}
	}
	return true
}

// equalMsgDicts is needed because maps.Equal won't work for non-comparable values.
func equalMsgDicts(left map[string]Msg, right map[string]Msg) bool {
	if len(left) != len(right) {
		return false
	}
	for k, v := range left {
		other, ok := right[k]
		if !ok || !v.Equal(other) {
			return false
		}
	}
	return true
}

func marshalMapWithSpaces(m map[string]Msg) ([]byte, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, ":", ": ")
	jsonString = strings.ReplaceAll(jsonString, ",", ", ")

	return []byte(jsonString), nil
}
