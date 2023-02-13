package tools

import (
	"errors"
	"fmt"
)

var ErrNil = errors.New("nil arg")

func NilPanic(args ...interface{}) {
	for i, v := range args {
		if v == nil {
			panic(
				fmt.Errorf("%w: #%d", ErrNil, i), //nolint:err113
			)
		}
	}
}

func Must[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(err)
	}
	return v
}

func Pointer[T any](v T) *T {
	return &v
}
