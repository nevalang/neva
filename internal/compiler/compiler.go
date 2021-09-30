package compiler

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	runtime "github.com/emil14/neva/internal/runtime/program"
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
		// ParseProgram([]byte) (program.Program)
		ParseModule([]byte) (program.Module, error)
	}

	// Translator creates runtime program representation.
	Translator interface {
		Translate(program.Program) (runtime.Program, error)
	}

	// Validator verifies that program is correct.
	Validator interface {
		Validate(program.Module) error // todo validate program
	}

	// Coder creates bytecode for given runtime program.
	Coder interface {
		Code(runtime.Program) ([]byte, error)
	}

	// Storage is an abstraction that allowes retrieve modules
	Storage interface {
		Program(descriptorPath string) (map[string][]byte, string, error)
	}

	// RemoteModuleParams is data needed to get module
	RemoteModuleParams struct {
		repo                string
		major, minor, patch uint64
		descriptorPath      string
	}
)

// Compiler compiles source code into bytecode.
type Compiler struct {
	parser     Parser
	validator  Validator
	translator Translator
	coder      Coder
	storage    Storage
	operators  map[string]program.Operator
}

func (c Compiler) PreCompile(descriptorPath string) (runtime.Program, error) {
	return c.preCompile(descriptorPath)
}

func (c Compiler) Compile(descriptorPath string) ([]byte, error) {
	prog, err := c.preCompile(descriptorPath)
	if err != nil {
		return nil, err
	}

	bytecode, err := c.coder.Code(prog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return bytecode, nil
}

func (c Compiler) preCompile(descriptorPath string) (runtime.Program, error) {
	scope, root, err := c.storage.Program(descriptorPath)
	if err != nil {
		return runtime.Program{}, err
	}

	prog, err := c.compileComponent(scope, root, c.defaultScope(len(scope)))
	if err != nil {
		return runtime.Program{}, err
	}

	rprog, err := c.translator.Translate(prog)
	if err != nil {
		return runtime.Program{}, err
	}

	return rprog, nil
}

func (c Compiler) compileComponent(
	scope map[string][]byte,
	component string,
	cache map[string]program.Component,
) (program.Program, error) {
	var mod program.Module

	if _, ok := cache[component]; !ok {
		parsed, err := c.parser.ParseModule(scope["root"])
		if err != nil {
			return program.Program{}, fmt.Errorf("%w: %v", ErrParsing, err)
		}

		if err := c.validator.Validate(parsed); err != nil {
			return program.Program{}, fmt.Errorf("%w: %v", ErrValidation, err)
		}

		mod = parsed
	}

	for name, dep := range mod.Deps {
		_, err := c.compileComponent(scope, name, cache)
		if err != nil {
			return program.Program{}, nil
		}

		if err := dep.Compare(mod.Interface()); err != nil {
			return program.Program{}, err
		}
	}

	cache[component] = mod

	return program.Program{
		Root:       component,
		Components: cache,
	}, nil
}

func (c Compiler) defaultScope(padding int) map[string]program.Component {
	m := make(map[string]program.Component, len(c.operators)+padding)
	for i := range c.operators {
		m[c.operators[i].Name] = c.operators[i]
	}
	return m
}

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

func MustNew(p Parser, v Validator, t Translator, c Coder, ops map[string]program.Operator) Compiler {
	cmp, err := New(p, v, t, c, ops)
	if err != nil {
		panic(err)
	}
	return cmp
}
