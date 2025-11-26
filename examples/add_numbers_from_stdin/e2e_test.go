package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	for i := 0; i < 1; i++ {
		t.Run("Add numbers from stdin", func(t *testing.T) {
			out := e2e.RunWithStdin(t, "3\n4\n\n", "run", "add_numbers_from_stdin")
			require.Equal(
				t,
				"7\n",
				out,
			)
		})
	}
}
