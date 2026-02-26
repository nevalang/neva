package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})

	// String-to-bytes cast should preserve byte content through roundtrip.
	require.Equal(t, "Hello, bytes!\n", out)
}
