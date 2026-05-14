package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	_, stderr := e2e.Run(t, []string{"run", "."}, e2e.WithCode(1))
	require.Contains(
		t,
		stderr,
		`panic: {"child": {"tag": "None"}, "text": "Get \"definitely%20not%20a%20valid%20URL\": unsupported protocol scheme \"\""}`,
	)
	require.Contains(t, stderr, "panic cause dataflow trace\n")
}
