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

	cmd := exec.Command("neva", "run", "stream_product")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		`{"first": 0, "second": 0}
{"first": 0, "second": 1}
{"first": 0, "second": 2}
{"first": 1, "second": 0}
{"first": 1, "second": 1}
{"first": 1, "second": 2}
{"first": 2, "second": 0}
{"first": 2, "second": 1}
{"first": 2, "second": 2}
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
