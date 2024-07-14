package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"main/main.neva:19:8 Incompatible types: in:in -> add: Subtype and supertype must both be either literals or instances, except if supertype is union: expression any, constraint { data int, idx int, last bool }\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
