package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		out, _ := e2e.Run(t, []string{"run", "main"})
		require.Equal(
			t,
			"[0,1,2,3,4,5,6,7,8,9]\n",
			out,
			"iteration: %d", i,
		)
	}
}
