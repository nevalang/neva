package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that Len def works with the list of integers.
func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "list_len")

	require.Equal(
		t,
		"5\n",
		out,
	)
}
