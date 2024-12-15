package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkMessagePassing(b *testing.B) {
	err := os.Chdir("..")
	require.NoError(b, err)

	wd, err := os.Getwd()
	require.NoError(b, err)
	defer os.Chdir(wd)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd := exec.Command("neva", "run", "message_passing")
		_, err := cmd.CombinedOutput()
		require.NoError(b, err)
		require.Equal(b, 0, cmd.ProcessState.ExitCode())
	}
}
