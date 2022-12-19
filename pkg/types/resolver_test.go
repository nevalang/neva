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
		func() testcase {
			return testcase{
				name: "invalid expr",
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(ts.Expr{}).Return(errors.New(""))
				},
				wantErr: ts.ErrInvalidExpr,
			}
		},
		func() testcase {
			expr := h.InstExpr("int")
			return testcase{
				name:  "inst ref type not in scope",
				expr:  expr,
				scope: map[string]ts.Def{},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
				},
				wantErr: ts.ErrNoRefType,
			}
		},
		func() testcase {
			expr := h.InstExpr("int")
			return testcase{
				name: "inst args < ref type params",
				expr: expr,
				scope: map[string]ts.Def{
					"int": h.NativeType("int", ts.Param{}),
				},
				initValidatorMock: func(v *Mockvalidator) {
					v.EXPECT().Validate(expr).Return(nil)
				},
				wantErr: ts.ErrInstArgsLen,
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
