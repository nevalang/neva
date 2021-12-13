package compiler

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	runtime "github.com/emil14/neva/internal/runtime/program"
)

type (
	Compiler struct {
		storage    Storage
		parser     ModuleParser
		checker    Checker
		translator Translator
		coder      Coder
		opsIO      map[PkgComponentRef]program.IO
	}

	ModuleParser interface {
		Parse([]byte) (program.Module, error)
	}

	Checker interface {
		Check(program.Program) error
	}

	Translator interface {
		Translate(program.Program) (runtime.Program, error)
	}

	Coder interface {
		Code(runtime.Program) ([]byte, error)
	}

	Storage interface {
		Pkg(string, map[PkgComponentRef]struct{}) (Pkg, error)
	}

	Pkg struct {
		Exec      string
		Operators []PkgComponentRef
		Modules   map[PkgComponentRef][]byte
		Scope     map[string]PkgComponentRef
		Meta      PkgMeta
	}

	PkgComponentRef struct {
		Pkg, Name string
	}

	PkgMeta struct {
		CompilerVersion string
	}
)

var (
	ErrParsing    = errors.New("failed to parse module")
	ErrValidation = errors.New("module is invalid")
	ErrInternal   = errors.New("internal error")
)

func (c Compiler) Version() string {
	return "0.0.1"
}

type CompileResult struct {
	compiled    runtime.Program
	precompiled program.Program
}

func (c Compiler) Compile(descriptorPath string) (CompileResult, error) {
	pkg, err := c.storage.Pkg(descriptorPath, c.opsSet())
	if err != nil {
		return CompileResult{}, err
	}

	if v := c.Version(); v != pkg.Meta.CompilerVersion {
		return CompileResult{}, fmt.Errorf(
			"wrong compiler version: want %s, got %s", pkg.Meta.CompilerVersion, v,
		)
	}

	ops, err := c.pkgOps(pkg)
	if err != nil {
		return CompileResult{}, err
	}

	mods, err := c.pkgMods(pkg)
	if err != nil {
		return CompileResult{}, fmt.Errorf("%w: %v", ErrParsing, err)
	}

	scope, err := c.pkgScope(pkg, ops, mods)
	if err != nil {
		return CompileResult{}, err
	}

	precompiled := program.Program{
		Root:  pkg.Exec,
		Scope: scope,
	}

	if err := c.checker.Check(precompiled); err != nil {
		return CompileResult{}, err
	}

	compiled, err := c.translator.Translate(precompiled)
	if err != nil {
		return CompileResult{}, err
	}

	return CompileResult{compiled, precompiled}, nil
}

func (Compiler) pkgScope(
	pkg Pkg,
	ops map[PkgComponentRef]program.Operator,
	mods map[PkgComponentRef]program.Module,
) (map[string]program.Component, error) {
	scope := make(map[string]program.Component, len(pkg.Scope))

	for alias, ref := range pkg.Scope {
		op, ok := ops[ref]
		if ok {
			scope[alias] = program.Component{
				Type:     program.OperatorComponent,
				Operator: op,
			}
		}

		mod, ok := mods[ref]
		if !ok {
			return nil, fmt.Errorf("")
		}

		scope[alias] = program.Component{
			Type:   program.ModuleComponent,
			Module: mod,
		}
	}

	return scope, nil
}

func (c Compiler) pkgMods(pkg Pkg) (map[PkgComponentRef]program.Module, error) {
	mods := make(map[PkgComponentRef]program.Module, len(pkg.Modules))
	for ref, bb := range pkg.Modules {
		mod, err := c.parser.Parse(bb)
		if err != nil {
			return nil, err
		}
		mods[ref] = mod
	}
	return mods, nil
}

func (c Compiler) pkgOps(pkg Pkg) (map[PkgComponentRef]program.Operator, error) {
	ops := make(map[PkgComponentRef]program.Operator, len(pkg.Operators))

	for _, opRef := range pkg.Operators {
		io, ok := c.opsIO[opRef]
		if !ok {
			return nil, fmt.Errorf("operator not found %s", opRef)
		}

		ops[opRef] = program.Operator{IO: io}
	}

	return ops, nil
}

func (c Compiler) opsSet() map[PkgComponentRef]struct{} {
	opsSet := make(map[PkgComponentRef]struct{}, len(c.opsIO))
	for ref := range c.opsIO {
		opsSet[ref] = struct{}{}
	}
	return opsSet
}

func New(
	parser ModuleParser,
	checker Checker,
	translator Translator,
	coder Coder,
	store Storage,
	opsIO map[PkgComponentRef]program.IO,
) (Compiler, error) {
	if parser == nil || checker == nil || translator == nil || store == nil || opsIO == nil {
		return Compiler{}, fmt.Errorf("nil deps")
	}

	return Compiler{
		parser:     parser,
		checker:    checker,
		translator: translator,
		coder:      coder,
		opsIO:      opsIO,
		storage:    store,
	}, nil
}

func MustNew(
	parser ModuleParser,
	checker Checker,
	translator Translator,
	coder Coder,
	store Storage,
	opsIO map[PkgComponentRef]program.IO,
) Compiler {
	cmp, err := New(parser, checker, translator, coder, store, opsIO)
	if err != nil {
		panic(err)
	}

	return cmp
}
