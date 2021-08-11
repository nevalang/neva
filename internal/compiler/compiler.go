package compiler

import (
	"fmt"

	cprogram "github.com/emil14/stream/internal/compiler/program"
	rprogram "github.com/emil14/stream/internal/runtime/program"
)

type compiler struct {
	parser     Parser
	validator  Validator
	translator Translator
	coder      Coder
}

func (c compiler) Compile(src []byte) ([]byte, error) {
	mod, err := c.parser.Parse(src)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Validate(mod); err != nil {
		return nil, err
	}

	bb, err := c.coder.Code(c.translator.Translate(mod))
	if err != nil {
		return nil, err
	}

	return bb, nil
}

type Translator interface {
	Translate(cprogram.Module) rprogram.Program
}

type Coder interface {
	Code(rprogram.Program) ([]byte, error)
}

func New(p Parser, v Validator, t Translator, c Coder) (compiler, error) {
	if p == nil || v == nil || t == nil || c == nil {
		return compiler{}, fmt.Errorf("failed to build compiler")
	}

	return compiler{
		parser:     p,
		validator:  v,
		translator: t,
		coder:      c,
	}, nil
}
