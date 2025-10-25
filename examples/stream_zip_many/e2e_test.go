package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	repoRoot := filepath.Clean(filepath.Join(wd, "..", ".."))
	examplesDir := filepath.Join(repoRoot, "examples")

	binaryPath := filepath.Join(t.TempDir(), "neva")

	build := exec.Command("go", "build", "-o", binaryPath, "./cmd/neva")
	build.Dir = repoRoot

	out, err := build.CombinedOutput()
	require.NoError(t, err, string(out))

	expectedOutput := `[1,10,100]
[2,20,200]
[3,30,300]
`

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			cmd := exec.Command(binaryPath, "run", "stream_zip_many")
			cmd.Dir = examplesDir

			out, err := cmd.CombinedOutput()
			require.NoError(t, err, string(out))
			require.Equal(t, expectedOutput, string(out))

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
