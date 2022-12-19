package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResolver_Resolve(t *testing.T) {
	tests := []struct {
		name string

		expr  types.Expr
		scope map[string]types.Def

		initValidatorMock func(v *Mockvalidator)
		initCheckerMock   func(c *Mockchecker)

		want    types.Expr
		wantErr error
	}{
		{
			name:              "",
			expr:              types.Expr{},
			scope:             map[string]types.Def{},
			initValidatorMock: func(v *Mockvalidator) {},
			want:              types.Expr{},
			wantErr:           nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			v := NewMockvalidator(ctrl)
			tt.initValidatorMock(v)
			c := NewMockchecker(ctrl)
			tt.initCheckerMock(c)

			got, err := types.CustomResolver(v, c).Resolve(tt.expr, tt.scope)
			require.Equal(t, got, tt.wantErr)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
