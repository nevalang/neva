package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		out, _ := e2e.Run(t, []string{"run", "map_list"})
		require.Equal(
			t,
			"[49,29,19,99]\n",
			out,
		)
	}
}
