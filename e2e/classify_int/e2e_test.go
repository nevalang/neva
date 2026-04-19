package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	stdout, _ := e2e.Run(t, []string{"run", "src"})
	require.Equal(
		t,
		stdout,
		"negative :(\n",
	)
}
