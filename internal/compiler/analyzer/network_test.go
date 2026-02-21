package analyzer

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
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

// createSingleElementUnion creates a union type with a single element matching the given type.
// It's used only by unit tests.
func (a Analyzer) createSingleElementUnion(expr ts.Expr) ts.Expr {
	// if the expression is already a union, return it as-is
	if expr.Lit != nil && expr.Lit.Union != nil {
		return expr
	}

	// create a single-element union
	// for primitive types like int, create union { int }
	// for complex types, create union with the type name as the tag
	if expr.Inst != nil {
		typeName := expr.Inst.Ref.String()
		// create a new instance expression with the same type
		tagExpr := ts.Expr{
			Inst: &ts.InstExpr{
				Ref:  expr.Inst.Ref,
				Args: expr.Inst.Args,
			},
		}
		return ts.Expr{
			Lit: &ts.LitExpr{
				Union: map[string]*ts.Expr{
					typeName: &tagExpr,
				},
			},
		}
	}

	// if the expression is a literal, we need to handle it differently
	if expr.Lit != nil {
		// for literal expressions, we can't easily create a union
		// this shouldn't happen for operator operands, but let's handle it
		return expr
	}

	// fallback: return the expression as-is if we can't create a union
	return expr
}

func TestIsMissingAliasNodeName(t *testing.T) {
	t.Parallel()

	require.True(t, isMissingAliasNodeName(src.MissingAliasNodeName(1)))
	require.False(t, isMissingAliasNodeName("handler"))
}
