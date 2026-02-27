package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})

	// Int-to-float cast should preserve numeric value.
	require.Equal(t, "42\n", out)
}
