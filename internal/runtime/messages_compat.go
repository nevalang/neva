package runtime

import "github.com/nevalang/neva/internal/runtime/messages"

// Message value aliases remain while runtime callers migrate to messages.
type (
	Msg         = messages.Msg
	BoolMsg     = messages.BoolMsg
	IntMsg      = messages.IntMsg
	FloatMsg    = messages.FloatMsg
	StringMsg   = messages.StringMsg
	BytesMsg    = messages.BytesMsg
	ListMsg     = messages.ListMsg
	DictMsg     = messages.DictMsg
	StructMsg   = messages.StructMsg
	StructField = messages.StructField
	UnionMsg    = messages.UnionMsg
)

func NewBoolMsg(value bool) BoolMsg { return messages.NewBoolMsg(value) }

func NewIntMsg(value int64) IntMsg { return messages.NewIntMsg(value) }

func NewFloatMsg(value float64) FloatMsg { return messages.NewFloatMsg(value) }

func NewStringMsg(value string) StringMsg { return messages.NewStringMsg(value) }

func NewBytesMsg(value []byte) BytesMsg { return messages.NewBytesMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewListMsg(value []Msg) Msg { return messages.NewListMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewListBoolMsg(value []bool) Msg { return messages.NewListBoolMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewListIntMsg(value []int64) Msg { return messages.NewListIntMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewListFloatMsg(value []float64) Msg { return messages.NewListFloatMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewListStringMsg(value []string) Msg { return messages.NewListStringMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewDictMsg(value map[string]Msg) Msg { return messages.NewDictMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewDictBoolMsg(value map[string]bool) Msg { return messages.NewDictBoolMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewDictIntMsg(value map[string]int64) Msg { return messages.NewDictIntMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewDictFloatMsg(value map[string]float64) Msg { return messages.NewDictFloatMsg(value) }

//nolint:ireturn // Compatibility API preserves the existing message interface return type.
func NewDictStringMsg(value map[string]string) Msg { return messages.NewDictStringMsg(value) }

func AsListUntyped(list ListMsg) ([]Msg, bool) { return messages.AsListUntyped(list) }

func AsListBools(list ListMsg) ([]bool, bool) { return messages.AsListBools(list) }

func AsListInts(list ListMsg) ([]int64, bool) { return messages.AsListInts(list) }

func AsListFloats(list ListMsg) ([]float64, bool) { return messages.AsListFloats(list) }

func AsListStrings(list ListMsg) ([]string, bool) { return messages.AsListStrings(list) }

func AsDictUntyped(dict DictMsg) (map[string]Msg, bool) { return messages.AsDictUntyped(dict) }

func AsDictBools(dict DictMsg) (map[string]bool, bool) { return messages.AsDictBools(dict) }

func AsDictInts(dict DictMsg) (map[string]int64, bool) { return messages.AsDictInts(dict) }

func AsDictFloats(dict DictMsg) (map[string]float64, bool) { return messages.AsDictFloats(dict) }

func AsDictStrings(dict DictMsg) (map[string]string, bool) { return messages.AsDictStrings(dict) }

func NewStructField(name string, value Msg) StructField {
	return messages.NewStructField(name, value)
}

func NewStructMsg(fields []StructField) StructMsg { return messages.NewStructMsg(fields) }

func NewUnionMsg(tag string, data Msg) UnionMsg { return messages.NewUnionMsg(tag, data) }

func Uint8Index(index int) uint8 { return messages.Uint8Index(index) }
