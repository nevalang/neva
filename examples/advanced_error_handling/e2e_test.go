package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	_, stderr := e2e.Run(t, []string{"run", "."})
	require.Equal(
		t,
		`panic: {"text": "Get \"definitely%20not%20a%20valid%20URL\": unsupported protocol scheme \"\""}
`,
		stderr,
	)
}
