package messages

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

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

func NewStringMsg(s string) StringMsg {
	return StringMsg{
		internalMsg: internalMsg{},
		v:           s,
	}
}

func equalString(left StringMsg, right Msg) bool {
	return equalStringValue(left.v, right)
}

func equalStringValue(left string, right Msg) bool {
	rightTyped, ok := right.(StringMsg)
	return ok && left == rightTyped.v
}

// StringConcat returns the concatenation of left and right.
func StringConcat(left, right Msg) StringMsg { return NewStringMsg(left.Str() + right.Str()) }

// StringIsGreater reports whether left is lexicographically greater than right.
func StringIsGreater(left, right Msg) BoolMsg { return NewBoolMsg(left.Str() > right.Str()) }

// StringIsLesser reports whether left is lexicographically less than right.
func StringIsLesser(left, right Msg) BoolMsg { return NewBoolMsg(left.Str() < right.Str()) }

// StringJoin joins the string values in list with separator. It preserves the
// typed string-list fast path and relies on the compiler invariant for boxed
// list<string> values.
func StringJoin(list ListMsg, separator string) string {
	builder := strings.Builder{}
	if values, ok := ListAsStrings(list); ok {
		for index := range values {
			if index > 0 {
				builder.WriteString(separator)
			}
			builder.WriteString(values[index])
		}
		return builder.String()
	}

	for index, value := range list.Untyped() {
		if index > 0 {
			builder.WriteString(separator)
		}
		builder.WriteString(value.Str())
	}
	return builder.String()
}

// StringAt returns the rune at index. Negative indexes count from the end.
// The boolean is false when index is outside the string bounds.
func StringAt(value Msg, index int64) (StringMsg, bool) {
	runes := []rune(value.Str())
	item, found := listAt(runes, index)
	if !found {
		return StringMsg{}, false
	}
	return NewStringMsg(string(item)), true
}

// StringSlice returns the normalized rune-indexed range of value.
func StringSlice(value Msg, from, to int64) StringMsg {
	runes := []rune(value.Str())
	start, end := listSliceBounds(from, to, int64(len(runes)))
	return NewStringMsg(string(runes[start:end]))
}

// StringSplit returns a typed list of the substrings of value separated by delimiter.
//
//nolint:ireturn // Msg is the value-layer contract.
func StringSplit(value Msg, delimiter Msg) Msg {
	return NewListStringMsg(strings.Split(value.Str(), delimiter.Str()))
}

// StringToLower returns the Unicode lowercase form of value.
func StringToLower(value Msg) StringMsg { return NewStringMsg(strings.ToLower(value.Str())) }

// StringToUpper returns the Unicode uppercase form of value.
func StringToUpper(value Msg) StringMsg { return NewStringMsg(strings.ToUpper(value.Str())) }

// StringFromBytes converts bytes to a string value.
func StringFromBytes(value Msg) StringMsg { return NewStringMsg(string(value.Bytes())) }

// BytesFromString converts a string to a bytes value.
func BytesFromString(value Msg) BytesMsg { return NewBytesMsg([]byte(value.Str())) }

// StringFromIntCodePoint converts an integer Unicode code point to a string.
// Invalid code points become the Unicode replacement character.
func StringFromIntCodePoint(value Msg) StringMsg {
	codePoint := value.Int()
	if codePoint < 0 || codePoint > unicode.MaxRune || (codePoint >= 0xD800 && codePoint <= 0xDFFF) {
		return NewStringMsg(string(unicode.ReplacementChar))
	}
	// #nosec G115 -- guarded by Unicode range checks above.
	return NewStringMsg(string(rune(codePoint)))
}

// StringFromInt formats value as a base-10 integer string.
func StringFromInt(value Msg) StringMsg { return NewStringMsg(strconv.FormatInt(value.Int(), 10)) }

// StringFromBool formats value as a boolean string.
func StringFromBool(value Msg) StringMsg { return NewStringMsg(strconv.FormatBool(value.Bool())) }

// StringFromFloat formats value with format, precision, and bit size.
func StringFromFloat(value Msg, format Msg, precision, bitSize Msg) StringMsg {
	formatByte := byte('g')
	if formatString := format.Str(); len(formatString) > 0 {
		formatByte = formatString[0]
	}
	return NewStringMsg(strconv.FormatFloat(value.Float(), formatByte, int(precision.Int()), int(bitSize.Int())))
}

// StringFromIntBase formats value using base.
func StringFromIntBase(value, base Msg) StringMsg {
	return NewStringMsg(strconv.FormatInt(value.Int(), int(base.Int())))
}

// BoolFromString parses value as a boolean.
func BoolFromString(value Msg) (BoolMsg, error) {
	parsed, err := strconv.ParseBool(value.Str())
	return NewBoolMsg(parsed), err
}

// FloatFromString parses value as a floating-point value with bitSize.
func FloatFromString(value, bitSize Msg) (FloatMsg, error) {
	parsed, err := strconv.ParseFloat(value.Str(), int(bitSize.Int()))
	return NewFloatMsg(parsed), err
}

// IntFromString parses value as a base-10 platform-sized integer.
func IntFromString(value Msg) (IntMsg, error) {
	parsed, err := strconv.Atoi(value.Str())
	return NewIntMsg(int64(parsed)), err
}

// IntFromStringBase parses value with base and bitSize.
func IntFromStringBase(value, base, bitSize Msg) (IntMsg, error) {
	parsed, err := strconv.ParseInt(value.Str(), int(base.Int()), int(bitSize.Int()))
	return NewIntMsg(parsed), err
}

// StringRegexpSubmatch returns the full regular-expression match followed by
// its submatches as a typed string list.
//
//nolint:ireturn // Msg is the value-layer contract.
func StringRegexpSubmatch(pattern, value Msg) (Msg, error) {
	expression, err := regexp.Compile(pattern.Str())
	if err != nil {
		return nil, err //nolint:wrapcheck // Runtime function preserves the regexp package error text as its public error.
	}
	return NewListStringMsg(expression.FindStringSubmatch(fmt.Sprint(value))), nil
}

// StringFormat substitutes $<index> placeholders in template with args.
func StringFormat(template string, args []Msg) (string, error) {
	usedArgs := make(map[int]bool)
	var result strings.Builder
	result.Grow(len(template))

	for position := 0; position < len(template); {
		argumentIndex, nextPosition, hasPlaceholder, err := stringPlaceholderAt(template, position)
		if err != nil {
			return "", err
		}
		if hasPlaceholder {
			if argumentIndex >= len(args) {
				return "", fmt.Errorf("template refers to arg %d, but only %d args given", argumentIndex, len(args))
			}
			usedArgs[argumentIndex] = true
			fmt.Fprint(&result, args[argumentIndex])
			position = nextPosition
			continue
		}

		result.WriteByte(template[position])
		position++
	}

	if len(usedArgs) != len(args) {
		return "", fmt.Errorf("not all arguments are used in the template: got %v, used %v", len(args), len(usedArgs))
	}
	return result.String(), nil
}

func stringPlaceholderAt(template string, start int) (int, int, bool, error) {
	if start >= len(template) || template[start] != '$' {
		return 0, start, false, nil
	}

	digitStart := start + 1
	digitEnd := digitStart
	for digitEnd < len(template) && template[digitEnd] >= '0' && template[digitEnd] <= '9' {
		digitEnd++
	}
	if digitStart == digitEnd {
		return 0, start, false, nil
	}

	index, err := strconv.Atoi(template[digitStart:digitEnd])
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid placeholder %q: %w", template[digitStart:digitEnd], err)
	}
	return index, digitEnd, true, nil
}
