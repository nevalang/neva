package test

// in this file we only check that code in the examples folder works as expected.

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

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "0_do_nothing")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

func TestEcho(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "1_echo")

	cmd.Stdin = strings.NewReader("yo\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"yo\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

func TestStructSelectorWithVerbose(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "7_struct_selector/verbose")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"Charley\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

func TestStructSelectorWithSugar(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "7_struct_selector/with_sugar")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"Charley\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
