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

	// Test successful case with "Alice"
	out := e2e.RunWithStdin(t, "Alice\n", "run", "switch")
	require.Equal(t, "Enter the name: ALICE\n", out)

	// Test panic case with "Bob"
	out = e2e.RunWithStdin(t, "Bob\n", "run", "switch")
	require.Equal(t, "Enter the name: bob\n", out)

	// Test panic case with "Charlie"
	out = e2e.RunWithStdin(t, "Charlie\n", "run", "switch")
	require.Equal(t, "Enter the name: panic: Charlie\n", out)
}
