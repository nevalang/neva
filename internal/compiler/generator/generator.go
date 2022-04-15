package generator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/runtime"
)

type (
	Translator interface {
		Translate(compiler.Program) (runtime.Program, error)
	}
	Encoder interface {
		Encode(runtime.Program) ([]byte, error)
	}
)

var (
	ErrTranslator = errors.New("translator")
	ErrEncoder    = errors.New("encoder")
)

type Generator struct {
	translator Translator
	encoder    Encoder
}

func (g Generator) Generate(cprog compiler.Program) ([]byte, error) {
	rprog, err := g.translator.Translate(cprog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTranslator, err)
	}

	bb, err := g.encoder.Encode(rprog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEncoder, err)
	}

	return bb, nil
}
