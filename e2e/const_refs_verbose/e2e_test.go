package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// Nested map has only one key because keys are unordered
// so having order in test will make it flacky.
func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		`{"l": [1, 2, 3], "d": {"key": 1}}
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
