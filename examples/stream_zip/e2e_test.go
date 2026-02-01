package test

import (
	"fmt"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	expectedOutput := `{"left": 0, "right": "a"}
{"left": 1, "right": "b"}
{"left": 2, "right": "c"}
`

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "."})
			require.Equal(t, expectedOutput, out)
		})
	}
}
