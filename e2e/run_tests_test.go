// In this file we test files designed specifically for e2e.
package test

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
