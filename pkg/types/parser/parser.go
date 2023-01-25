package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
)

var (
	ErrEmptyStr               = errors.New("empty string")
	ErrBraceExprLen           = errors.New("string that strats with '[' must be at least 4 characters long")
	ErrNoCloseBrace           = errors.New("string that strats with '[' must contain ']'")
	ErrArrSize                = errors.New("string betwee '[ and ']' must be integer")
	ErrArrType                = errors.New("string after '[' is not valid type expression")
	ErrMissingCurlyClose      = errors.New("non-union expression that strats with '{' must have '}' at the end")
	ErrUnionEl                = errors.New("failed to parse union element")
	ErrInvalidCurlyEl         = errors.New("elements inside curly braces must all be record pairs or enum elements")
	ErrTooMuchPartsForCurlyEl = errors.New("element inside curly braces cannot have more than 2 parts")
	ErrRecField               = errors.New("cannot parse record field")
	ErrMissingAngleClose      = errors.New("strings that has '<' must also has '>'")
	ErrEmptyAngleBrackets     = errors.New("string with '<>' must not contain arguments")
	ErrInstArg                = errors.New("could not parse inst argument")
)

// TODO make API to extend parser
func Parse(s string) (ts.Expr, error) { //nolint:funlen,gocognit
	s = strings.TrimSpace(s)

	rr := []rune(s)
	charCount := len(rr)
	if charCount == 0 {
		return ts.Expr{}, ErrEmptyStr
	}

	if unionEls := strings.Split(s, "|"); len(unionEls) > 1 { // union
		exprs := make([]ts.Expr, 0, len(unionEls))
		for i, el := range unionEls {
			expr, err := Parse(el)
			if err != nil {
				return ts.Expr{}, fmt.Errorf("%w: #%d, err %v", ErrUnionEl, i, err)
			}
			exprs = append(exprs, expr)
		}
		return h.Union(exprs...), nil
	}

	if strings.HasPrefix(s, "[") { // arr
		if charCount < 4 {
			return ts.Expr{}, fmt.Errorf("%w: got %d", ErrBraceExprLen, charCount)
		}

		closingIdx := strings.Index(s, "]")
		if closingIdx == -1 {
			return ts.Expr{}, ErrNoCloseBrace
		}

		betweenBraces := strings.TrimSpace(s[1:closingIdx])

		size, err := strconv.ParseInt(betweenBraces, 10, 64)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("%w: %v", ErrArrSize, err)
		}

		afterBraces := strings.TrimSpace(s[closingIdx+1:]) // we sure charCount>=4

		arrType, err := Parse(afterBraces)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
		}

		return h.Arr(int(size), arrType), nil
	}

	// record or enum
	if strings.HasPrefix(s, "{") { //nolint:nestif
		isRecord, isEnum := false, false

		closingIdx := strings.Index(s, "}")
		if closingIdx != charCount-1 {
			return ts.Expr{}, ErrMissingCurlyClose
		}

		betweenCurlyBraces := strings.TrimSpace(s[1:closingIdx])

		if betweenCurlyBraces == "" { // empty rec {}
			return h.Rec(nil), nil
		}

		els := strings.Split(betweenCurlyBraces, ",") // we know that's len(els) >= 1
		rec := make(map[string]ts.Expr, len(els))     // allocate memory for both record and enum
		enum := make([]string, 0, len(els))           // to complete computation in one cycle
		for _, el := range els {                      // els are record fields or enum elements
			parts := strings.Split(strings.TrimSpace(el), " ") // record field and its type or just enum element

			switch { // we don't handle len(parts) == 0 because we know there's someting between braces
			case len(parts) > 2:
				return ts.Expr{}, fmt.Errorf("%w: %v", ErrTooMuchPartsForCurlyEl, parts)
			case len(parts) == 2:
				isRecord = true
			case len(parts) == 1:
				isEnum = true
			} // at this point we have (isRecord || isEnum) == true

			if isRecord && isEnum {
				return ts.Expr{}, ErrInvalidCurlyEl
			}

			if isEnum {
				enum = append(enum, el) // enum's el is basically just a string, no need to parse anything there
				continue
			}

			expr, err := Parse(parts[1]) // no need to check isRecord because of how switch works, we know it's true
			if err != nil {
				return ts.Expr{}, fmt.Errorf("%w: %v", ErrRecField, err)
			}
			rec[parts[0]] = expr
		}

		if isEnum {
			return h.Enum(enum...), nil
		}

		return h.Rec(rec), nil // no need to check isRecord also
	} // at this point we know it's inst

	openIdx := strings.Index(s, "<")
	if openIdx == -1 {
		return h.Inst(s), nil
	}

	closeIdx := strings.LastIndex(s, ">")
	if closeIdx == -1 {
		return ts.Expr{}, ErrMissingAngleClose
	}

	betweenAngleBrackets := strings.TrimSpace(s[openIdx+1 : closeIdx])
	if betweenAngleBrackets == "" {
		return ts.Expr{}, ErrEmptyAngleBrackets
	}

	args := strings.Split(betweenAngleBrackets, ",")
	exprs := make([]ts.Expr, 0, len(args))
	for _, arg := range args {
		expr, err := Parse(arg)
		if err != nil {
			return ts.Expr{}, ErrInstArg
		}
		exprs = append(exprs, expr)
	}

	return h.Inst(s[0:openIdx], exprs...), nil
}
