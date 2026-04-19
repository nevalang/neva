package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := range 100 {
		t.Logf("Running iteration %d", i)
		out, _ := e2e.Run(t, []string{"run", "."})
		require.Equal(
			t,
			"a\nb\nc\nd\n",
			out,
			"iteration %d failed\n", i,
		)
	}
}
