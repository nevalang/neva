//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for range 10 {
		out, _ := e2e.Run(t, []string{"run", "."})
		require.Equal(
			t,
			"3\n",
			out,
		)
	}
}
