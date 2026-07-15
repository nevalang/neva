package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Split must preserve each branch's stream lifecycle and item order.
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(t, "{\"else\": [1, 3], \"then\": [2, 4]}\n", out)
}
