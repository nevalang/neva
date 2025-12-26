package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that Len def works with the list of integers.
func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "list_len"})

	require.Equal(
		t,
		"5\n",
		out,
	)
}
