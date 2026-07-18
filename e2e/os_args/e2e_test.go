package test

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Intent: verify std/os.Args exposes the runtime argv list.
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.True(t, strings.HasPrefix(out, "[\""))
	require.False(t, strings.Contains(out, ","))
	require.True(t, strings.HasSuffix(out, "\n"))
}
