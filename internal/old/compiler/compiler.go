package compiler

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/old/compiler/program"
	"github.com/emil14/neva/internal/pkg/utils"
)

type (
	Storage interface {
		Pkg(string) (Pkg, error)
	}

	Checker interface {
		Check(program.Program) error
	}

	Translator interface {
		Translate(program.Program) ([]byte, error)
	}

	Pkg struct {
		Root      string
		Scope     map[string]ScopeRef
		Operators map[OpRef]program.ComponentIO
		Modules   map[string][]byte
		Meta      Meta
	}

	ScopeRef struct {
		Type ScopeRefType
		Pkg  string
		Name string
	}

	OpRef struct {
		Pkg, Name string
	}

	Meta struct {
		CompilerVersion string
	}
)

type ScopeRefType uint8

const (
	ScopeRefOperator ScopeRefType = iota + 1
	ScopeRefModule
)

var (
	ErrParsing    = errors.New("failed to parse module")
	ErrValidation = errors.New("module is invalid")
	ErrInternal   = errors.New("internal error")
)

type Compiler struct {
	checker    Checker
	translator Translator
}

func (c Compiler) version() string {
	return "0.0.1"
}

func (c Compiler) Compile(src Pkg) ([]byte, error) {
	if v := c.version(); v != src.Meta.CompilerVersion {
		return nil, fmt.Errorf(
			"wrong compiler version: want %s, got %s", src.Meta.CompilerVersion, v,
		)
	}

	scope, err := c.pkgScope(src, nil)
	if err != nil {
		return nil, err
	}

	precompiled := program.Program{
		RootComponent: src.Root,
		Components:    scope,
	}

	if err := c.checker.Check(precompiled); err != nil {
		return nil, err
	}

	return c.translator.Translate(precompiled)
}

func (Compiler) pkgScope(
	pkg Pkg,
	mods map[string]program.Module,
) (map[string]program.Component, error) {
	scope := make(map[string]program.Component, len(pkg.Scope))

	for alias, ref := range pkg.Scope {
		if ref.Type == ScopeRefOperator {
			// opref := OpRef{
			// 	Pkg:  ref.Pkg,
			// 	Name: ref.Name,
			// }

			// op, ok := ops[opref]
			// if ok {
			// 	scope[alias] = program.Component{
			// 		Type:       program.OperatorComponent,
			// 		OperatorIO: op,
			// 	}
			// }
		}

		mod, ok := mods[ref.Name]
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

func New(
	checker Checker,
	translator Translator,
	store Storage,
) (Compiler, error) {
	if err := utils.NilArgs(); err != nil {
		return Compiler{}, err
	}

	return Compiler{
		checker:    checker,
		translator: translator,
		// storage:    store,
	}, nil
}

func MustNew(
	checker Checker,
	translator Translator,
	store Storage,
) Compiler {
	cmp, err := New(checker, translator, store)
	utils.MaybePanic(err)
	return cmp
}

func NewOperatorsIO() map[string]map[string]program.ComponentIO {
	return map[string]map[string]program.ComponentIO{
		"math": {
			"mul":       program.ComponentIO{},
			"remainder": program.ComponentIO{},
		},
		"logic": {
			"more":   program.ComponentIO{},
			"filter": program.ComponentIO{},
		},
	}
}
