package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/require"
)

var h ts.Helper

func TestLiteralExpr_Empty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		lit  ts.LitExpr
		want bool
	}{
		{
			name: "all 4 fields: arr, enum, union and rec are empty",
			lit:  ts.LitExpr{nil, nil, nil, nil},
			want: true,
		},
		{
			name: "arr not empty",
			lit:  ts.LitExpr{&ts.ArrLit{}, nil, nil, nil},
			want: false,
		},
		{
			name: "rec not empty",
			lit:  ts.LitExpr{nil, map[string]ts.Expr{}, nil, nil},
			want: false,
		},
		{
			name: "enum not empty",
			lit:  ts.LitExpr{nil, nil, []string{}, nil},
			want: false,
		},
		{
			name: "union not empty",
			lit:  ts.LitExpr{nil, nil, nil, []ts.Expr{}},
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
		lit  ts.LitExpr
		want ts.LiteralType
	}{
		{
			name: "unknown",
			lit:  ts.LitExpr{nil, nil, nil, nil},
			want: ts.EmptyLitType,
		},
		{
			name: "arr",
			lit:  ts.LitExpr{&ts.ArrLit{}, nil, nil, nil},
			want: ts.ArrLitType,
		},
		{
			name: "rec",
			lit:  ts.LitExpr{nil, map[string]ts.Expr{}, nil, nil},
			want: ts.RecLitType,
		},
		{
			name: "enum",
			lit:  ts.LitExpr{nil, nil, []string{}, nil},
			want: ts.EnumLitType,
		},
		{
			name: "union",
			lit:  ts.LitExpr{nil, nil, nil, []ts.Expr{}},
			want: ts.UnionLitType,
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
		inst ts.InstExpr
		want bool
	}{
		{
			name: "default inst (empty ref and nil args)",
			inst: ts.InstExpr{
				Ref:  "",
				Args: nil,
			},
			want: true,
		},
		{
			name: "empty ref and empty list args",
			inst: ts.InstExpr{
				Ref:  "",
				Args: []ts.Expr{},
			},
			want: true,
		},
		{
			name: "empty ref and non empty list args",
			inst: ts.InstExpr{
				Ref:  "",
				Args: []ts.Expr{{}}, // content doesn't matter here
			},
			want: false,
		},
		{
			name: "non-empty ref and non empty list args",
			inst: ts.InstExpr{
				Ref:  "t",
				Args: []ts.Expr{{}},
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

func TestDef_String(t *testing.T) {
	tests := []struct {
		name string
		def  ts.Def
		want string
	}{
		{
			name: "",
			def: h.Def(
				h.Inst("vec", h.Inst("T")),
				h.Param("T", h.Inst("int")),
			),
			want: "<T int> = vec<T>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.def.String(); got != tt.want {
				t.Errorf("Def.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
