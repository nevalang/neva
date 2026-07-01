package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Filter must close the output stream when the final input item is rejected.
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(t, "[2,4,6,8]\n", out)
}
