package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestTrace_String(t *testing.T) {
	tests := []struct {
		trace func() ts.Trace
		want  string
	}{
		{
			trace: func() ts.Trace {
				t1 := ts.NewTrace(nil, "t1")
				t2 := ts.NewTrace(&t1, "t2")
				return ts.NewTrace(&t2, "t3")
			},
			want: "[t1, t2, t3]",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.want, func(t *testing.T) {
			got := tt.trace().String()
			assert.Equal(t, tt.want, got)
		})
	}
}
