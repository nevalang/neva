package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.RunWithStdin(t, "yo\n", "run", "main")
	require.Equal(
		t,
		"yo\n",
		out,
	)
}
