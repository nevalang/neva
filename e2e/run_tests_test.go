package test

// in this file we test files designed specifically for e2e.

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var wd string

func init() { wd, _ = os.Getwd() }

// There is special case where constant has float type but integer literal.
func TestFloatConstWithIntLit(t *testing.T) {
	err := os.Chdir("./tests/float_const_with_int_lit")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "main")

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
func TestConnWithOnlyPortAddr(t *testing.T) {
	err := os.Chdir("./tests/conn_with_only_port_addr")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	expected := strings.TrimSpace(
		"main/main.neva:8:2 Invalid connection",
	)
	require.Equal(
		t,
		expected,
		strings.TrimSpace(string(out)),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that struct selector works with port address sender.
func TestStructSelectorOnPortAddr(t *testing.T) {
	err := os.Chdir("./tests/struct_selector_on_port_addr")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"42\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
