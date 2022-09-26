package checker

import "github.com/emil14/neva/internal/runtime/src"

type Checker struct{}

func (c Checker) Check(src.Program) error {
	return nil
}
