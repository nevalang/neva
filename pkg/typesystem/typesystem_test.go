package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ts "github.com/nevalang/neva/pkg/typesystem"
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
			name: "all 4 fields: arr, enum, union and rec are empty",
			lit:  ts.LitExpr{nil, nil, nil, nil},
			want: true,
		},
		{
			name: "arr not empty",
			lit:  ts.LitExpr{&ts.ArrLit{}, nil, nil, nil},
			want: false,
		},
		{
			name: "rec not empty",
			lit:  ts.LitExpr{nil, map[string]ts.Expr{}, nil, nil},
			want: false,
		},
		{
			name: "enum not empty",
			lit:  ts.LitExpr{nil, nil, []string{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  ts.LitExpr{nil, nil, nil, []ts.Expr{}},
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
			lit:  ts.LitExpr{nil, nil, nil, nil},
			want: ts.EmptyLitType,
		},
		{
			name: "arr",
			lit:  ts.LitExpr{&ts.ArrLit{}, nil, nil, nil},
			want: ts.ArrLitType,
		},
		{
			name: "rec",
			lit:  ts.LitExpr{nil, map[string]ts.Expr{}, nil, nil},
			want: ts.StructLitType,
		},
		{
			name: "enum",
			lit:  ts.LitExpr{nil, nil, []string{}, nil},
			want: ts.EnumLitType,
		},
		{
			name: "union",
			lit:  ts.LitExpr{nil, nil, nil, []ts.Expr{}},
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
			name: "",
			def: h.Def(
				h.Inst("vec", h.Inst("T")),
				h.Param("T", h.Inst("int")),
			),
			want: "<T int> = vec<T>",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.def.String(); got != tt.want {
				t.Errorf("Def.String() = %v, want %v", got, tt.want)
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
			want: "<str>",
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
			want: "map<str>",
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
			want: "map<str, bool>",
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
								Ref: ts.DefaultStringer("vec"),
								Args: []ts.Expr{
									{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("bool")}},
								},
							},
						},
					},
				},
			},
			want: "map<str, vec<bool>>",
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
		// arr
		{
			name: "lit_expr_arr_with_size_0_and_without type",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Arr: &ts.ArrLit{},
				},
			},
			want: "[0]empty",
		},
		{
			name: "lit expr arr with size 0 and with type",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Arr: &ts.ArrLit{
						Expr: ts.Expr{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")}},
					},
				},
			},
			want: "[0]int",
		},
		{
			name: "lit expr arr with size 4096 and with type",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Arr: &ts.ArrLit{
						Size: 4096,
						Expr: ts.Expr{Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")}},
					},
				},
			},
			want: "[4096]int",
		},
		// rec
		{
			name: "lit expr rec no fields",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{},
				},
			},
			want: "{}",
		},
		{
			name: "lit_expr_rec_with_one_field",
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
			want: "{ name str }",
		},
		{ // FIXME flacky test (struct must be ordered)
			name: "lit_expr_rec_with_two_fields",
			expr: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{
						"name": {Inst: &ts.InstExpr{Ref: ts.DefaultStringer("string")}},
						"age":  {Inst: &ts.InstExpr{Ref: ts.DefaultStringer("int")}},
					},
				},
			},
			want: "{ name str, age int }",
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
			want: "int | str",
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
