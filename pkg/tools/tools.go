package tools

import "fmt"

func PanicWithNil(args ...interface{}) {
	for i, v := range args {
		if v == nil {
			panic(
				fmt.Errorf("nil arg #%d", i),
			)
		}
	}
}

func Must[T any](v T, err error) T {
	PanicMaybe(err)
	return v
}

func PanicMaybe(err error) {
	if err != nil {
		panic(err)
	}
}
