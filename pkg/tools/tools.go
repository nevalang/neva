package tools

import (
	"errors"
	"fmt"
)

var ErrNil = errors.New("nil arg")

// NilPanic panics if at least one argument is nil
func NilPanic(args ...interface{}) {
	for i, v := range args {
		if v == nil {
			panic(
				fmt.Errorf("%w: #%d", ErrNil, i), //nolint:err113
			)
		}
	}
}
