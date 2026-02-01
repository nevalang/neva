package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we do 10 iterations because there was a bug
	// that was only reproducible after a few runs
	for i := 0; i < 10; i++ {
		out, _ := e2e.Run(t, []string{"run", "main"})
		require.Equal(
			t,
			"[49,29,19,99]\n",
			out,
			"iteration %d", i,
		)
	}
}
