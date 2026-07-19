package test

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Intent: verify filesystem creation, inspection, mutation, and cleanup.
	out, _ := e2e.Run(t, []string{"run", "main"})

	trimmed := strings.TrimSpace(out)
	require.Contains(t, trimmed, "before=[")
	require.Contains(t, trimmed, "size0=0")
	require.Contains(t, trimmed, "size1=5")
	require.Contains(t, trimmed, "after=[]")
}
