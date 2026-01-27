package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run("", func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "main"})
			require.Equal(
				t,
				"1\n0\n",
				out,
				"iteration: %d", i,
			)
		})
	}
}
