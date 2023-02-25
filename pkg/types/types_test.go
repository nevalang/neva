package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ts "github.com/emil14/neva/pkg/types"
	types "github.com/emil14/neva/pkg/types"
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
			want: ts.RecLitType,
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

func TestInstExpr_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		inst ts.InstExpr
		want bool
	}{
		{
			name: "default inst (empty ref and nil args)",
			inst: ts.InstExpr{
				Ref:  "",
				Args: nil,
			},
			want: true,
		},
		{
			name: "empty ref and empty list args",
			inst: ts.InstExpr{
				Ref:  "",
				Args: []ts.Expr{},
			},
			want: true,
		},
		{
			name: "empty ref and non empty list args",
			inst: ts.InstExpr{
				Ref:  "",
				Args: []ts.Expr{{}}, // content doesn't matter here
			},
			want: false,
		},
		{
			name: "non-empty ref and non empty list args",
			inst: ts.InstExpr{
				Ref:  "t",
				Args: []ts.Expr{{}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.inst.Empty(), tt.want)
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
		expr types.Expr
		want string
	}{
		// insts
		{
			name: "empty_expr_(inst_with_empty_ref_and_no_args)",
			expr: types.Expr{},
			want: "empty",
		},
		{
			name: "inst expr with non-empty ref and no args",
			expr: types.Expr{
				Inst: types.InstExpr{Ref: "int"},
			},
			want: "int",
		},
		{
			name: "inst expr with empty refs and with args",
			expr: types.Expr{
				Inst: types.InstExpr{
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "str"}},
					},
				},
			},
			want: "<str>",
		},
		{
			name: "inst expr with non-empty refs and with args",
			expr: types.Expr{
				Inst: types.InstExpr{
					Ref: "map",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "str"}},
					},
				},
			},
			want: "map<str>",
		},
		{
			name: "inst expr with non-empty refs and with several args",
			expr: types.Expr{
				Inst: types.InstExpr{
					Ref: "map",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "str"}},
						{Inst: types.InstExpr{Ref: "bool"}},
					},
				},
			},
			want: "map<str, bool>",
		},
		{
			name: "inst expr with non-empty refs and with nested arg",
			expr: types.Expr{
				Inst: types.InstExpr{
					Ref: "map",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "str"}},
						{
							Inst: types.InstExpr{
								Ref: "vec",
								Args: []types.Expr{
									{Inst: types.InstExpr{Ref: "bool"}},
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
			expr: types.Expr{
				Lit: types.LitExpr{
					Enum: []string{},
				},
			},
			want: "{}",
		},
		{
			name: "lit expr enum with one el",
			expr: types.Expr{
				Lit: types.LitExpr{
					Enum: []string{"MONDAY"},
				},
			},
			want: "{ MONDAY }",
		},
		{
			name: "lit expr enum with two els",
			expr: types.Expr{
				Lit: types.LitExpr{
					Enum: []string{"MONDAY", "TUESDAY"},
				},
			},
			want: "{ MONDAY, TUESDAY }",
		},
		// arr
		{
			name: "lit_expr_arr_with_size_0_and_without type",
			expr: types.Expr{
				Lit: types.LitExpr{
					Arr: &types.ArrLit{},
				},
			},
			want: "[0]empty",
		},
		{
			name: "lit expr arr with size 0 and with type",
			expr: types.Expr{
				Lit: types.LitExpr{
					Arr: &types.ArrLit{
						Expr: types.Expr{Inst: types.InstExpr{Ref: "int"}},
					},
				},
			},
			want: "[0]int",
		},
		{
			name: "lit expr arr with size 4096 and with type",
			expr: types.Expr{
				Lit: types.LitExpr{
					Arr: &types.ArrLit{
						Size: 4096,
						Expr: types.Expr{Inst: types.InstExpr{Ref: "int"}},
					},
				},
			},
			want: "[4096]int",
		},
		// rec
		{
			name: "lit expr rec no fields",
			expr: types.Expr{
				Lit: types.LitExpr{
					Rec: map[string]types.Expr{},
				},
			},
			want: "{}",
		},
		{
			name: "lit expr rec with one field",
			expr: types.Expr{
				Lit: types.LitExpr{
					Rec: map[string]types.Expr{
						"name": {Inst: types.InstExpr{Ref: "str"}},
					},
				},
			},
			want: "{ name str }",
		},
		{ // FIXME flacky test (struct must be ordered)
			name: "lit expr rec with two fields",
			expr: types.Expr{
				Lit: types.LitExpr{
					Rec: map[string]types.Expr{
						"name": {Inst: types.InstExpr{Ref: "str"}},
						"age":  {Inst: types.InstExpr{Ref: "int"}},
					},
				},
			},
			want: "{ name str, age int }",
		},
		// union
		{
			name: "lit expr empty union", // not a valid expr
			expr: types.Expr{
				Lit: types.LitExpr{
					Union: []types.Expr{},
				},
			},
			want: "",
		},
		{
			name: "lit expr union with one el", // not a valid expr
			expr: types.Expr{
				Lit: types.LitExpr{
					Union: []types.Expr{
						{Inst: types.InstExpr{Ref: "int"}},
					},
				},
			},
			want: "int",
		},
		{
			name: "lit expr union with two els",
			expr: types.Expr{
				Lit: types.LitExpr{
					Union: []types.Expr{
						{Inst: types.InstExpr{Ref: "int"}},
						{Inst: types.InstExpr{Ref: "str"}},
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
