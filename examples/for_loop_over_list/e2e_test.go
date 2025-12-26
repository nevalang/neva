package test

import (
	"fmt"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// We do 100 attempts to prove that implementation is correct
	// and order of elements is always the same.
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "for_loop_over_list"})
			require.Equal(
				t,
				"1\n2\n3\n",
				out,
				"Unexpected output",
			)
		})
	}
}
