package test

import (
	"fmt"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	expectedOutput := `[1,10,100]
[2,20,200]
[3,30,300]
`

	for i := range 10 {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "."})
			require.Equal(t, expectedOutput, out)
		})
	}
}
