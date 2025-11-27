package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Nested map has only one key because keys are unordered
// so having order in test will make it flacky.
func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "const_refs")
	require.Equal(
		t,
		`{"d": {"key": 1}, "l": [1, 2, 3]}
`,
		out,
	)
}
