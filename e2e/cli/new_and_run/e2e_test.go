package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// This function tests `neva new` followed by `neva run`.
func Test(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
	}()

	e2e.Run(t, "new", ".")

	out := e2e.RunCombined(t, "run", "src")
	require.Equal(t, "Hello, World!\n", out)
}
