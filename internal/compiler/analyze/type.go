package analyze

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrScopeGetLocalType = errors.New("scope get local type")
)

// analyzeType builds text expr that refers to given type and tries to resolve it.
func (a Analyzer) analyzeType(name string, scope Scope) (ts.Def, map[src.EntityRef]struct{}, error) {
	def, err := scope.getLocalType(name)
	if err != nil {
		return ts.Def{}, nil, errors.Join(ErrScopeGetLocalType, err)
	}

	testExpr := ts.Expr{
		Inst: ts.InstExpr{
			Ref:  name,
			Args: a.buildTestExprArgs(def.Params),
		},
	}

	if _, err = a.resolver.Resolve(testExpr, scope); err != nil {
		return ts.Def{}, nil, fmt.Errorf("%w: %v", errors.Join(ErrResolver, err), testExpr)
	}

	return def, scope.visited, nil
}

// buildTestExprArgs takes params and returns args that can be used to build the test expression
func (Analyzer) buildTestExprArgs(params []ts.Param) []ts.Expr {
	args := make([]ts.Expr, 0, len(params))
	for _, param := range params {
		if param.Constr.Empty() {
			args = append(args, h.Inst("int"))
		} else {
			args = append(args, param.Constr)
		}
	}
	return args
}
