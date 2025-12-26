package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		out := e2e.Run(t, []string{"run", "reduce_list"})
		require.Equal(
			t,
			"55\n",
			out,
		)
	}
}
