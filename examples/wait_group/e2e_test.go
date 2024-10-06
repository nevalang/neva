package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "wait_group")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	t.Log(string(out))

	expected := []string{"Hello", "World!", "Neva"}
	actual := strings.Split(strings.TrimSpace(string(out)), "\n")
	require.ElementsMatch(t, expected, actual)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
