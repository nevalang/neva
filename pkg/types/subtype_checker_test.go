package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestSubTypeChecker_SubTypeCheck(t *testing.T) { //nolint:maintidx
	tests := []struct {
		name    string
		expr    types.Expr
		constr  types.Expr
		wantErr error
	}{
		// Instantiations
		{
			name:   "expr and constr are default values (empty insts)", // ''<> <: ''<>
			expr:   types.Expr{},
			constr: types.Expr{},
		},
		{
			name: "expr is inst and constr is lit (not union)", // ''<> !<: {}
			expr: types.Expr{
				Inst: types.InstExpr{Ref: ""},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{EnumLit: []string{}}, // content doesn't matter here
			},
			wantErr: types.ErrDiffTypes,
		},
		{
			name: "expr is lit (not union) and constr is inst", // {} !<: ''<>
			expr: types.Expr{
				Lit: types.LiteralExpr{EnumLit: []string{}}, // content doesn't matter here
			},
			constr: types.Expr{
				Inst: types.InstExpr{Ref: ""},
			},
			wantErr: types.ErrDiffTypes,
		},
		{
			name:    "expr and constr are insts has with different refs", // a<> !<: b<>
			expr:    types.Expr{Inst: types.InstExpr{Ref: "a"}},
			constr:  types.Expr{Inst: types.InstExpr{Ref: "b"}},
			wantErr: types.ErrDiffRefs,
		},
		{
			name: "expr inst, same refs, but expr has less args", // a<> !<: b<int>
			expr: types.Expr{Inst: types.InstExpr{Ref: "a"}},
			constr: types.Expr{
				Inst: types.InstExpr{
					Ref: "a",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "int"}}, // arg itself doesn't matter here
					},
				},
			},
			wantErr: types.ErrArgsLen,
		},
		{
			name: "expr inst, same refs and args count, but one arg is incompatible", // a<str> !<: a<int>
			expr: types.Expr{
				Inst: types.InstExpr{
					Ref: "a",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "str"}}, // str !<: int
					},
				},
			},
			constr: types.Expr{
				Inst: types.InstExpr{
					Ref: "a",
					Args: []types.Expr{
						{Inst: types.InstExpr{Ref: "int"}},
					},
				},
			},
			wantErr: types.ErrArg,
		},
		// arr
		{
			name: "expr and constr has diff lit types (constr not union)",
			expr: types.Expr{
				Lit: types.LiteralExpr{EnumLit: []string{}},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{ArrLit: &types.ArrLit{}},
			},
			wantErr: types.ErrDiffLitTypes,
		},
		{
			name: "expr's arr lit has lesser size than constr",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 1}, // expr doesn't matter here
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{Size: 2},
				},
			},
			wantErr: types.ErrLitArrSize,
		},
		{
			name: "expr's arr has incompat type",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{
						Size: 2, // same size itself won't cause any problem
						Expr: types.Expr{
							Inst: types.InstExpr{Ref: "a"}, // type will
						},
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{
						Size: 2,
						Expr: types.Expr{
							Inst: types.InstExpr{Ref: "b"},
						},
					},
				},
			},
			wantErr: types.ErrArrDiffType,
		},
		{
			name: "expr and constr arrs, expr is bigger and have compat type",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{
						Size: 3, // bigger size won't cause any problem
						Expr: types.Expr{
							Inst: types.InstExpr{Ref: "a"}, // same type
						},
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					ArrLit: &types.ArrLit{
						Size: 2,
						Expr: types.Expr{
							Inst: types.InstExpr{Ref: "a"},
						},
					},
				},
			},
			wantErr: nil,
		},
		// enum
		{
			name: "expr and constr enums, expr is bigger",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a", "b"},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a"},
				},
			},
			wantErr: types.ErrBigEnum,
		},
		{
			name: "expr and constr enums, expr not bigger but contain diff el",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a", "d"}, // d != b
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a", "b", "c"},
				},
			},
			wantErr: types.ErrEnumEl,
		},
		{
			name: "expr and constr enums, expr not bigger and all reqired els are the same",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a", "b"},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					EnumLit: []string{"a", "b", "c"}, // c el won't cause any problem
				},
			},
			wantErr: nil,
		},
		// rec
		{
			name: "expr and constr recs, expr has less fields",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{}, // 0 fields is ok
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{
						"a": {}, // expr itself doesn't matter here
					},
				},
			},
			wantErr: types.ErrRecLen,
		},
		{
			name: "expr and constr recs, expr leaks field",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{ // both has 1 field
						"b": {}, // expr itself doesn't matter here
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{
						"a": {}, // but this field is missing
					},
				},
			},
			wantErr: types.ErrRecNoField,
		},
		{
			name: "expr and constr recs, expr has incompat field",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{ // both has 1 field
						"b": {}, // b field itself won't cause any problems
						"a": {}, // this one will
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{
						"a": {
							Inst: types.InstExpr{Ref: "x"}, // not same as in expr
						},
					},
				},
			},
			wantErr: types.ErrRecField,
		},
		{
			name: "expr and constr recs, expr has all constr fields, all fields compatible",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{ // both has 1 field
						"a": {
							Inst: types.InstExpr{Ref: "x"}, // not same as in expr
						},
						"b": {}, // b field itself won't cause any problems
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					RecLit: map[string]types.Expr{
						"a": {
							Inst: types.InstExpr{Ref: "x"}, // not same as in expr
						},
					},
				},
			},
			wantErr: nil,
		},
		// union
		{
			name: "expr inst, constr union. expr incompat with all els",
			expr: types.Expr{
				Inst: types.InstExpr{Ref: "x"}, // not compat with both a and b
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: types.ErrUnion,
		},
		{
			name: "expr not union, constr is. expr is compat with one el",
			expr: types.Expr{
				Inst: types.InstExpr{Ref: "b"},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "expr and constr are unions, expr has more els",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
						{Inst: types.InstExpr{Ref: "c"}},
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: types.ErrUnionsLen,
		},
		{
			name: "expr and constr are unions, same size but incompat expr el",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "c"}},
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "x"}}, // this
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
						{Inst: types.InstExpr{Ref: "c"}},
					},
				},
			},
			wantErr: types.ErrUnions,
		},
		{
			name: "expr and constr are unions, expr is less and compat",
			expr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "c"}},
						{Inst: types.InstExpr{Ref: "a"}},
					},
				},
			},
			constr: types.Expr{
				Lit: types.LiteralExpr{
					UnionLit: []types.Expr{
						{Inst: types.InstExpr{Ref: "a"}},
						{Inst: types.InstExpr{Ref: "b"}},
						{Inst: types.InstExpr{Ref: "c"}},
					},
				},
			},
			wantErr: nil,
		},
	}

	s := types.SubTypeChecker{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(
				t,
				s.SubtypeCheck(tt.expr, tt.constr),
				tt.wantErr,
			)
		})
	}
}
