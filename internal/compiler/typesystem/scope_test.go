// This file implements Scope interface specifically for tests.
package typesystem_test

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

var ErrDefaultScope = errors.New("default scope")

type TestScope map[string]ts.Def

func (s TestScope) IsTopType(expr ts.Expr) bool {
	if expr.Inst == nil {
		return false
	}
	return expr.Inst.Ref.String() == "any"
}

func (s TestScope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	v, ok := s[ref.String()]
	if !ok {
		return ts.Def{}, nil, ErrDefaultScope
	}
	return v, s, nil
}
