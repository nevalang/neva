package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we run N times to make sure https://github.com/nevalang/neva/issues/872 is fixed
	for range 10 {
		cmd := exec.Command("neva", "run", "main")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
		require.Equal(
			t,
			"success!\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
