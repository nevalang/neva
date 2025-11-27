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
	out, err := e2e.RunExpectingError(t, "run", "main")
	require.Contains(
		t,
		out+err,
		"main/main.neva:12:1: port 'in:start' is used twice\n",
	)
}
