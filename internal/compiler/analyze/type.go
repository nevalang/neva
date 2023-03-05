package analyze

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrScopeGetLocalType = errors.New("scope get local type")
	ErrResolver          = errors.New("type expression  resolver")
)

func (a Analyzer) analyzeType(name string, scope Scope) (ts.Def, map[src.EntityRef]struct{}, error) {
	def, err := scope.getLocalType(name)
	if err != nil {
		return ts.Def{}, nil, errors.Join(ErrScopeGetLocalType, err)
	}

	testExpr := ts.Expr{
		Inst: ts.InstExpr{
			Ref:  name,
			Args: a.getTestExprArgs(def),
		},
	}

	// TODO return simplified defs (type t1 pkg1.t0<t0> // t1<int> -> vec<int>)
	if _, err = a.Resolver.Resolve(testExpr, scope); err != nil {
		return ts.Def{}, nil, fmt.Errorf("%w: %v", errors.Join(ErrResolver, err), testExpr)
	}

	return def, scope.visited, nil
}

func (Analyzer) getTestExprArgs(def ts.Def) []ts.Expr {
	args := make([]ts.Expr, 0, len(def.Params))
	for _, param := range def.Params {
		if param.Constr.Empty() {
			args = append(args, h.Inst("int"))
		} else {
			args = append(args, param.Constr)
		}
	}
	return args
}
