package test

// in this file we test files designed specifically for e2e.

import (
	"os"
	"os/exec"
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
		"",
		string(out),
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
	require.Equal(
		t,
		"main/main.neva:8:2 Invalid connection\n",
		string(out),
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

// There was a bug with order of channels in IR func-call that lead to wrong answer from subtractor.
func TestOrderDependendWithArrInport(t *testing.T) {
	err := os.Chdir("./tests/order_dependend_with_arr_inport")
	require.NoError(t, err)

	defer os.Chdir(wd)

	// we do 100 attempts because there was a floating bug caused by unordered map that was
	for i := 0; i < 100; i++ {
		cmd := exec.Command("neva", "run", "main")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		require.Equal(
			t,
			"-4\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}

// There was a bug when compiler couldn't parse [-215]
func TestListWithNegInt(t *testing.T) {
	err := os.Chdir("./tests/list_with_neg_nums")
	require.NoError(t, err)

	defer os.Chdir(wd)

	for i := 0; i < 100; i++ {
		cmd := exec.Command("neva", "run", "main")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		require.Equal(
			t,
			"[-215]\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}

// Check that compiler throws human-readable error when type arguments in Node are incompatible.
func TestIncompatCompTypeArg(t *testing.T) {
	err := os.Chdir("./tests/incompat_comp_type_arg")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"main/main.neva:2:9 Incompatible types: want int | float | string, got any\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Check program with comments is parsed without errors.
func TestComments(t *testing.T) {
	err := os.Chdir("./tests/comments")
	require.NoError(t, err)

	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"<empty>\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}