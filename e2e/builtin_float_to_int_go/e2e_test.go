package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})

	// Float-to-int cast should truncate toward zero, like Go.
	require.Equal(t, "1\n", out)
}
