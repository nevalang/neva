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
		Parse([]byte) (program.Modules, error)
	}

	// Validator verifies that program is correct.
	Validator interface {
		Validate(program.Modules) error // todo validate program
	}

	// Translator creates runtime program representation.
	Translator interface {
		Translate(program.Program) (rprog.Program, error)
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
	operators  map[string]program.Operator
}

// Compile compiles source code to bytecode.
func (c Compiler) Compile(src []byte) ([]byte, error) {
	mod, err := c.parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsing, err)
	}

	if err := c.validator.Validate(mod); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	prog := program.Program{
		Root: "root",
		// TODO
		Components: map[string]program.Component{
			"root": mod,
			"*":    c.operators["*"],
		},
	}

	rprog, err := c.translator.Translate(prog)
	if err != nil {
		return nil, err
	}

	bb, err := c.coder.Code(rprog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return bb, nil
}

// New creates new Compiler. It returns error if at least one of the deps is nil.
func New(p Parser, v Validator, t Translator, c Coder, ops map[string]program.Operator) (Compiler, error) {
	if p == nil || v == nil || t == nil || c == nil || ops == nil {
		return Compiler{}, fmt.Errorf("%w: failed to build compiler", ErrInternal)
	}

	return Compiler{
		parser:     p,
		validator:  v,
		translator: t,
		coder:      c,
		operators:  ops,
	}, nil
}

// MustNew creates Compiler or panics.
func MustNew(p Parser, v Validator, t Translator, c Coder, ops map[string]program.Operator) Compiler {
	cmp, err := New(p, v, t, c, ops)
	if err != nil {
		panic(err)
	}
	return cmp
}
