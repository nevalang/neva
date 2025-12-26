package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Note: previous test used stdin "yo\n" but asserted compiler error?
	// The error "port 'in:start' is used twice" is a compiler error.
	// Compiler doesn't read stdin. The stdin was likely ignored or copy-pasted.
	// I will use RunExpectingError.
	out := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1), e2e.WithStderr())
	require.Contains(
		t,
		out,
		"main/main.neva:12:1: port 'in:start' is used twice\n",
	)
}
