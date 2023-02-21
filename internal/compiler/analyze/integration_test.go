//go:build integration
// +build integration

package analyze_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emil14/neva/internal/compiler/analyze"
	"github.com/emil14/neva/internal/compiler/src"
)

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		prog    src.Prog
		wantErr error
	}

	tests := []testcase{
		{
			name:    "",
			prog:    src.Prog{},
			wantErr: nil,
		},
	}

	a := analyze.Analyzer{
		Resolver: TestDefaultResolver,
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := a.Analyze(context.Background(), tt.prog)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
