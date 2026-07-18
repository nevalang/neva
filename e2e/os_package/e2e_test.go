package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Intent: verify env access, expansion, and lookup behavior end-to-end.
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(
		t,
		"value\nfalse\nvalue\n",
		out,
	)
}
