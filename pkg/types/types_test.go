package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestLiteralExpr_Empty(t *testing.T) {
	tests := []struct {
		name string
		lit  types.LiteralExpr
		want bool
	}{
		{
			name: "all 4 fields: arr, enum, union and rec are empty",
			lit:  types.LiteralExpr{nil, nil, nil, nil},
			want: true,
		},
		{
			name: "arr not empty",
			lit:  types.LiteralExpr{&types.ArrLit{}, nil, nil, nil},
			want: false,
		},
		{
			name: "rec not empty",
			lit:  types.LiteralExpr{nil, map[string]types.Expr{}, nil, nil},
			want: false,
		},
		{
			name: "enum not empty",
			lit:  types.LiteralExpr{nil, nil, []string{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  types.LiteralExpr{nil, nil, nil, []types.Expr{}},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.lit.Empty(), tt.want)
		})
	}
}

func TestLiteralExpr_Type(t *testing.T) {
	tests := []struct {
		name string
		lit  types.LiteralExpr
		want types.LiteralType
	}{
		{
			name: "unknown",
			lit:  types.LiteralExpr{nil, nil, nil, nil},
			want: types.EmptyLitType,
		},
		{
			name: "arr",
			lit:  types.LiteralExpr{&types.ArrLit{}, nil, nil, nil},
			want: types.ArrLitType,
		},
		{
			name: "rec",
			lit:  types.LiteralExpr{nil, map[string]types.Expr{}, nil, nil},
			want: types.RecLitType,
		},
		{
			name: "enum",
			lit:  types.LiteralExpr{nil, nil, []string{}, nil},
			want: types.EnumLitType,
		},
		{
			name: "union",
			lit:  types.LiteralExpr{nil, nil, nil, []types.Expr{}},
			want: types.UnionLitType,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.lit.Type(), tt.want)
		})
	}
}
