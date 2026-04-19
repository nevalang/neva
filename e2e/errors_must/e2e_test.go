package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we run N times to make sure https://github.com/nevalang/neva/issues/872 is fixed
	for range 10 {
		out, _ := e2e.Run(t, []string{"run", "main"})
		require.Equal(
			t,
			"success!\n",
			out,
		)
	}
}
