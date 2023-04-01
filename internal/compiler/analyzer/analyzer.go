package analyzer

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/nevalang/internal/compiler"
	"github.com/nevalang/nevalang/internal/compiler/helper"
	"github.com/nevalang/nevalang/pkg/tools"
	ts "github.com/nevalang/nevalang/pkg/types"
)

var (
	ErrRootPkgMissing       = errors.New("program must have root pkg")
	ErrMainPkgNotFound      = errors.New("root pkg not found")
	ErrMainPkgNotExecutable = errors.New("main pkg must be executable (have root component)")
	ErrPkg                  = errors.New("analyze package")
	ErrTypeParams           = errors.New("interface type parameters")
	ErrIO                   = errors.New("io")
	ErrValidator            = errors.New("type validator")
	ErrResolver             = errors.New("type expression  resolver")
)

var h helper.Helper

// TODO make sure there's no two or more connections with the same sender
type Analyzer struct {
	resolver  TypeExprResolver
	checker   SubtypeChecker
	validator TypeValidator
	native    map[compiler.EntityRef]struct{} // set of components that must not have implementation
}

type (
	TypeExprResolver interface {
		Resolve(ts.Expr, ts.Scope) (ts.Expr, error)
	}
	SubtypeChecker interface {
		// TODO try to hide terminator params
		Check(ts.Expr, ts.Expr, ts.TerminatorParams) error
	}
	TypeValidator interface {
		ValidateParams([]ts.Param) error
	}
)

func (a Analyzer) Analyze(ctx context.Context, prog compiler.Program) (compiler.Program, error) {
	mainPkg, ok := prog.Pkgs["main"]
	if !ok {
		return compiler.Program{}, ErrMainPkgNotFound
	}

	if err := a.analyzeMainPkg(mainPkg, prog.Pkgs); err != nil {
		return compiler.Program{}, errors.Join(ErrExecutablePkg, err)
	}

	resolvedPkgs := make(map[string]compiler.Pkg, len(prog.Pkgs))
	for pkgName := range prog.Pkgs {
		resolvedPkg, err := a.analyzePkg(pkgName, prog.Pkgs)
		if err != nil {
			return compiler.Program{}, fmt.Errorf("%w: found in %v", errors.Join(ErrPkg, err), pkgName)
		}
		resolvedPkgs[pkgName] = resolvedPkg
	}

	return compiler.Program{
		Pkgs: resolvedPkgs,
	}, nil
}

// analyzeTypeParameters validates type parameters and resolves their constraints.
func (a Analyzer) analyzeTypeParameters(
	params []ts.Param,
	scope Scope,
) ([]ts.Param, map[compiler.EntityRef]struct{}, error) {
	if err := a.validator.ValidateParams(params); err != nil {
		return nil, nil, errors.Join(ErrValidator, err)
	}

	resolvedParams := make([]ts.Param, len(params))
	for i, param := range params {
		if param.Constr.Empty() {
			continue
		}
		resolvedConstr, err := a.resolver.Resolve(param.Constr, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrResolver, err), param.Name)
		}
		resolvedParams[i] = ts.Param{
			Name:   param.Name,
			Constr: resolvedConstr,
		}
	}

	return resolvedParams, nil, nil
}

func (a Analyzer) analyzeIO(io compiler.IO, scope Scope, params []ts.Param) (compiler.IO, map[compiler.EntityRef]struct{}, error) {
	return compiler.IO{}, nil, nil
}

func (Analyzer) mergeUsed(used ...map[compiler.EntityRef]struct{}) map[compiler.EntityRef]struct{} {
	result := map[compiler.EntityRef]struct{}{}
	for _, u := range used {
		for k := range u {
			result[k] = struct{}{}
		}
	}
	return result
}

func MustNew(r TypeExprResolver, c SubtypeChecker, v TypeValidator) Analyzer {
	tools.NilPanic(r, c, v)
	return Analyzer{
		resolver:  r,
		checker:   c,
		validator: v,
	}
}
