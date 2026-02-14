package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg/core"
)

func TestValidator_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment
		name    string
		expr    ts.Expr
		wantErr error
	}{
		// both or nothing
		{
			name:    "empty lit and empty inst (default expr)",
			expr:    ts.Expr{},
			wantErr: ts.ErrExprMustBeInstOrLit,
		},
		{
			name: "non-empty lit and inst",
			expr: ts.Expr{
				Lit:  &ts.LitExpr{Union: map[string]*ts.Expr{"a": nil}},
				Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
			},
			wantErr: ts.ErrExprMustBeInstOrLit,
		},
		// non-empty inst
		{
			name:    "empty lit and non-empty inst",
			expr:    h.Inst("int"),
			wantErr: nil,
		},
		// struct
		{
			name:    "empty struct (non-empty lit)",
			expr:    h.Struct(nil),
			wantErr: nil,
		},
		// non-empty struct
		{
			name: "empty struct lit",
			expr: h.Struct(map[string]ts.Expr{
				"foo": h.Inst("int"),
			}),
			wantErr: nil,
		},
		// union
		{
			name: "union with 0 elements",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 1 tag-only element",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 2 tag-only elements",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil, "b": nil},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 3 tag-only elements",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil, "b": nil, "c": nil},
				},
			},
			wantErr: nil,
		},
		// unions with type expressions
		{
			name: "union with one type expression",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"a": {
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "union with with two tags, one member is tag-only",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"a": nil,
						"b": {
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "int"},
							},
						}},
				},
			},
			wantErr: nil,
		},
	}

	v := ts.Validator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.ErrorIs(
				t, v.Validate(tt.expr), tt.wantErr,
			)
		})
	}
}
