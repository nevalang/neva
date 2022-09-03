package checker

import "github.com/emil14/neva/internal/runtime"

type Checker struct{}

func (c Checker) Check(runtime.Program) error {
	return nil
}
