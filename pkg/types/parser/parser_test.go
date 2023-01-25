package parser_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/emil14/neva/pkg/types/parser"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		s       string
		want    ts.Expr
		wantErr error
	}{
		{
			s:    "{}",
			want: h.Rec(nil),
		},
		{
			s:    "[256]u8",
			want: h.Arr(256, h.Inst("u8")),
		},
		{
			s:    "t",
			want: h.Inst("t"),
		},
		{
			s:       "t<>",
			wantErr: parser.ErrEmptyAngleBrackets,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			got, err := parser.Parse(tt.s)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
