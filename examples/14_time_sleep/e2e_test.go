package test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("../")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "14_time_sleep")

	before := time.Now()
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(t, "", string(out))
	require.Greater(t, time.Since(before).Seconds(), float64(1))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
