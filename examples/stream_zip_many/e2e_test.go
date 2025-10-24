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

        expectedOutput := `[1,10,100]
[2,20,200]
[3,30,300]
`

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("Iteration %d", i), func(t *testing.T) {
			cmd := exec.Command("neva", "run", "stream_zip_many")

			out, err := cmd.CombinedOutput()
			require.NoError(t, err, string(out))
			require.Equal(t, expectedOutput, string(out))

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
