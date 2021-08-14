package operators

import "github.com/emil14/neva/internal/runtime"

func New() map[string]runtime.Operator {
	return map[string]runtime.Operator{
		"*": Mul,
	}
}
