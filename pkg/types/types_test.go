package types_test

import (
	"testing"

	"github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestLiteralExpr_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lit  types.LitExpr
		want bool
	}{
		{
			name: "all 4 fields: arr, enum, union and rec are empty",
			lit:  types.LitExpr{nil, nil, nil, nil},
			want: true,
		},
		{
			name: "arr not empty",
			lit:  types.LitExpr{&types.ArrLit{}, nil, nil, nil},
			want: false,
		},
		{
			name: "rec not empty",
			lit:  types.LitExpr{nil, map[string]types.Expr{}, nil, nil},
			want: false,
		},
		{
			name: "enum not empty",
			lit:  types.LitExpr{nil, nil, []string{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  types.LitExpr{nil, nil, nil, []types.Expr{}},
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
		lit  types.LitExpr
		want types.LiteralType
	}{
		{
			name: "unknown",
			lit:  types.LitExpr{nil, nil, nil, nil},
			want: types.EmptyLitType,
		},
		{
			name: "arr",
			lit:  types.LitExpr{&types.ArrLit{}, nil, nil, nil},
			want: types.ArrLitType,
		},
		{
			name: "rec",
			lit:  types.LitExpr{nil, map[string]types.Expr{}, nil, nil},
			want: types.RecLitType,
		},
		{
			name: "enum",
			lit:  types.LitExpr{nil, nil, []string{}, nil},
			want: types.EnumLitType,
		},
		{
			name: "union",
			lit:  types.LitExpr{nil, nil, nil, []types.Expr{}},
			want: types.UnionLitType,
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
		inst types.InstExpr
		want bool
	}{
		{
			name: "default inst (empty ref and nil args)",
			inst: types.InstExpr{
				Ref:  "",
				Args: nil,
			},
			want: true,
		},
		{
			name: "empty ref and empty list args",
			inst: types.InstExpr{
				Ref:  "",
				Args: []types.Expr{},
			},
			want: true,
		},
		{
			name: "empty ref and non empty list args",
			inst: types.InstExpr{
				Ref:  "",
				Args: []types.Expr{{}}, // content doesn't matter here
			},
			want: false,
		},
		{
			name: "non-empty ref and non empty list args",
			inst: types.InstExpr{
				Ref:  "t",
				Args: []types.Expr{{}},
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
