package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we do 100 attempts because there was a floating bug
	// caused by unordered map that was
	for i := 0; i < 100; i++ {
		out := e2e.Run(t, "run", "main")
		require.Equal(
			t,
			"-6\n",
			out,
		)
	}
}
