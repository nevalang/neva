package test

import (
	"os"
	"os/exec"
	"testing"

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

	cmd := exec.Command("neva", "run", "const_refs")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		`{"d": {"key": 1}, "l": [1, 2, 3]}
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
