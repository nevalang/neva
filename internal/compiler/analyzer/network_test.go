package analyzer

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
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
