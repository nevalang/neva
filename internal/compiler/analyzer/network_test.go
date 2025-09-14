package analyzer

import (
	"testing"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
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

// testScope implements ts.Scope interface for testing
type testScope struct{}

func (s *testScope) GetType(ref core.EntityRef) (ts.Def, ts.Scope, error) {
	return ts.Def{}, s, nil
}

func (s *testScope) IsTopType(expr ts.Expr) bool {
	return false
}

func TestGetOperatorConstraint(t *testing.T) {
	analyzer := Analyzer{}

	tests := []struct {
		name        string
		operator    src.BinaryOperator
		expected    string
		description string
	}{
		{
			name:        "add_operator",
			operator:    src.AddOp,
			expected:    "union { int, float, string }",
			description: "+ operator should support int, float, string",
		},
		{
			name:        "multiply_operator",
			operator:    src.MulOp,
			expected:    "union { int, float }",
			description: "* operator should support int, float",
		},
		{
			name:        "modulo_operator",
			operator:    src.ModOp,
			expected:    "int",
			description: "% operator should support int only",
		},
		{
			name:        "equal_operator",
			operator:    src.EqOp,
			expected:    "any",
			description: "== operator should support any type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			binary := src.Binary{
				Operator: tt.operator,
				Meta:     core.Meta{},
			}

			constraint, err := analyzer.getOperatorConstraint(binary)
			if err != nil {
				t.Fatalf("getOperatorConstraint failed: %v", err)
			}
			require.Equal(t, tt.expected, constraint.String(), tt.description)
		})
	}
}

func TestOperatorConstraintAndUnionCreation(t *testing.T) {
	analyzer := Analyzer{}

	tests := []struct {
		name        string
		operator    src.BinaryOperator
		leftType    ts.Expr
		rightType   ts.Expr
		description string
	}{
		{
			name:     "int_plus_int",
			operator: src.AddOp,
			leftType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			rightType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			description: "int + int should create proper unions and constraints",
		},
		{
			name:     "int_plus_string",
			operator: src.AddOp,
			leftType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			rightType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "string"},
				},
			},
			description: "int + string should create proper unions and constraints",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create binary expression
			binary := src.Binary{
				Operator: tt.operator,
				Meta:     core.Meta{},
			}

			// test constraint creation
			constraint, err := analyzer.getOperatorConstraint(binary)
			if err != nil {
				t.Fatalf("getOperatorConstraint failed: %v", err)
			}
			require.NotNil(t, constraint, "constraint should not be nil")

			// test union creation
			leftUnion := analyzer.createSingleElementUnion(tt.leftType)
			rightUnion := analyzer.createSingleElementUnion(tt.rightType)

			// verify left union structure
			require.Contains(t, leftUnion.Lit.Union, tt.leftType.Inst.Ref.Name, "left union should contain the type name")

			// verify right union structure
			require.Contains(t, rightUnion.Lit.Union, tt.rightType.Inst.Ref.Name, "right union should contain the type name")

			t.Logf("Constraint: %s", constraint.String())
			t.Logf("Left Union: %s", leftUnion.String())
			t.Logf("Right Union: %s", rightUnion.String())
		})
	}
}

func TestCheckOperatorOperandTypesWithTypeSystem(t *testing.T) {
	// test the function components without the complex resolver
	analyzer := Analyzer{}

	tests := []struct {
		name        string
		operator    src.BinaryOperator
		leftType    ts.Expr
		rightType   ts.Expr
		description string
	}{
		{
			name:     "int_plus_int",
			operator: src.AddOp,
			leftType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			rightType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			description: "int + int should create proper unions and constraints",
		},
		{
			name:     "int_plus_string",
			operator: src.AddOp,
			leftType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "int"},
				},
			},
			rightType: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "string"},
				},
			},
			description: "int + string should create proper unions and constraints",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create binary expression
			binary := src.Binary{
				Operator: tt.operator,
				Meta:     core.Meta{},
			}

			// test constraint creation
			constraint, err := analyzer.getOperatorConstraint(binary)
			if err != nil {
				t.Fatalf("getOperatorConstraint failed: %v", err)
			}
			require.NotNil(t, constraint, "constraint should not be nil")

			// test union creation
			leftUnion := analyzer.createSingleElementUnion(tt.leftType)
			rightUnion := analyzer.createSingleElementUnion(tt.rightType)

			// verify that the function creates the expected structures
			// this tests the core logic without the complex resolver
			require.Contains(t, leftUnion.Lit.Union, tt.leftType.Inst.Ref.Name, "left union should contain the type name")
			require.Contains(t, rightUnion.Lit.Union, tt.rightType.Inst.Ref.Name, "right union should contain the type name")

			t.Logf("Test: %s", tt.description)
			t.Logf("Constraint: %s", constraint.String())
			t.Logf("Left Union: %s", leftUnion.String())
			t.Logf("Right Union: %s", rightUnion.String())
		})
	}
}
