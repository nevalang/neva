package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
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
	MsgKindBytes
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

func (msg Msg) IsBytes() bool {
	return msg.kind == MsgKindBytes
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
	if msg.kind == MsgKindBool {
		return msg.bits == 1
	}
	panicUnexpectedBoolKind(msg.kind)
	return false
}

func (msg Msg) Int() int64 {
	if msg.kind == MsgKindInt {
		// #nosec G115 -- msg.bits stores int64 in two's complement form.
		return int64(msg.bits)
	}
	panicUnexpectedIntKind(msg.kind)
	return 0
}

func (msg Msg) Float() float64 {
	if msg.kind == MsgKindFloat {
		return math.Float64frombits(msg.bits)
	}
	panicUnexpectedFloatKind(msg.kind)
	return 0
}

func (msg Msg) Str() string {
	if msg.kind == MsgKindString {
		return msg.str
	}
	panicUnexpectedStrKind(msg.kind)
	return ""
}

func (msg Msg) Bytes() []byte {
	if msg.kind != MsgKindBytes {
		panicUnexpectedBytesKind(msg.kind)
	}
	value, ok := msg.val.([]byte)
	if !ok {
		panic("unexpected Bytes value type")
	}
	return value
}

func (msg Msg) List() []Msg {
	if msg.kind != MsgKindList {
		panicUnexpectedListKind(msg.kind)
	}
	list, ok := msg.val.([]Msg)
	if !ok {
		panic("unexpected List value type")
	}
	return list
}

func (msg Msg) Dict() map[string]Msg {
	if msg.kind != MsgKindDict {
		panicUnexpectedDictKind(msg.kind)
	}
	dict, ok := msg.val.(map[string]Msg)
	if !ok {
		panic("unexpected Dict value type")
	}
	return dict
}

func (msg Msg) Struct() StructMsg {
	if msg.kind != MsgKindStruct {
		panicUnexpectedStructKind(msg.kind)
	}
	structMsg, ok := msg.val.(StructMsg)
	if !ok {
		panic("unexpected Struct value type")
	}
	return structMsg
}

func (msg Msg) Union() UnionMsg {
	if msg.kind != MsgKindUnion {
		panicUnexpectedUnionKind(msg.kind)
	}
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
	case MsgKindBytes, MsgKindList, MsgKindDict, MsgKindStruct, MsgKindUnion:
		b, err := msg.MarshalJSON()
		if err != nil {
			panic(err)
		}
		return string(b)
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
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return json.Marshal(msg.Str())
	case MsgKindBytes:
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return json.Marshal(msg.Bytes())
	case MsgKindList:
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

//nolint:gocognit,gocyclo,cyclop,funlen // hot-path switch intentionally explicit for performance.
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
		return math.Float64frombits(msg.bits) == math.Float64frombits(other.bits)
	case MsgKindString:
		return msg.str == other.str
	case MsgKindBytes:
		left, leftOK := msg.val.([]byte)
		right, rightOK := other.val.([]byte)
		if !leftOK || !rightOK {
			panic("unexpected Bytes value type")
		}
		return bytes.Equal(left, right)
	case MsgKindList:
		left, leftOK := msg.val.([]Msg)
		right, rightOK := other.val.([]Msg)
		if !leftOK || !rightOK {
			panic("unexpected List value type")
		}
		return equalMsgLists(left, right)
	case MsgKindDict:
		left, leftOK := msg.val.(map[string]Msg)
		right, rightOK := other.val.(map[string]Msg)
		if !leftOK || !rightOK {
			panic("unexpected Dict value type")
		}
		return equalMsgDicts(left, right)
	case MsgKindStruct:
		left, leftOK := msg.val.(StructMsg)
		right, rightOK := other.val.(StructMsg)
		if !leftOK || !rightOK {
			panic("unexpected Struct value type")
		}
		return left.Equal(right)
	case MsgKindUnion:
		left, leftOK := msg.val.(UnionMsg)
		right, rightOK := other.val.(UnionMsg)
		if !leftOK || !rightOK {
			panic("unexpected Union value type")
		}
		return left.Equal(right)
	default:
		return false
	}
}

// Cold panic helpers keep hot accessors inline-friendly.
//
//go:noinline
func panicUnexpectedBoolKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Bool call on %s message", kind))
}

//go:noinline
func panicUnexpectedIntKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Int call on %s message", kind))
}

//go:noinline
func panicUnexpectedFloatKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Float call on %s message", kind))
}

//go:noinline
func panicUnexpectedStrKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Str call on %s message", kind))
}

//go:noinline
func panicUnexpectedBytesKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Bytes call on %s message", kind))
}

//go:noinline
func panicUnexpectedListKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected List call on %s message", kind))
}

//go:noinline
func panicUnexpectedDictKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Dict call on %s message", kind))
}

//go:noinline
func panicUnexpectedStructKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Struct call on %s message", kind))
}

//go:noinline
func panicUnexpectedUnionKind(kind MsgKind) {
	panic(fmt.Sprintf("unexpected Union call on %s message", kind))
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
	case MsgKindBytes:
		return "bytes"
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

// NewStringMsg wraps a string runtime value.
func NewStringMsg(s string) Msg {
	return Msg{kind: MsgKindString, str: s}
}

// NewBytesMsg wraps a bytes runtime value.
func NewBytesMsg(v []byte) Msg {
	return Msg{kind: MsgKindBytes, val: v}
}

// NewListMsg wraps a list runtime value.
func NewListMsg(v []Msg) Msg {
	return Msg{kind: MsgKindList, val: v}
}

// NewDictMsg wraps a dict runtime value.
func NewDictMsg(d map[string]Msg) Msg {
	return Msg{kind: MsgKindDict, val: d}
}

// StructMsg stores struct fields in runtime representation.
type StructMsg struct {
	fields []StructField
}

func (msg StructMsg) Struct() StructMsg { return msg }

// Get returns the value of a field by name.
// It panics if the field is not found.
// It uses linear scan to find the field.
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

	return marshalStructFieldsWithSpaces(fields)
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

// StructField is a helper to construct structs via runtime.newstruct api without exposing fields.
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

// NewStructField constructs a StructField with provided name and value.
func NewStructField(name string, value Msg) StructField {
	return StructField{name: name, value: value}
}

// NewStructMsg builds a struct message from a slice of StructField.
// Underlying struct representation remains unchanged for now.
func NewStructMsg(fields []StructField) Msg { return newStructMsg(fields).Msg() }

func NewStructValue(fields []StructField) StructMsg { return newStructMsg(fields) }

// UnionMsg represents tagged union runtime data.
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
	b, err := msg.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
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
		return fmt.Appendf(nil, `{ "tag": %q }`, msg.tag), nil
	}

	dataJSON, err := msg.data.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(nil, `{ "tag": %q, "data": %s }`, msg.tag, dataJSON), nil
}

// Equal implements strict equality for UnionMsg messages.
// If one union has data and another doesn't, it returns false.
// It returns false if tags are different.
// It returns false if data is different.
// Tags are compared as Go strings and data is compared recursively using Equal method.
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
	return UnionMsg{tag: tag, hasData: false}.Msg()
}

func NewUnionValue(tag string, data Msg) UnionMsg {
	return newUnionMsg(tag, data)
}

func NewUnionValueNoData(tag string) UnionMsg {
	return UnionMsg{tag: tag, hasData: false}
}

// --- OPERATIONS ---

// Match compares two messages and returns true if they match, false otherwise.
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

//nolint:gocognit // hot-path loop intentionally explicit for performance.
func equalMsgLists(left []Msg, right []Msg) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		leftItem := left[i]
		rightItem := right[i]
		if leftItem.kind != rightItem.kind {
			return false
		}
		//nolint:exhaustive // non-primitive kinds intentionally handled by default fallback.
		switch leftItem.kind {
		case MsgKindBool, MsgKindInt:
			if leftItem.bits != rightItem.bits {
				return false
			}
		case MsgKindFloat:
			if math.Float64frombits(leftItem.bits) != math.Float64frombits(rightItem.bits) {
				return false
			}
		case MsgKindString:
			if leftItem.str != rightItem.str {
				return false
			}
		default:
			if !leftItem.Equal(rightItem) {
				return false
			}
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

func marshalMapWithSpaces(msgMap map[string]Msg) ([]byte, error) {
	if len(msgMap) == 0 {
		return []byte("{}"), nil
	}

	keys := make([]string, 0, len(msgMap))
	for k := range msgMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]byte, 0, len(msgMap)*16+2)
	out = append(out, '{')
	for idx := range keys {
		if idx > 0 {
			out = append(out, ',', ' ')
		}

		keyJSON, err := json.Marshal(keys[idx])
		if err != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return nil, err
		}
		out = append(out, keyJSON...)
		out = append(out, ':', ' ')

		valueJSON, err := marshalNestedJSON(msgMap[keys[idx]])
		if err != nil {
			return nil, err
		}
		out = append(out, valueJSON...)
	}
	out = append(out, '}')

	return out, nil
}

func marshalStructFieldsWithSpaces(fields []StructField) ([]byte, error) {
	out := make([]byte, 0, len(fields)*16+2)
	out = append(out, '{')
	for idx := range fields {
		if idx > 0 {
			out = append(out, ',', ' ')
		}

		keyJSON, err := json.Marshal(fields[idx].name)
		if err != nil {
			//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			return nil, err
		}
		out = append(out, keyJSON...)
		out = append(out, ':', ' ')

		valueJSON, err := marshalNestedJSON(fields[idx].value)
		if err != nil {
			return nil, err
		}
		out = append(out, valueJSON...)
	}
	out = append(out, '}')

	return out, nil
}

func marshalNestedJSON(msg Msg) ([]byte, error) {
	switch msg.Kind() {
	case MsgKindInvalid, MsgKindBool, MsgKindInt, MsgKindFloat, MsgKindString, MsgKindBytes:
		return msg.MarshalJSON()
	case MsgKindList:
		return marshalNestedListWithSpaces(msg.List())
	case MsgKindDict:
		return marshalMapWithSpaces(msg.Dict())
	case MsgKindStruct:
		return marshalStructFieldsWithSpaces(msg.Struct().fields)
	case MsgKindUnion:
		return marshalNestedUnionCompact(msg.Union())
	default:
		panic(fmt.Sprintf("unexpected nested JSON marshal for unknown message kind: %d", msg.Kind()))
	}
}

func marshalNestedListWithSpaces(list []Msg) ([]byte, error) {
	if len(list) == 0 {
		return []byte("[]"), nil
	}

	var out []byte
	out = append(out, '[')
	for i := range list {
		if i > 0 {
			out = append(out, ',', ' ')
		}
		itemJSON, err := marshalNestedJSON(list[i])
		if err != nil {
			return nil, err
		}
		out = append(out, itemJSON...)
	}
	out = append(out, ']')

	return out, nil
}

func marshalNestedUnionCompact(msg UnionMsg) ([]byte, error) {
	if !msg.hasData {
		return fmt.Appendf(nil, `{"tag": %q}`, msg.tag), nil
	}

	dataJSON, err := marshalNestedJSON(msg.data)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, `{"tag": %q, "data": %s}`, msg.tag, dataJSON), nil
}
