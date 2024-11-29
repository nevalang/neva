package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	// Test successful case with "Alice"
	cmd := exec.Command("neva", "run", "switch_fan_out")
	cmd.Stdin = strings.NewReader("Alice\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(t, "Enter the name: ALICEalice\n", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Test panic case with "Bob"
	cmd = exec.Command("neva", "run", "switch_fan_out")
	cmd.Stdin = strings.NewReader("Bob\n")
	out, _ = cmd.CombinedOutput()
	require.Equal(t, "Enter the name: panic: Bob\n", string(out))
}
