package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	for i := 0; i < 1; i++ {
		out := e2e.Run(t, "run", "advanced_error_handling")
		require.Equal(
			t,
			`panic: {"text": "Get \"definitely%20not%20a%20valid%20URL\":  unsupported protocol scheme \"\""}
`,
			out,
		)
	}
}
