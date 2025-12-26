package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 10; i++ {
		out, _ := e2e.Run(t, []string{"run", "add_numbers"})
		require.Equal(
			t,
			"3\n",
			out,
		)
	}
}
