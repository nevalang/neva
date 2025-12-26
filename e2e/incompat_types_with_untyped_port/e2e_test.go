package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, err := e2e.RunExpectingError(t, "run", "main")
	require.Contains(
		t,
		out+err,
		"Incompatible types: in:data -> println: subtype instance must have same ref as supertype: got any, want int",
	)
}
