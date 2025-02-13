// Remember - Go runs benchmark function twice:
// first time for calibration, second time for actual benchmark.
// This might affect working with relative paths!

package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkMessagePassing(b *testing.B) {
	// Store original working directory
	originalWd, err := os.Getwd()
	require.NoError(b, err)

	// Change to parent directory
	err = os.Chdir("..")
	require.NoError(b, err)

	// Ensure we return to original directory after benchmark
	defer func() {
		err := os.Chdir(originalWd)
		require.NoError(b, err)
	}()

	// Reset timer after setup
	b.ResetTimer()

	for b.Loop() {
		cmd := exec.Command("neva", "run", "message_passing")
		out, err := cmd.CombinedOutput()
		require.NoError(b, err, string(out))
	}
}
