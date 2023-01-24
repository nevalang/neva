package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/emil14/neva/pkg/types/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    ts.Expr
		wantErr error
	}{
		{
			s:       "t",
			want:    h.Inst("t"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			got, err := parser.Parse(tt.s)
			require.Equal(t, got, tt.want)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
