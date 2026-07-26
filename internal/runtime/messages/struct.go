package messages

import (
	"encoding/json"
	"fmt"
)

type StructMsg struct {
	internalMsg
	fields []StructField
}

func (msg StructMsg) Struct() StructMsg { return msg }

// StructGetField returns the value of name in value. It panics when the field is
// absent because that contradicts the compiler's static struct-type guarantee.
//
//nolint:ireturn // Msg is the value-layer contract.
func StructGetField(value StructMsg, name string) Msg {
	if field, ok := value.get(name); ok {
		return field
	}
	panic(fmt.Sprintf("field %q not found", name))
}

// Get returns the value of name in msg.
// It forwards to StructGetField for compatibility with existing internal consumers.
//
//nolint:ireturn // Compatibility method returns the message contract.
func (msg StructMsg) Get(name string) Msg { return StructGetField(msg, name) }

// StructGetPath returns the nested field selected by path.
//
//nolint:ireturn // Msg is the value-layer contract.
func StructGetPath(value Msg, path []string) Msg {
	for _, name := range path {
		value = StructGetField(value.Struct(), name)
	}
	return value
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

func equalStructValue(left StructMsg, right Msg) bool {
	rightTyped, ok := right.(StructMsg)
	return ok && equalStructs(left, rightTyped)
}

func equalStructs(left, right StructMsg) bool {
	if len(left.fields) != len(right.fields) {
		return false
	}
	for i := range left.fields {
		rightValue, ok := right.get(left.fields[i].name)
		if !ok || !Equal(left.fields[i].value, rightValue) {
			return false
		}
	}
	return true
}

//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
