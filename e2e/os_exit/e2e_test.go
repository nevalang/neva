package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	_, stderr := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1))
	require.Contains(t, stderr, "failed to run generated executable: exit status 2")
	require.NotContains(t, stderr, "exit cause dataflow trace")
	require.NotContains(t, stderr, "panic cause dataflow trace")
	require.NotContains(t, stderr, "exit: 2")
}
