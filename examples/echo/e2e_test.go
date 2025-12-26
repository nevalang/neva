package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run("Echo Test", func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "echo"}, e2e.WithStdin("yo\n"))
			require.Equal(
				t,
				"yo\n",
				out,
			)
		})
	}
}
