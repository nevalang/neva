package analyzer

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/stretchr/testify/require"
)

func TestCreateSingleElementUnion(t *testing.T) {
	a := Analyzer{}

	// test with int type
	intType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		},
	}

	result := a.createSingleElementUnion(intType)

	// should create union { int }
	require.Equal(t, "int", result.Lit.Union["int"].Inst.Ref.Name)

	// test with string type
	stringType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "string"},
		},
	}

	result2 := a.createSingleElementUnion(stringType)

	// should create union { string }
	require.Equal(t, "string", result2.Lit.Union["string"].Inst.Ref.Name)

	// test with existing union (should return as-is)
	existingUnion := ts.Expr{
		Lit: &ts.LitExpr{
			Union: map[string]*ts.Expr{
				"int": {
					Inst: &ts.InstExpr{
						Ref: core.EntityRef{Name: "int"},
					},
				},
			},
		},
	}

	result3 := a.createSingleElementUnion(existingUnion)
	require.Equal(t, existingUnion, result3)
}

// testScope implements ts.Scope interface for testing
type testScope struct{}

func (s *testScope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	return ts.Def{}, s, nil
}

func (s *testScope) IsTopType(expr ts.Expr) bool {
	return false
}
