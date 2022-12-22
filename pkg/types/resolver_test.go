package types_test

import (
	"errors"
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResolver_Resolve(t *testing.T) {
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
		// func() testcase {
		// 	return testcase{
		// 		name:              "invalid expr",
		// 		initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(ts.Expr{}).Return(errors.New("")) },
		// 		wantErr:           ts.ErrInvalidExpr,
		// 	}
		// },
		// func() testcase {
		// 	expr := h.Inst("int")
		// 	return testcase{
		// 		name:              "inst ref type not in scope",
		// 		expr:              expr,
		// 		initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(expr).Return(nil) },
		// 		scope:             map[string]ts.Def{},
		// 		wantErr:           ts.ErrNoRefType,
		// 	}
		// },
		// func() testcase { // expr = vec<>, scope = { vec<t> }
		// 	expr := h.Inst("vec")
		// 	return testcase{
		// 		name:              "inst args < ref type params",
		// 		expr:              expr,
		// 		initValidatorMock: func(v *Mockvalidator) { v.EXPECT().Validate(expr).Return(nil) },
		// 		scope: map[string]ts.Def{
		// 			"vec": h.NativeDef("vec", ts.Param{Name: "t"}),
		// 		},
		// 		wantErr: ts.ErrInstArgsLen,
		// 	}
		// },
		func() testcase { // expr = map<t1>, scope = { map<t t2>, t1 , t2 }
			expr := h.Inst("map", h.Inst("t1"))
			constr := h.Inst("t2")
			return testcase{
				name: "inst arg has incompat arg",
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
				wantErr: ts.ErrSubtype,
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
