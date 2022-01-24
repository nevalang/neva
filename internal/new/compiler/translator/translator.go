package translator

import (
	"github.com/emil14/neva/internal/new/compiler"
	"github.com/emil14/neva/internal/new/runtime"
)

type Coder interface {
	Encode(runtime.Program) ([]byte, error)
}

type Translator struct {
	coder Coder
}

func (t Translator) Translate(compiler.Program) ([]byte, error) {
	return nil, nil
}
