package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	expectedOutput := `{"first": 0, "second": "a"}
{"first": 1, "second": "b"}
{"first": 2, "second": "c"}
`

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			cmd := exec.Command("neva", "run", "stream_zip")

			out, err := cmd.CombinedOutput()
			require.NoError(t, err, string(out))
			require.Equal(t, expectedOutput, string(out))

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
