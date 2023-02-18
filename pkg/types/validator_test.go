package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"

	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		expr    ts.Expr
		wantErr error
	}{
		// both or nothing
		{
			name:    "empty lit and empty inst (default expr)",
			expr:    ts.Expr{},
			wantErr: ts.ErrInvalidExprType,
		},
		{
			name: "non-empty lit and inst",
			expr: ts.Expr{
				Lit:  ts.LitExpr{Enum: []string{"a"}},
				Inst: ts.InstExpr{Ref: "int"},
			},
			wantErr: ts.ErrInvalidExprType,
		},
		// non-empty inst
		{
			name:    "empty lit and non-empty inst",
			expr:    h.Inst("int"),
			wantErr: nil,
		},
		// rec
		{
			name:    "empty rec (non-empty lit)",
			expr:    h.Rec(nil),
			wantErr: nil,
		},
		// non-empty rec
		{
			name: "empty rec lit",
			expr: h.Rec(map[string]ts.Expr{
				"foo": h.Inst("int"),
			}),
			wantErr: nil,
		},
		// arr
		{
			name: "array of 0 elements",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 0},
				},
			},
			wantErr: ts.ErrArrSize,
		},
		{
			name: "array of 1 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 1},
				},
			},
			wantErr: ts.ErrArrSize,
		},
		{
			name: "array of 2 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 2},
				},
			},
			wantErr: nil,
		},
		{
			name: "array of 3 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 3},
				},
			},
			wantErr: nil,
		},
		// union
		{
			name: "union of 0 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{},
				},
			},
			wantErr: ts.ErrUnionLen,
		},
		{
			name: "union of 1 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{{}},
				},
			},
			wantErr: ts.ErrUnionLen,
		},
		{
			name: "union of 2 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{{}, {}},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 3 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{{}, {}, {}},
				},
			},
			wantErr: nil,
		},
		// enum
		{
			name: "enum of 0 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{},
				},
			},
			wantErr: ts.ErrEnumLen,
		},
		{
			name: "enum of 1 element",
			expr: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{""},
				},
			},
			wantErr: ts.ErrEnumLen,
		},
		{
			name:    "enum of 2 duplicate element",
			expr:    h.Enum("a", "a"),
			wantErr: ts.ErrEnumDupl,
		},
		{
			name:    "enum of 2 diff element",
			expr:    h.Enum("a", "b"),
			wantErr: nil,
		},
		{
			name:    "enum of 3 diff element",
			expr:    h.Enum("a", "b", "c"),
			wantErr: nil,
		},
		{
			name:    "enum of 3 els with dupl",
			expr:    h.Enum("a", "b", "a"),
			wantErr: ts.ErrEnumDupl,
		},
	}

	v := ts.Validator{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(
				t, v.Validate(tt.expr), tt.wantErr,
			)
		})
	}
}
