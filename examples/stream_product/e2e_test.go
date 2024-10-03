package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedOutput = `{"first": 0, "second": 0}
{"first": 0, "second": 1}
{"first": 0, "second": 2}
{"first": 1, "second": 0}
{"first": 1, "second": 1}
{"first": 1, "second": 2}
{"first": 2, "second": 0}
{"first": 2, "second": 1}
{"first": 2, "second": 2}
`

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	for i := 0; i < 1; i++ {
		t.Run(fmt.Sprintf("Run %d", i+1), func(t *testing.T) {
			cmd := exec.Command("neva", "run", "stream_product")

			out, err := cmd.CombinedOutput()
			require.NoError(t, err)
			require.Equal(t, expectedOutput, string(out))

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
