package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "advanced_error_handling")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		`panic: {"text": "Get \"definitely%20not%20a%20valid%20URL\":  unsupported protocol scheme \"\""}
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
