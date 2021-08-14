package translator

import (
	cprog "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type translator struct{}

func (t translator) Translate(cprog.Program) rprog.Program {
	return rprog.Program{}
}

func MustNew() translator {
	return translator{}
}
