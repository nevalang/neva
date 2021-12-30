package utils

import "fmt"

func MaybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func NilArgs(vv ...interface{}) error {
	for i, v := range vv {
		if v == nil {
			return fmt.Errorf("nil arg #%d", i)
		}
	}
	return nil
}
