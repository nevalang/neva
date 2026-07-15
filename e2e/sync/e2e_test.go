package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Sync must wait for the delayed handler before it releases the next event.
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(t, "1\n2\n3\n", out)
}
