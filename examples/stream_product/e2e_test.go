package test

import (
	"fmt"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

var expectedOutput = `{"first": 0, "second": 0}
{"first": 0, "second": 1}
{"first": 0, "second": 2}
{"first": 0, "second": 3}
{"first": 1, "second": 0}
{"first": 1, "second": 1}
{"first": 1, "second": 2}
{"first": 1, "second": 3}
{"first": 2, "second": 0}
{"first": 2, "second": 1}
{"first": 2, "second": 2}
{"first": 2, "second": 3}
{"first": 3, "second": 0}
{"first": 3, "second": 1}
{"first": 3, "second": 2}
{"first": 3, "second": 3}
`

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run(fmt.Sprintf("Run %d", i+1), func(t *testing.T) {
			out, _ := e2e.Run(t, []string{"run", "stream_product"})
			require.Equal(t, expectedOutput, out)
		})
	}
}
