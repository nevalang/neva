package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func TestDebugRuntimeValidation(t *testing.T) {
	_ = os.RemoveAll("gen")
	_ = os.RemoveAll("src")
	_ = os.Remove("neva.yml")
	_ = os.Remove("neva.yaml")

	t.Cleanup(func() {
		_ = os.RemoveAll("gen")
		_ = os.RemoveAll("src")
		_ = os.Remove("neva.yml")
		_ = os.Remove("neva.yaml")
	})

	e2e.Run(t, []string{"new", "."})

	// Ensure the debug validation flag doesn't interfere with running the program.
	out, _ := e2e.Run(t, []string{"run", "--debug-runtime-validation", "src"})
	require.Equal(t, "Hello, World!\n", out)

	// Build Go output with debug validation enabled so runtime sources are emitted.
	e2e.Run(t, []string{"build", "--target=go", "--debug-runtime-validation", "--output=gen", "src"})

	// Verify the generated runtime includes the debug validation helper.
	debugPath := filepath.Join("gen", "runtime", "debug_validation.go")
	debugBytes, err := os.ReadFile(debugPath)
	require.NoError(t, err)
	require.Contains(t, string(debugBytes), "func DebugValidation")

	// Build the generated Go module to ensure the code compiles.
	buildCmd := exec.Command("go", "build", "-o", "app", ".")
	buildCmd.Dir = "gen"
	buildOut, err := buildCmd.CombinedOutput()
	require.NoError(t, err, string(buildOut))

	// Execute the built binary and assert the hello world output.
	binaryPath := filepath.Join("gen", "app")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	runCmd := exec.Command(binaryPath)
	runOut, err := runCmd.CombinedOutput()
	require.NoError(t, err, string(runOut))
	require.Equal(t, "Hello, World!\n", string(runOut))
}
