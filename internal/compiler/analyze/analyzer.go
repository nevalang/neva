package analyze

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	"github.com/emil14/neva/pkg/tools"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrRootPkgMissing       = errors.New("program must have root pkg")
	ErrRootPkgNotFound      = errors.New("root pkg not found")
	ErrMainPkgNotExecutable = errors.New("main pkg must be executable (have root component)")
	ErrPkg                  = errors.New("analyze package")
	ErrTypeParams           = errors.New("interface type parameters")
	ErrIO                   = errors.New("io")
	ErrValidator            = errors.New("type validator")
	ErrResolver             = errors.New("type expression  resolver")
)

var h src.Helper

type Analyzer struct {
	Resolver  TypeExprResolver
	Checker   SubtypeChecker
	validator TypeValidator
}

type (
	TypeExprResolver interface {
		Resolve(ts.Expr, ts.Scope) (ts.Expr, error)
	}
	SubtypeChecker interface {
		Check(ts.Expr, ts.Expr, ts.TerminatorParams) error
	}
	TypeValidator interface {
		ValidateParams([]ts.Param) error
	}
)

// Analyze checks that:
// Program has ref to root pkg;
// Root pkg exist;
// Root pkg has root component ref;
// All pkgs are analyzed;
func (a Analyzer) Analyze(ctx context.Context, prog src.Prog) (src.Prog, error) {
	if prog.MainPkg == "" {
		return src.Prog{}, ErrRootPkgMissing
	}

	mainPkg, ok := prog.Pkgs[prog.MainPkg]
	if !ok {
		return src.Prog{}, fmt.Errorf("%w: %v", ErrRootPkgNotFound, prog.MainPkg)
	}

	if mainPkg.MainComponent == "" {
		return src.Prog{}, ErrMainPkgNotExecutable
	}

	resolvedPkgs := make(map[string]src.Pkg, len(prog.Pkgs))
	for pkgName := range prog.Pkgs {
		resolvedPkg, err := a.analyzePkg(pkgName, prog.Pkgs)
		if err != nil {
			return src.Prog{}, fmt.Errorf("%w: found in %v", errors.Join(ErrPkg, err), pkgName)
		}
		resolvedPkgs[pkgName] = resolvedPkg
	}

	return src.Prog{
		Pkgs:    resolvedPkgs,
		MainPkg: prog.MainPkg,
	}, nil
}

func (a Analyzer) analyzeTypeParameters(
	params []ts.Param,
	scope Scope,
	args map[string]ts.Expr,
) ([]ts.Param, map[src.EntityRef]struct{}, error) {
	if args != nil && len(args) != len(params) {
		panic("len")
	}

	if err := a.validator.ValidateParams(params); err != nil {
		return nil, nil, errors.Join(ErrValidator, err)
	}

	resolvedParams := make([]ts.Param, len(params))
	for i, param := range params {
		if param.Constr.Empty() {
			if _, ok := args[param.Name]; ok {
				panic("")
			}
			continue
		}
		resolvedConstr, err := a.Resolver.Resolve(param.Constr, scope)
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

func (a Analyzer) analyzeIO(io src.IO, scope Scope, params []ts.Param) (src.IO, map[src.EntityRef]struct{}, error) {
	return src.IO{}, nil, nil
}

func (Analyzer) mergeUsed(used ...map[src.EntityRef]struct{}) map[src.EntityRef]struct{} {
	result := map[src.EntityRef]struct{}{}
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
		Resolver:  r,
		Checker:   c,
		validator: v,
	}
}
