// Package parser implements source code parsing.
package parser

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseFile(t *testing.T) {
	tests := []struct {
		name    string
		bb      []byte
		want    src.File
		wantErr error
	}{
		{
			name: "only_use_statement_but_with_dot",
			bb: []byte(`
				use {
					std/tmp
					github.com/nevalang/neva/pkg/typesystem
					some/really/deeply/nested/path/to/local/package/in/the/project
				}
			`),
			want: src.File{
				Imports: map[string]string{
					"tmp":        "std/tmp",
					"typesystem": "github.com/nevalang/neva/pkg/typesystem",
					"project":    "some/really/deeply/nested/path/to/local/package/in/the/project",
				},
			},
		},
	}

	p := Parser{
		debug: false,
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := p.ParseFile(context.Background(), tt.bb)
			require.Equal(t, got, tt.want)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
