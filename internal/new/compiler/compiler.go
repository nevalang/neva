package compiler

import (
	"errors"
	"fmt"
)

type (
	PkgManager interface {
		Pkg(string) (Pkg, error)
	}
	ModuleParser interface {
		Parse(map[string][]byte) (map[string]Module, error)
	}
	ProgramChecker interface {
		Check(Program) error
	}
	ProgramTranslator interface {
		Translate(Program) ([]byte, error)
	}
)

var (
	ErrPkgManager     = errors.New("package manager")
	ErrModParser      = errors.New("module parser")
	ErrProgChecker    = errors.New("program checker")
	ErrProgTranslator = errors.New("program translator")
)

type Compiler struct {
	packager   PkgManager
	parser     ModuleParser
	checker    ProgramChecker
	translator ProgramTranslator

	opsIO map[ComponentRef]IO // sure?
}

func (c Compiler) Compile(path string) ([]byte, error) {
	prog, err := c.PreCompile(path)
	if err != nil {
		return nil, fmt.Errorf("precompile: %w", err)
	}

	bb, err := c.translator.Translate(prog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProgTranslator, err)
	}

	return bb, nil
}

func (c Compiler) PreCompile(path string) (Program, error) {
	pkg, err := c.packager.Pkg(path)
	if err != nil {
		return Program{}, fmt.Errorf("%w: %v", ErrPkgManager, err)
	}

	mods, err := c.parser.Parse(pkg.Modules)
	if err != nil {
		return Program{}, fmt.Errorf("%w: %v", ErrModParser, err)
	}

	prog := Program{
		RootModule: pkg.RootComponent,
		Operators:  map[string]ComponentRef{}, // todo
		Modules:    mods,
	}

	if err := c.checker.Check(prog); err != nil {
		return Program{}, fmt.Errorf("%w: %v", ErrProgChecker, err)
	}

	return prog, nil
}
