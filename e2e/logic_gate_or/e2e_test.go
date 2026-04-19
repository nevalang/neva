package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that logical OR works
func TestOR(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(
		t,
		"false\n",
		out,
	)
}
