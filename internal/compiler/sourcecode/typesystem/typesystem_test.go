package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var h ts.Helper

func TestLiteralExpr_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lit  ts.LitExpr
		want bool
	}{
		{
			name: "all 4 fields: arr, enum, union and struct are empty",
			lit:  ts.LitExpr{nil, nil, nil},
			want: true,
		},
		{
			name: "struct not empty",
			lit:  ts.LitExpr{map[string]ts.Expr{}, nil, nil},
			want: false,
		},
		{
			name: "enum not empty",
			lit:  ts.LitExpr{nil, []string{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  ts.LitExpr{nil, nil, []ts.Expr{}},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.lit.Empty(), tt.want)
		})
	}
}

func TestLiteralExpr_Type(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lit  ts.LitExpr
		want ts.LiteralType
	}{
		{
			name: "unknown",
			lit:  ts.LitExpr{nil, nil, nil},
			want: ts.EmptyLitType,
		},
		{
			name: "struct",
			lit:  ts.LitExpr{map[string]ts.Expr{}, nil, nil},
			want: ts.StructLitType,
		},
		{
			name: "enum",
			lit:  ts.LitExpr{nil, []string{}, nil},
			want: ts.EnumLitType,
		},
		{
			name: "union",
			lit:  ts.LitExpr{nil, nil, []ts.Expr{}},
			want: ts.UnionLitType,
		},
	}

	for _, tt := range tests {
		tt := tt
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
		tt := tt
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
				Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")},
			},
			want: "int",
		},
		{
			name: "inst_expr_with_empty_refs_and_with_args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: nil,
					Args: []ts.Expr{
						{
							Inst: &ts.InstExpr{
								Ref: ts.DefaultStringer("string"),
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
					Ref: ts.DefaultStringer("map"),
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("string")}},
					},
				},
			},
			want: "map<string>",
		},
		{
			name: "inst expr with non-empty refs and with several args",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: ts.DefaultStringer("map"),
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("string")}},
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("bool")}},
					},
				},
			},
			want: "map<string, bool>",
		},
		{
			name: "inst expr with non-empty refs and with nested arg",
			expr: ts.Expr{
				Inst: &ts.InstExpr{
					Ref: ts.DefaultStringer("map"),
					Args: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("string")}},
						{
							Inst: &ts.InstExpr{
								Ref: ts.DefaultStringer("list"),
								Args: []ts.Expr{
									{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("bool")}},
								},
							},
						},
					},
				},
			},
			want: "map<string, list<bool>>",
		},
		// Lits
		// enum
		{
			name: "lit_expr_empty_enum",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Enum: []string{},
				},
			},
			want: "{}",
		},
		{
			name: "lit_expr_enum_with_one_el",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Enum: []string{"MONDAY"},
				},
			},
			want: "{ MONDAY }",
		},
		{
			name: "lit_expr_enum_with_two_els",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Enum: []string{"MONDAY", "TUESDAY"},
				},
			},
			want: "{ MONDAY, TUESDAY }",
		},
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
								Ref: ts.DefaultStringer("string"),
							},
						},
					},
				},
			},
			want: "{ name string }",
		},
		// union
		{
			name: "lit expr empty union", // not a valid expr
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: []ts.Expr{},
				},
			},
			want: "",
		},
		{
			name: "lit expr union with one el", // not a valid expr
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")}},
					},
				},
			},
			want: "int",
		},
		{
			name: "lit expr union with two els",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Union: []ts.Expr{
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")}},
						{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("string")}},
					},
				},
			},
			want: "int | string",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(
				t, tt.want, tt.expr.String(),
			)
		})
	}
}
