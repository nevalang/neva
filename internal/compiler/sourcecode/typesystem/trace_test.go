package typesystem_test

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/stretchr/testify/assert"
)

func TestTrace_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		trace func() ts.Trace
		want  string
	}{
		{
			trace: func() ts.Trace {
				t1 := ts.NewTrace(nil, ts.DefaultStringer("t1"))
				t2 := ts.NewTrace(&t1, ts.DefaultStringer("t2"))
				return ts.NewTrace(&t2, ts.DefaultStringer("t3"))
			},
			want: "[t1, t2, t3]",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			got := tt.trace().String()
			assert.Equal(t, tt.want, got)
		})
	}
}
