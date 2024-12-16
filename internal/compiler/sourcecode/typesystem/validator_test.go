package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
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
			name: "union of 0 element",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 1 element",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 2 element",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil, "b": nil},
				},
			},
			wantErr: nil,
		},
		{
			name: "union of 3 element",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{"a": nil, "b": nil, "c": nil},
				},
			},
			wantErr: nil,
		},
		// TODO add unit tests for unions with type expressions
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
