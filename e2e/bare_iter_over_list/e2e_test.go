package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		`{"data": 50, "idx": 0, "last": false}
{"data": 30, "idx": 1, "last": false}
{"data": 20, "idx": 2, "last": false}
{"data": 100, "idx": 3, "last": true}
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
