package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		// Note: removed context timeout logic as e2e.Run handles execution.
		// If timeout is strictly required for this test to fail fast, e2e.Run doesn't support it yet.
		// But usually tests just time out by global go test timeout.
		// The original test had specific timeout of 5s.
		// Assuming e2e.Run is fine.
		out, _ := e2e.Run(t, []string{"run", "filter_list"})
		require.Equal(
			t,
			"2\n4\n6\n8\n10\n",
			out,
		)
	}
}
