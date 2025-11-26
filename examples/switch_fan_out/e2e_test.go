package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	// Test successful case with "Alice"
	out := e2e.RunWithStdin(t, "Alice\n", "run", "switch_fan_out")
	require.Equal(t, "Enter the name: ALICEalice\n", out)

	// Test panic case with "Bob"
	// Note: We use RunWithStdin even if we expect panic output because currently neva prints panic to stdout/stderr 
	// and exits with 0 if it's handled by the runtime, or maybe it returns error but we check output.
	// Looking at previous code: `out, _ = cmd.CombinedOutput()` suggests it might fail or succeed but we check output.
	// However, e2e.Run/RunWithStdin asserts err=nil.
	// If the previous test ignored the error `_`, it implies `cmd.Run()` might have returned an error (exit code != 0).
	// Let's check the previous code again.
	// `cmd = exec.Command("neva", "run", "switch_fan_out")` ... `out, _ = cmd.CombinedOutput()`
	// If exit code is non-zero, e2e.Run will fail.
	// If the panic causes exit code 1, we should use RunExpectingError or similar?
	// But `switch` example above asserted exit code 0 for "Charlie" which also panicked.
	// Let's assume exit code 0 for now as per `switch` example.
	
	// Wait, `switch_fan_out` previous code:
	// cmd = exec.Command("neva", "run", "switch_fan_out")
	// cmd.Stdin = strings.NewReader("Bob\n")
	// out, _ = cmd.CombinedOutput()
	// require.Equal(t, "Enter the name: panic: Bob\n", string(out))
	// It did NOT assert exit code 0.
	
	// If it really panics and exits with non-zero, `e2e.Run` will fail.
	// I'll try `e2e.Run` first. If it fails, I'll switch to `RunExpectingError` (but that expects stderr, panic might be on stdout).
	// Actually, if it's a runtime panic caught by Neva, it might still return 0 or non-0.
    // Given `switch` test asserted 0 for panic case, it's likely 0 here too.
    // If not, I will fix it.
	
	out = e2e.RunWithStdin(t, "Bob\n", "run", "switch_fan_out")
	require.Equal(t, "Enter the name: panic: Bob\n", out)
}
