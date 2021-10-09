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
	ErrValidation = errors.New("module is invalid")
	ErrInternal   = errors.New("internal error")
)

// Compiler deps.
type (
	// Parser parses source code into compiler program representation.
	Parser interface {
		ParseModule([]byte) (program.Module, error)
		Program(program.Program) ([]byte, error)
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

	// Storage is an abstraction that allowes retrieve packages.
	Storage interface {
		Pkg(string) (Pkg, error)
	}

	// Pkg describes package.
	Pkg struct {
		Root    string
		Modules map[string][]byte
	}
)

// Compiler compiles source code into bytecode.
type Compiler struct {
	storage    Storage
	parser     Parser
	validator  Validator
	translator Translator
	coder      Coder
	operators  map[string]program.Operator
}

func (c Compiler) PreCompile(pkgDescriptorPath string) (runtime.Program, program.Program, error) {
	return c.preCompile(pkgDescriptorPath)
}

// func (c Compiler) Compile(pkgDescriptorPath string) ([]byte, error) {
// 	prog, err := c.preCompile(pkgDescriptorPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bytecode, err := c.coder.Code(prog)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
// 	}

// 	return bytecode, nil
// }

func (c Compiler) preCompile(pkgDescriptorPath string) (runtime.Program, program.Program, error) {
	pkg, err := c.storage.Pkg(pkgDescriptorPath)
	if err != nil {
		return runtime.Program{}, program.Program{}, err
	}

	scope := c.defaultScope(len(pkg.Modules))
	for k, v := range pkg.Modules {
		mod, err := c.compileModule(v)
		if err != nil {
			return runtime.Program{}, program.Program{}, err
		}
		scope[k] = mod
	}

	prog, err := c.buildProgram(scope, pkg.Root)
	if err != nil {
		return runtime.Program{}, program.Program{}, err
	}

	rprog, err := c.translator.Translate(prog)
	if err != nil {
		return runtime.Program{}, program.Program{}, err
	}

	return rprog, prog, nil
}

// TODO
func (c Compiler) buildProgram(scope map[string]program.Component, root string) (program.Program, error) {
	return program.Program{
		Root:       root,
		Components: scope,
	}, nil
}

func (c Compiler) compileModule(mod []byte) (program.Module, error) {
	parsed, err := c.parser.ParseModule(mod)
	if err != nil {
		return program.Module{}, fmt.Errorf("%w: %v", ErrParsing, err)
	}

	if err := c.validator.Validate(parsed); err != nil {
		return program.Module{}, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	return parsed, nil
}

func (c Compiler) compileProgram(
	modules map[string][]byte,
	root string,
	scope map[string]program.Component,
) (program.Program, error) {
	var mod program.Module

	if _, ok := scope[root]; !ok {
		parsed, err := c.parser.ParseModule(modules[root])
		if err != nil {
			return program.Program{}, fmt.Errorf("%w: %v", ErrParsing, err)
		}

		if err := c.validator.Validate(parsed); err != nil {
			return program.Program{}, fmt.Errorf("%w: %v", ErrValidation, err)
		}

		mod = parsed
	}

	for name, dep := range mod.Deps {
		prog, err := c.compileProgram(modules, name, scope)
		if err != nil {
			return program.Program{}, err
		}

		subroot, ok := prog.Components[prog.Root]
		if !ok {
			return program.Program{}, fmt.Errorf("%w", ErrInternal)
		}

		if err := subroot.Interface().Compare(dep); err != nil {
			return program.Program{}, err
		}
	}

	scope[root] = mod

	return program.Program{
		Root:       root,
		Components: scope,
	}, nil
}

// defaultScope returns scope with operators and free space for modules.
func (c Compiler) defaultScope(padding int) map[string]program.Component {
	m := make(map[string]program.Component, len(c.operators)+padding)
	for i := range c.operators {
		m[c.operators[i].Name] = c.operators[i]
	}
	return m
}

func New(p Parser, v Validator, t Translator, c Coder, s Storage, ops map[string]program.Operator) (Compiler, error) {
	if p == nil || v == nil || t == nil || c == nil || s == nil || ops == nil {
		return Compiler{}, fmt.Errorf("%w: failed to build compiler", ErrInternal)
	}

	return Compiler{
		parser:     p,
		validator:  v,
		translator: t,
		coder:      c,
		storage:    s,
		operators:  ops,
	}, nil
}

func MustNew(p Parser, v Validator, t Translator, c Coder, s Storage, ops map[string]program.Operator) Compiler {
	cmp, err := New(p, v, t, c, s, ops)
	if err != nil {
		panic(err)
	}
	return cmp
}
