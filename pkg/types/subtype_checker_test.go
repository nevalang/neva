package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSubTypeChecker_SubTypeCheck(t *testing.T) { //nolint:maintidx
	tests := []struct {
		name      string
		subType   ts.Expr
		trace1    ts.Trace
		superType ts.Expr
		trace2    ts.Trace
		scope     map[string]ts.Def
		checker   func(*MockrecursionCheckerMockRecorder)
		wantErr   error
		enabled   bool
	}{
		// Instantiations
		{
			name:      "arg and constr are default values (empty insts)",
			subType:   ts.Expr{},
			superType: ts.Expr{},
			wantErr:   nil,
		},
		//  kinds
		{
			name:      "arg inst, constr lit (not union)", // int <: {}
			subType:   h.Inst("int"),
			superType: h.Enum(),
			wantErr:   ts.ErrDiffExprTypes,
		},
		{
			name:      "constr inst, arg lit (not union)", // {} <: int
			subType:   h.Enum(),
			superType: h.Inst("int"),
			wantErr:   ts.ErrDiffExprTypes,
		},
		// diff refs
		{
			name:      "insts, diff refs, no args", // int <: bool (no need to check vice versa, they resolved)
			subType:   h.Inst("int"),
			superType: h.Inst("bool"),
			wantErr:   ts.ErrDiffRefs,
		},
		{
			name:      "insts, same refs, no args", // int <: int
			subType:   h.Inst("int"),
			superType: h.Inst("int"),
			wantErr:   nil,
		},
		// args count
		{
			name:      "insts, arg has less args", // vec <: vec<int>
			subType:   h.Inst("vec"),
			superType: h.Inst("vec", h.Inst("int")),
			wantErr:   ts.ErrArgsCount,
		},
		{
			name:      "insts, arg has same args count", // vec<int> <: vec<int>
			subType:   h.Inst("vec", h.Inst("int")),
			superType: h.Inst("vec", h.Inst("int")),
			wantErr:   nil,
		},
		{
			name:      "insts, arg has more args count", // vec<int, str> <: vec<int>
			subType:   h.Inst("vec", h.Inst("int"), h.Inst("str")),
			superType: h.Inst("vec", h.Inst("int")),
			wantErr:   nil,
		},
		// args compatibility
		{
			name:    "insts, one arg's arg incompat", // vec<str> <: vec<int|str>
			subType: h.Inst("vec", h.Inst("str")),
			superType: h.Inst(
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
			subType: h.Inst(
				"vec",
				h.Union(
					h.Inst("str"),
					h.Inst("int"),
				),
			),
			superType: h.Inst("vec", h.Inst("int")),
			wantErr:   ts.ErrArgNotSubtype,
		},
		// arr
		{
			name: "expr and constr has diff lit types (constr not union)",
			subType: ts.Expr{
				Lit: ts.LitExpr{Enum: []string{}},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{Arr: &ts.ArrLit{}},
			},
			wantErr: ts.ErrDiffLitTypes,
		},
		{
			name: "expr's arr lit has lesser size than constr",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 1}, // expr doesn't matter here
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{Size: 2},
				},
			},
			wantErr: ts.ErrLitArrSize,
		},
		{
			name: "expr's arr has incompat type",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{
						Size: 2, // same size itself won't cause any problem
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "a"}, // type will
						},
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Arr: &ts.ArrLit{
						Size: 3, // bigger size won't cause any problem
						Expr: ts.Expr{
							Inst: ts.InstExpr{Ref: "a"}, // same type
						},
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a", "b"},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a"},
				},
			},
			wantErr: ts.ErrBigEnum,
		},
		{
			name: "expr and constr enums, expr not bigger but contain diff el",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a", "d"}, // d != b
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a", "b", "c"},
				},
			},
			wantErr: ts.ErrEnumEl,
		},
		{
			name: "expr and constr enums, expr not bigger and all reqired els are the same",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a", "b"},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Enum: []string{"a", "b", "c"}, // c el won't cause any problem
				},
			},
			wantErr: nil,
		},
		// rec
		{
			name: "expr and constr recs, expr has less fields",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{}, // 0 fields is ok
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{
						"a": {}, // expr itself doesn't matter here
					},
				},
			},
			wantErr: ts.ErrRecLen,
		},
		{
			name: "expr and constr recs, expr leaks field",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"b": {}, // expr itself doesn't matter here
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{
						"a": {}, // but this field is missing
					},
				},
			},
			wantErr: ts.ErrRecNoField,
		},
		{
			name: "expr and constr recs, expr has incompat field",
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"b": {}, // b field itself won't cause any problems
						"a": {}, // this one will
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Rec: map[string]ts.Expr{ // both has 1 field
						"a": {
							Inst: ts.InstExpr{Ref: "x"}, // not same as in expr
						},
						"b": {}, // b field itself won't cause any problems
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Inst: ts.InstExpr{Ref: "x"}, // not compat with both a and b
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Inst: ts.InstExpr{Ref: "b"},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
						{Inst: ts.InstExpr{Ref: "c"}},
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "c"}},
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "x"}}, // this
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
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
			subType: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "c"}},
						{Inst: ts.InstExpr{Ref: "a"}},
					},
				},
			},
			superType: ts.Expr{
				Lit: ts.LitExpr{
					Union: []ts.Expr{
						{Inst: ts.InstExpr{Ref: "a"}},
						{Inst: ts.InstExpr{Ref: "b"}},
						{Inst: ts.InstExpr{Ref: "c"}},
					},
				},
			},
			wantErr: nil,
		},
		{ // vec<t1> t1, vec<t2> t2, { vec, t1=vec<t1>, t2=vec<t2> }
			enabled:   true,
			name:      "arg and constr contain diff refs but they are actually recursive",
			subType:   h.Inst("vec", h.Inst("t1")),
			trace1:    ts.NewTrace(nil, "t1"),
			superType: h.Inst("vec", h.Inst("t2")),
			trace2:    ts.NewTrace(nil, "t2"),
			scope: map[string]ts.Def{
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"t2":  h.Def(h.Inst("vec", h.Inst("t2"))),
			},
			checker: func(mcmr *MockrecursionCheckerMockRecorder) {
				mcmr.Check(gomock.Any(), gomock.Any()).Return(false, nil).AnyTimes()
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		if !tt.enabled {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			recChecker := NewMockrecursionChecker(ctrl)
			if tt.checker != nil {
				tt.checker(recChecker.EXPECT())
			}

			checker := ts.NewSubtypeChecker(recChecker)

			require.ErrorIs(
				t,
				checker.Check(tt.subType, tt.trace1, tt.superType, tt.trace2, tt.scope),
				tt.wantErr,
			)
		})
	}
}
