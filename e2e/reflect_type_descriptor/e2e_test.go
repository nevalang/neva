package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Test verifies that users can construct finite descriptors, including a
// recursive back-edge, using the public std/reflect types.
func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Empty(t, out)
}
