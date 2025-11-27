package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.Run(t, "run", "main")
	require.Equal(
		t,
		"50\n30\n20\n100\n",
		out,
	)
}
