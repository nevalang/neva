package test

// in this file we test files designed specifically for e2e.

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(
		t,
		"42\n",
		out,
	)
}
