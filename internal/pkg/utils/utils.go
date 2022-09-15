package utils

import "fmt"

// ErrFatal panics if err != nil
func ErrFatal(err error) {
	if err != nil {
		panic(err)
	}
}

// NilArgs returns non nil error if at least 1 arg is nil
func NilArgs(args ...interface{}) error {
	for i, v := range args {
		if v == nil {
			return fmt.Errorf("nil arg #%d", i)
		}
	}
	return nil
}

// PanicOnNil is NilArgs wrapped by ErrFatal
func PanicOnNil(args ...interface{}) {
	ErrFatal(
		NilArgs(args),
	)
}
