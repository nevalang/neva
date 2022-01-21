package translator

import (
	"github.com/emil14/neva/internal/old/compiler/program"
	runtime "github.com/emil14/neva/internal/old/runtime/program"
)

type Coder interface {
	Code(runtime.Program) ([]byte, error)
}

type Translator struct {
	coder Coder
}

func (t Translator) Translate(prog program.Program) ([]byte, error) {
	return nil, nil
}

func New() Translator {
	return Translator{}
}
