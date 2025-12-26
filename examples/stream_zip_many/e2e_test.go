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

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			out := e2e.Run(t, []string{"run", "stream_zip_many"}, e2e.WithStderr())
			require.Equal(t, expectedOutput, out)
		})
	}
}
