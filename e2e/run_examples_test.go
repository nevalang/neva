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

// Check that simple example with interface compiles and prints empty message.
func TestInterfacesSimple(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "3_interfaces/simple")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"<empty>\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that example with interface and imports from third-party module compiles and prints empty message.
func TestInterfacesWithImports(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "3_interfaces/with_imports")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"<empty>\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that math example with adding numbers by explicitly using port streamer works as expected.
func TestMathAddNumbersVerbose(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "4_math/add_numbers_verbose")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"3\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that math example with adding numbers by using port bridge works as expected.
func TestMathAddNumbersWithBridge(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "4_math/add_numbers_with_bridge")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"3\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that math example with multiplying numbers by using port bridge works as expected.
func TestMathMultiplyNumbers(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "4_math/multiply_numbers")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"6\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check that regex.Submatcher returns exact 3 sub-strings on a given example.
func TestRegexSubmatch(t *testing.T) {
	err := os.Chdir("../examples")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "10_regex_submatch")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"[axxxbyc xxx y]\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
