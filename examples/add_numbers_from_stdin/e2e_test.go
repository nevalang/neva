package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run("Add numbers from stdin", func(t *testing.T) {
			out := e2e.RunExampleWithStdin(t, "3\n4\n\n", "add_numbers_from_stdin")
			require.Equal(
				t,
				"7\n",
				out,
			)
		})
	}
}
