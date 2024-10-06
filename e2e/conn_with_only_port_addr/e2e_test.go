package test

import (
	"os/exec"
	"testing"

	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Contains(
		t,
		string(out),
		parser.ErrEmptyConnDef.Error(),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
