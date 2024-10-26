package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	cmd.Stdin = strings.NewReader("yo\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"main/main.neva:4:19: port 'out:stop' is used twice\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}