package compiler

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

// Compiler errors.
var (
	ErrParsing    = errors.New("failed to parse module")
	ErrValidation = errors.New("module invalid")
	ErrInternal   = errors.New("internal error")
)

// Compiler deps.
type (
	// Parser parses source code into compiler program representation.
	Parser interface {
		Parse([]byte) (program.Module, error)
	}

	// Validator verifies that program is correct.
	Validator interface {
		Validate(program.Module) error // todo validate program
	}

	// Translator creates runtime program representation.
	Translator interface {
		Translate(program.Program) rprog.Program
	}

	// Coder creates bytecode for given runtime program.
	Coder interface {
		Code(rprog.Program) ([]byte, error)
	}
)

// Compiler compiles source code to bytecode.
type Compiler struct {
	parser     Parser
	validator  Validator
	translator Translator
	coder      Coder
}

// Compile compiles source code down to bytecode.
func (cmplr Compiler) Compile(src []byte) ([]byte, error) {
	mod, err := cmplr.parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsing, err)
	}

	if err := cmplr.validator.Validate(mod); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	prog := program.Program{
		Components: map[string]program.Component{
			"root": mod,
		},
	}

	bb, err := cmplr.coder.Code(cmplr.translator.Translate(prog))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return bb, nil
}

// New creates new Compiler. It returns error if at least one of the deps is nil.
func New(p Parser, v Validator, t Translator, c Coder) (Compiler, error) {
	if p == nil || v == nil || t == nil || c == nil {
		return Compiler{}, fmt.Errorf("%w: failed to build compiler", ErrInternal)
	}

	return Compiler{
		parser:     p,
		validator:  v,
		translator: t,
		coder:      c,
	}, nil
}

// MustNew creates Compiler or panics.
func MustNew(p Parser, v Validator, t Translator, c Coder) Compiler {
	cmp, err := New(p, v, t, c)
	if err != nil {
		panic(err)
	}
	return cmp
}
