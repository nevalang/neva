package test

import (
	"fmt"
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

	// We do 100 attempts to prove that implementation is correct
	// and order of elements is always the same.
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			out := e2e.Run(t, "run", "for_loop_over_list")
			require.Equal(
				t,
				"1\n2\n3\n",
				out,
				"Unexpected output",
			)
		})
	}
}
