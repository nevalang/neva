package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		out := e2e.RunExampleCombined(t, "advanced_error_handling")
		require.Equal(
			t,
			`panic: {"text": "Get \"definitely%20not%20a%20valid%20URL\":  unsupported protocol scheme \"\""}
`,
			out,
		)
	}
}
