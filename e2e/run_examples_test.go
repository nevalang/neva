//go:build e2e
// +build e2e

package test

// in this file we only check that code in examples folder works as expected.

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoNothing(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	cmd := exec.Command("neva", "run", "0_do_nothing")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		0,
		len(strings.TrimSpace(string(out))),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

func TestEcho(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	cmd := exec.Command("neva", "run", "1_echo")

	cmd.Stdin = strings.NewReader("yo\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"yo",
		strings.TrimSpace(string(out)),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}