package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run("Add numbers from stdin", func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "main"}, e2e.WithStdin("3\n4\n\n"))
			require.Equal(
				t,
				"7\n",
				out,
			)
		})
	}
}
