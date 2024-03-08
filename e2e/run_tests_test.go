//go:build e2e
// +build e2e

package test

// in this file we test files designed specifically for e2e.

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// There is special case where constant has float type but integer literal.
func TestFloatConstWithIntLit(t *testing.T) {
	err := os.Chdir("./mod")
	require.NoError(t, err)

	cmd := exec.Command("neva", "run", "float_const_with_int_lit")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		0,
		len(strings.TrimSpace(string(out))),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Expect normal error message and not go panic trace in case of bad connection.
func TestConnectionWithOnlyPortAddr(t *testing.T) {
	err := os.Chdir("./mod")
	require.NoError(t, err)

	cmd := exec.Command("neva", "run", "connection_with_only_port_addr")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	expected := strings.TrimSpace(
		"connection_with_only_port_addr/main.neva:8:2 Invalid connection, make sure you have both sender and receiver",
	)
	require.Equal(
		t,
		expected,
		strings.TrimSpace(string(out)),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
