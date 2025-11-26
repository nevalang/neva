package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, err := e2e.RunExpectingError(t, "run", "main")
	require.Equal(t, "", out) // Expecting error usually means empty stdout for compiler errors, but let's check. Original code ignored stdout.
	require.Contains(
		t,
		string(out)+err, // The original test checked combined output. e2e.RunExpectingError returns stdout and stderr separately.
		"main/main.neva:4:1: array inport 'printf:args' is used incorrectly: slot 1 is missing\n",
	)
}
