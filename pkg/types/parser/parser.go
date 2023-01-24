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
	ErrEmptyStr          = errors.New("empty string")
	ErrBraceExprLen      = errors.New("string that strats with '[' must be at least 4 characters long")
	ErrNoCloseBrace      = errors.New("string that strats with '[' must contain ']'")
	ErrArrSize           = errors.New("string betwee '[ and ']' must be integer")
	ErrArrType           = errors.New("string after '[' is not valid type expression")
	ErrMissingCurlyClose = errors.New("non-union expression that strats with '{' must have '}' at the end")
	ErrUnionEl           = errors.New("failed to parse union element")
	// ErrMissingCommas = errors.New("string must contain parts separated by commas")
)

func Parse(s string) (ts.Expr, error) { //nolint:funlen
	s = strings.TrimSpace(s)

	rr := []rune(s)
	charCount := len(rr)
	if charCount == 0 {
		return ts.Expr{}, ErrEmptyStr
	}

	if unionEls := strings.Split(s, "|"); len(unionEls) != 0 { // union
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

		afterBraces := strings.TrimSpace(s[closingIdx:])

		arrType, err := Parse(afterBraces)
		if err != nil {
			return ts.Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
		}

		return h.Arr(int(size), arrType), nil
	}

	if strings.HasPrefix(s, "{") { // record or enum
		isRecord := false
		isEnum := false

		closingIdx := strings.Index(s, "}")
		if closingIdx != charCount-1 {
			return ts.Expr{}, ErrMissingCurlyClose
		}

		betweenCurlyBraces := strings.TrimSpace(s[1:closingIdx])

		els := strings.Split(betweenCurlyBraces, ",")
		for _, parts := range els {
			if isRecord && isEnum {
				panic("")
			}

			parts := strings.Split(strings.TrimSpace(parts), " ")
			if len(parts) > 2 {
				panic("") // too many for rec and enum
			}
			if len(parts) == 2 {
				isRecord = true
			} else {
				isEnum = true
			}

			// TODO
		}
	}

	// TODO enums, records, insts

	return ts.Expr{}, nil
}
