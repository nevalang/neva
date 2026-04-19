package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	_, stderr := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1))
	require.Contains(
		t,
		stderr,
		"Incompatible types: in:data -> println: subtype instance must have same ref as supertype: got any, want int",
	)
}
