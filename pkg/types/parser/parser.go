package parser

import (
	"errors"
	"strconv"
	"strings"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
)

var (
	ErrEmptyStr     = errors.New("empty string")
	ErrBraceExprLen = errors.New("string that strats with '[' must be at least 4 characters long")
	ErrNoCloseBrace = errors.New("string that strats with '[' contain ']'")
	ErrArrSize      = errors.New("string betwee '[ and ']' must be integer")
	ErrArrType      = errors.New("string after '[' is not valid type expression")
)

func Parse(s string) (ts.Expr, error) {
	rr := []rune(s)
	l := len(rr)
	if l == 0 {
		return ts.Expr{}, ErrEmptyStr
	}

	unionEls := strings.Split(s, "|")
	if len(unionEls) != 0 {
		els := make([]ts.Expr, 0, len(unionEls))
		for _, el := range unionEls {
			expr, err := Parse(el)
			if err != nil {
				return ts.Expr{}, err
			}
			els = append(els, expr)
		}
		return h.Union(els...), nil
	}

	if strings.HasPrefix(s, "[") {
		if l < 4 {
			return ts.Expr{}, ErrBraceExprLen
		}

		closing := strings.Index(s, "]")
		if closing == -1 {
			return ts.Expr{}, ErrNoCloseBrace
		}

		sizeint, err := strconv.ParseInt(s[1:closing], 10, 64)
		if err != nil {
			return ts.Expr{}, ErrArrSize
		}

		v, err := Parse(s[closing:])
		if err != nil {
			return ts.Expr{}, ErrArrType
		}

		return h.Arr(int(sizeint), v), nil
	}

	// TODO enums, records, insts

	return ts.Expr{}, nil
}
