package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/stretchr/testify/require"
)

func TestSubTypeChecker_SubTypeCheck(t *testing.T) { //nolint:maintidx
	tests := []struct {
		name    string
		arg     ts.Expr
		constr  ts.Expr
		wantErr error
	}{
		// Instantiations
		{
			name:    "arg and constr are default values (empty insts)",
			arg:     ts.Expr{},
			constr:  ts.Expr{},
			wantErr: nil,
		},
		//  kinds
		{
			name:    "arg inst, constr lit (not union)", // int <: {}
			arg:     h.Inst("int"),
			constr:  h.Enum(),
			wantErr: ts.ErrDiffExprTypes,
		},
		{
			name:    "constr inst, arg lit (not union)", // {} <: int
			arg:     h.Enum(),
			constr:  h.Inst("int"),
			wantErr: ts.ErrDiffExprTypes,
		},
		// diff refs
		{
			name:    "insts, diff refs, no args", // int <: bool (no need to check vice versa, they resolved)
			arg:     h.Inst("int"),
			constr:  h.Inst("bool"),
			wantErr: ts.ErrDiffRefs,
		},
		{
			name:    "insts, same refs, no args", // int <: int
			arg:     h.Inst("int"),
			constr:  h.Inst("int"),
			wantErr: nil,
		},
		// args count
		{
			name:    "insts, arg has less args", // vec <: vec<int>
			arg:     h.Inst("vec"),
			constr:  h.Inst("vec", h.Inst("int")),
			wantErr: ts.ErrArgsCount,
		},
		{
			name:    "insts, arg has same args count", // vec<int> <: vec<int>
			arg:     h.Inst("vec", h.Inst("int")),
			constr:  h.Inst("vec", h.Inst("int")),
			wantErr: nil,
		},
		{
			name:    "insts, arg has more args count", // vec<int, str> <: vec<int>
			arg:     h.Inst("vec", h.Inst("int"), h.Inst("str")),
			constr:  h.Inst("vec", h.Inst("int")),
			wantErr: nil,
		},
		// args compatibility
		{
			name: "insts, one arg's arg incompat", // vec<str> <: vec<int|str>
			arg:  h.Inst("vec", h.Inst("str")),
			constr: h.Inst(
				"vec",
				h.Union(
					h.Inst("str"),
					h.Inst("int"),
				),
			),
			wantErr: nil,
		},
		{
			name: "insts, constr arg incompat", // vec<str|int> <: vec<int>
			arg: h.Inst(
				"vec",
				h.Union(
					h.Inst("str"),
					h.Inst("int"),
				),
			),
			constr:  h.Inst("vec", h.Inst("int")),
			wantErr: ts.ErrArgNotSubtype,
		},
		// arr
		{
			name: "expr and constr has diff lit types (constr not union)",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{Enum: []string{}},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{Arr: &ts.ArrLit{}},
			},
			wantErr: ts.ErrDiffLitTypes,
		},
		{
			name: "expr's arr lit has lesser size than constr",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{Size: 1}, // expr doesn't matter here
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{Size: 2},
				},
			},
			wantErr: ts.ErrLitArrSize,
		},
		{
			name: "expr's arr has incompat type",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{
						Size: 2, // same size itself won't cause any problem
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "a"}, // type will
						},
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{
						Size: 2,
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "b"},
						},
					},
				},
			},
			wantErr: ts.ErrArrDiffType,
		},
		{
			name: "expr and constr arrs, expr is bigger and have compat type",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{
						Size: 3, // bigger size won't cause any problem
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "a"}, // same type
						},
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Arr: &ts.ArrLit{
						Size: 2,
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "a"},
						},
					},
				},
			},
			wantErr: nil,
		},
		// enum
		{
			name: "expr and constr enums, expr is bigger",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a", "b"},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a"},
				},
			},
			wantErr: ts.ErrBigEnum,
		},
		{
			name: "expr and constr enums, expr not bigger but contain diff el",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a", "d"}, // d != b
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a", "b", "c"},
				},
			},
			wantErr: ts.ErrEnumEl,
		},
		{
			name: "expr and constr enums, expr not bigger and all reqired els are the same",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a", "b"},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Enum: []string{"a", "b", "c"}, // c el won't cause any problem
				},
			},
			wantErr: nil,
		},
		// rec
		{
			name: "expr and constr recs, expr has less fields",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{}, // 0 fields is ok
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{
						"a": {}, // expr itself doesn't matter here
					},
				},
			},
			wantErr: ts.ErrRecLen,
		},
		{
			name: "expr and constr recs, expr leaks field",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"b": {}, // expr itself doesn't matter here
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{
						"a": {}, // but this field is missing
					},
				},
			},
			wantErr: ts.ErrRecNoField,
		},
		{
			name: "expr and constr recs, expr has incompat field",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"b": {}, // b field itself won't cause any problems
						"a": {}, // this one will
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{
						"a": {
							Inst: ts.InstExpr{Ref: "x"}, // not same as in expr
						},
					},
				},
			},
			wantErr: ts.ErrRecField,
		},
		{
			name: "expr and constr recs, expr has all constr fields, all fields compatible",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"a": {
							Inst: ts.InstExpr{Ref: "x"}, // not same as in expr
						},
						"b": {}, // b field itself won't cause any problems
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Rec: map[string]ts.Expr{
						"a": {
							Inst: ts.InstExpr{Ref: "x"}, // not same as in expr
						},
					},
				},
			},
			wantErr: nil,
		},
		// union
		{
			name: "expr inst, constr union. expr incompat with all els",
			arg: ts.Expr{
				Inst: ts.InstExpr{Ref: "x"}, // not compat with both a and b
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: ts.ErrUnion,
		},
		{
			name: "expr not union, constr is. expr is compat with one el",
			arg: ts.Expr{
				Inst: ts.InstExpr{Ref: "b"},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "expr and constr are unions, expr has more els",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
						{Inst: ts.InstExpr{Ref: "c"}},
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
					},
				},
			},
			wantErr: ts.ErrUnionsLen,
		},
		{
			name: "expr and constr are unions, same size but incompat expr el",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "c"}},
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "x"}}, // this
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
						{Inst: ts.InstExpr{Ref: "c"}},
					},
				},
			},
			wantErr: ts.ErrUnions,
		},
		{
			name: "expr and constr are unions, expr is less and compat",
			arg: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "c"}},
						{Inst: ts.InstExpr{Ref: "a"}},
					},
				},
			},
			constr: ts.Expr{
				Lit: ts.LiteralExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
						{Inst: ts.InstExpr{Ref: "c"}},
					},
				},
			},
			wantErr: nil,
		},
		// name: "expr and constr are unions, expr is same and compat",
	}

	s := ts.SubTypeChecker{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.ErrorIs(
				t,
				s.SubtypeCheck(tt.arg, tt.constr),
				tt.wantErr,
			)
		})
	}
}
