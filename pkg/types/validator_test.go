package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		expr    types.Expr
		wantErr error
	}{
		// non checkable cases
		{
			name:    "inst (default expr is empty lit and default inst)",
			expr:    types.Expr{}, // empty lit means inst
			wantErr: nil,
		},
		{
			name: "rec lit",
			expr: types.Expr{
				Lit: types.LiteralExpr{RecLit: map[string]types.Expr{}}, // fields doesn't matter here
			},
			wantErr: nil,
		},
		// inst and lit at the same time
		{
			name: "non empty lit and inst with non empty ref",
			expr: types.Expr{
				Lit:  types.LiteralExpr{EnumLit: []string{}}, // it's ok to have "" ref for inst
				Inst: types.InstExpr{Ref: "x"},               // but it's not ok to have a non-empty ref for lit
			},
			wantErr: types.ErrInvalidExprType,
		},
		{
			name: "non empty lit and inst with non empty ref",
			expr: types.Expr{
				Lit:  types.LiteralExpr{EnumLit: []string{}},
				Inst: types.InstExpr{Args: []types.Expr{{}}}, // same for len(args)>0
			},
			wantErr: types.ErrInvalidExprType,
		},
		// arr
		{
			name: "array of 0 elements",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 0},
				},
			},
			wantErr: types.ErrArrSize,
		},
		{
			name: "array of 1 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 1},
				},
			},
			wantErr: types.ErrArrSize,
		},
		{
			name: "array of 2 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 2},
				},
			},
			wantErr: nil,
		},
		{
			name: "array of 3 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 3},
				},
			},
			wantErr: nil,
		},
		// union
		{
			name: "union of 0 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{},
				},
			},
			wantErr: types.ErrUnionLen,
		},
		{
			name: "union of 1 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{{}},
				},
			},
			wantErr: types.ErrUnionLen,
		},
		{
			name: "union of 2 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{{}, {}},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 3 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{{}, {}, {}},
				},
			},
			wantErr: nil,
		},
		// enum
		{
			name: "enum of 0 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{},
				},
			},
			wantErr: types.ErrEnumLen,
		},
		{
			name: "enum of 1 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{""},
				},
			},
			wantErr: types.ErrEnumLen,
		},
		{
			name: "enum of 2 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"", ""},
				},
			},
			wantErr: nil,
		},
		{
			name: "enum of 3 element",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"", "", ""},
				},
			},
			wantErr: nil,
		},
	}

	v := types.Validator{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(
				t, v.Validate(tt.expr), tt.wantErr,
			)
		})
	}
}
