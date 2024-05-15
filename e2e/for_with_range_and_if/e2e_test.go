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
		`is odd
is odd
is even
is even
is odd
is odd
is even
is even
is odd
is odd
is even
is even
is odd
is odd
is even
is even
is odd
is odd
is even
is even
`,
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
