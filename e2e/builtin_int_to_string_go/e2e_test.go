package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})

	// Int-to-string cast should follow Go's Unicode code-point conversion.
	require.Equal(t, "*\n", out)
}
