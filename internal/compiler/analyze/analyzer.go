package analyze

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrRootPkgMissing       = errors.New("program must have root pkg")
	ErrRootPkgNotFound      = errors.New("root pkg not found")
	ErrMainPkgNotExecutable = errors.New("main pkg must be executable (have root component)")
	ErrPkg                  = errors.New("analyze package")
	ErrTypeParams           = errors.New("interface type parameters")
	ErrIO                   = errors.New("io")
)

var h src.Helper

type Analyzer struct {
	Resolver TypeSystem
	Compator Compator
}

type (
	TypeSystem interface {
		Resolve(ts.Expr, ts.Scope) (ts.Expr, error)
	}
	Compator interface {
		Check(ts.Expr, ts.Expr, ts.TerminatorParams) error
	}
)

// Analyze checks that:
// Program has ref to root pkg;
// Root pkg exist;
// Root pkg has root component ref;
// All pkgs are analyzed;
func (a Analyzer) Analyze(ctx context.Context, prog src.Prog) (src.Prog, error) {
	if prog.RootPkg == "" {
		return src.Prog{}, ErrRootPkgMissing
	}

	mainPkg, ok := prog.Pkgs[prog.RootPkg]
	if !ok {
		return src.Prog{}, fmt.Errorf("%w: %v", ErrRootPkgNotFound, prog.RootPkg)
	}

	if mainPkg.RootComponent == "" {
		return src.Prog{}, ErrMainPkgNotExecutable
	}

	resolvedPkgs := make(map[string]src.Pkg, len(prog.Pkgs))
	for pkgName := range prog.Pkgs {
		resolvedPkg, err := a.analyzePkg(pkgName, prog.Pkgs)
		if err != nil {
			return src.Prog{}, fmt.Errorf("%w: %v", errors.Join(ErrPkg, err), pkgName)
		}
		resolvedPkgs[pkgName] = resolvedPkg
	}

	return src.Prog{
		Pkgs:    resolvedPkgs,
		RootPkg: prog.RootPkg,
	}, nil
}

func (a Analyzer) analyzeTypeParameters(params []ts.Param, scope Scope) ([]ts.Param, map[src.EntityRef]struct{}, error) {
	return nil, nil, nil
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
