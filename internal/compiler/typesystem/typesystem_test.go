package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

var h ts.Helper

func TestLiteralExpr_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment
		name string
		lit  ts.LitExpr
		want bool
	}{
		{
			name: "all 2 fields: structs and unions are empty",
			lit:  ts.LitExpr{nil, nil},
			want: true,
		},
		{
			name: "struct not empty",
			lit:  ts.LitExpr{map[string]ts.Expr{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  ts.LitExpr{nil, map[string]*ts.Expr{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.lit.Empty(), tt.want)
		})
	}
}

func TestLiteralExpr_Type(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment
		name string
		lit  ts.LitExpr
		want ts.LiteralType
	}{
		{
			name: "unknown",
			lit:  ts.LitExpr{nil, nil},
			want: ts.EmptyLitType,
		},
		{
			name: "struct",
			lit:  ts.LitExpr{map[string]ts.Expr{}, nil},
			want: ts.StructLitType,
		},
		{
			name: "union",
			lit:  ts.LitExpr{nil, map[string]*ts.Expr{}},
			want: ts.UnionLitType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.lit.Type(), tt.want)
		})
	}
}

func TestDef_String(t *testing.T) {
	tests := []struct {
		name string
		def  ts.Def
		want string
	}{
		{
			name: "<T_int>_=_list<T>",
			def: h.Def(
				h.Inst("list", h.Inst("T")),
				h.Param("T", h.Inst("int")),
			),
			want: "<T int> = list<T>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.def.String(); got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpr_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		expr ts.Expr
		want string
	}{
		// insts
		{
			name: "empty_expr_(inst_with_empty_ref_and_no_args)",
			expr: ts.Expr{},
			want: "empty",
		},
		{
			name: "inst_expr_with_non-empty_ref_and_no_args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
			},
			want: "int",
		},
		{
			name: "inst_expr_with_empty_refs_and_with_args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Args: []ts.Expr{
						{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "string"},
							},
						},
					},
				},
			},
			want: "<string>",
		},
		{
			name: "inst expr with non-empty refs and with args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "map"},
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
					},
				},
			},
			want: "map<string>",
		},
		{
			name: "inst expr with non-empty refs and with several args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "map"},
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}}},
					},
				},
			},
			want: "map<string, bool>",
		},
		{
			name: "inst expr with non-empty refs and with nested arg",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: core.EntityRef{Name: "map"},
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
						{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "list"},
								Args: []ts.Expr{
									{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "bool"}}},
								},
							},
						},
					},
				},
			},
			want: "map<string, list<bool>>",
		},
		// --- Lits ---
		// struct
		{
			name: "lit expr struct no fields",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{},
				},
			},
			want: "{}",
		},
		{
			name: "lit_expr_struct_with_one_field",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{
						"name": {
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Name: "string"},
							},
						},
					},
				},
			},
			want: "{ name string }",
		},
		// union (tag-only)
		{
			name: "lit expr empty union",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{},
				},
			},
			want: "union {}",
		},
		{
			name: "lit expr union with one el",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
					},
				},
			},
			want: "union { int }",
		},
		{
			name: "lit expr union with two els",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"int":    {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
						"string": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
					},
				},
			},
			want: "union { int, string }",
		},
		// unions (tag-and-value)
		{
			name: "lit expr union with one el (tag-and-value)",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
					},
				},
			},
			want: "union { Int int }",
		},
		{
			name: "lit expr union with two el (tag-and-value)",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: map[string]*ts.Expr{
						"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
						"Str": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
					},
				},
			},
			want: "union { Int int, Str string }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(
				t, tt.want, tt.expr.String(),
			)
		})
	}
}
