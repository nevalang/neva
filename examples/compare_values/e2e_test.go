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

	for i := 0; i < 1; i++ {
		out := e2e.Run(t, "run", "compare_values")
		require.Equal(
			t,
			"They match\n",
			out,
		)
	}
}
