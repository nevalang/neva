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
		"{}\n",
		out,
	)
}
