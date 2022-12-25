package types_test

import (
	"errors"
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResolver_Resolve(t *testing.T) { //nolint:maintidx
	type testcase struct {
		name              string
		expr              ts.Expr
		scope             map[string]ts.Def
		initValidatorMock func(v *Mockvalidator)
		initCheckerMock   func(c *Mockchecker)
		want              ts.Expr
		wantErr           error
	}

	tests := []func() testcase{
		func() testcase {
			return testcase{
				name:              "invalid expr",
				initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(ts.Expr{}).Return(errors.New("")) },
				wantErr:           ts.ErrInvalidExpr,
			}
		},
		func() testcase { // expr = int, scope = {}
			expr := h.Inst("int")
			return testcase{
				name:              "inst expr refers to undefined",
				expr:              expr,
				initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(expr).Return(nil) },
				scope:             map[string]ts.Def{},
				wantErr:           ts.ErrUndefinedRef,
			}
		},
		func() testcase { // expr = vec<>, scope = { vec<t> = vec }
			expr := h.Inst("vec")
			return testcase{
				name:              "args < params",
				expr:              expr,
				initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(expr).Return(nil) },
				scope: map[string]ts.Def{
					"vec": h.NativeDef("vec", ts.Param{Name: "t"}),
				},
				wantErr: ts.ErrInstArgsLen,
			}
		},
		func() testcase { // expr = map<t1>, scope = { map<t t2> = map, t1 , t2 }
			expr := h.Inst("map", h.Inst("t1"))
			constr := h.Inst("t2")
			return testcase{
				name: "incompat arg",
				expr: expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(expr.Inst.Args[0]).Return(nil) // first recursive call
					v.EXPECT().Validate(constr).Return(nil)            // first recursive call
				},
				initCheckerMock: func(c *Mockchecker) {
					c.EXPECT().SubtypeCheck(expr.Inst.Args[0], constr).Return(errors.New(""))
				},
				scope: map[string]ts.Def{
					"map": h.NativeDef(
						"map",
						ts.Param{"t", constr},
					),
					"t1": h.NativeDef("t1"),
					"t2": h.NativeDef("t2"),
				},
				wantErr: ts.ErrIncompatArg,
			}
		},
		func() testcase { // expr = t1<t2>, scope = { t2, t1<t> = t3 }
			expr := h.Inst("t1", h.Inst("t2"))
			return testcase{
				name: "expr base type not found",
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.NativeDef("t3", ts.Param{Name: "t"}),
					"t2": h.NativeDef("t2"),
				},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(expr.Inst.Args[0]).Return(nil)
				},
				wantErr: ts.ErrNoBaseType,
			}
		},
		func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			return testcase{
				name: "constr undefined ref",
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.NativeDef("t3", ts.Param{Name: "t", Constraint: constr}),
					"t2": h.NativeDef("t2"),
				},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(expr.Inst.Args[0]).Return(nil)
					v.EXPECT().Validate(constr).Return(nil)
				},
				wantErr: ts.ErrConstr,
			}
		},
		func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1 }
			expr := h.Inst("t1", h.Inst("t2"))
			return testcase{
				name: "constr base type not found",
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.NativeDef("t3", ts.Param{Name: "t"}),
					"t2": h.NativeDef("t2"),
				},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(expr.Inst.Args[0]).Return(nil)
				},
				wantErr: ts.ErrNoBaseType,
			}
		},
		func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1, t3 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			return testcase{
				name: "invalid constr",
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.NativeDef("t3", ts.Param{Name: "t"}),
					"t2": h.NativeDef("t2"),
					"t3": h.NativeDef("t2"),
				},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(expr.Inst.Args[0]).Return(nil)
					v.EXPECT().Validate(constr).Return(errors.New(""))
				},
				wantErr: ts.ErrInvalidExpr,
			}
		},
		// Lits
		func() testcase { // expr = enum{}, scope = {}
			expr := h.Enum()
			return testcase{
				name:              "enum",
				expr:              expr,
				initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(expr).Return(nil) },
				want:              expr,
				wantErr:           nil,
			}
		},
		func() testcase { // expr = [2]t, scope = {}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				name:  "arr with unresolvable (undefined) type",
				scope: map[string]ts.Def{},
				expr:  expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(typ).Return(nil)
				},
				wantErr: ts.ErrArrType,
			}
		},
		func() testcase { // expr = [2]t, scope = {}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				name:  "arr with unresolvable (invalid) type",
				scope: map[string]ts.Def{"t": h.NativeDef("t")},
				expr:  expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(typ).Return(errors.New(""))
				},
				wantErr: ts.ErrArrType,
			}
		},
		func() testcase { // expr = [2]t, scope = {t=t}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				name:  "arr with resolvable type",
				scope: map[string]ts.Def{"t": h.NativeDef("t")},
				expr:  expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(typ).Return(nil)
				},
				want:    expr,
				wantErr: nil,
			}
		},
		func() testcase { // expr = t1 | t2, scope = {t1=t1}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			return testcase{
				name:  "union with unresolvable (undefined) element",
				scope: map[string]ts.Def{"t1": h.NativeDef("t1")},
				expr:  expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(t1).Return(nil)
					v.EXPECT().Validate(t2).Return(nil)
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
		func() testcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			return testcase{
				name: "union with unresolvable (invalid) element",
				scope: map[string]ts.Def{
					"t1": h.NativeDef("t1"),
					"t2": h.NativeDef("t2"),
				},
				expr: expr,
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
					v.EXPECT().Validate(t1).Return(nil)
					v.EXPECT().Validate(t2).Return(errors.New(""))
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
	}

	for _, tt := range tests {
		tt := tt
		tc := tt()

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			v := NewMockvalidator(ctrl)
			if tc.initValidatorMock != nil {
				tc.initValidatorMock(v)
			}

			c := NewMockchecker(ctrl)
			if tc.initCheckerMock != nil {
				tc.initCheckerMock(c)
			}

			got, err := ts.MustNewResolver(v, c).Resolve(tc.expr, tc.scope)
			require.Equal(t, got, tc.want)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
