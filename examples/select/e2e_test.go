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

	for i := 0; i < 100; i++ {
		t.Logf("Running iteration %d", i)
		out := e2e.Run(t, "run", "select")
		require.Equal(
			t,
			"a\nb\nc\nd\n",
			out,
			"iteration %d failed\n", i,
		)
	}
}
