package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Regression: union literal payloads must be checked before generic receivers erase the concrete tag type.
	_, stderr := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1))
	require.Contains(t, stderr, "Union literal payload type")
}
