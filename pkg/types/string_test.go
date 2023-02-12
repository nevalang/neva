package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestExpr_String(t *testing.T) {
	tests := []struct {
		name string
		expr types.Expr
		want string
	}{
		// insts
		{
			name: "empty expr (inst with empty ref and no args)",
			expr: types.Expr{},
			want: "",
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
			name: "lit expr empty enum",
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
			name: "lit expr arr with size 0 and without type",
			expr: types.Expr{
				Lit: types.LitExpr{
					Arr: &types.ArrLit{},
				},
			},
			want: "[0]",
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
			require.Equal(
				t, tt.want, tt.expr.String(),
			)
		})
	}
}
