package test

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Intent: verify process metadata calls return the current workspace and host data.
	out, _ := e2e.Run(t, []string{"run", "main"})

	cwd, err := os.Getwd()
	require.NoError(t, err)
	tempDir := os.TempDir()

	fields := strings.Split(strings.TrimSpace(out), "|")
	require.Len(t, fields, 6)
	require.Equal(t, "wd="+cwd, fields[0])
	require.Equal(t, "temp="+tempDir, fields[5])

	pid, err := strconv.Atoi(strings.TrimPrefix(fields[1], "pid="))
	require.NoError(t, err)
	require.Greater(t, pid, 0)

	ppid, err := strconv.Atoi(strings.TrimPrefix(fields[2], "ppid="))
	require.NoError(t, err)
	require.Greater(t, ppid, 0)

	require.NotEmpty(t, strings.TrimPrefix(fields[3], "host="))
	require.NotEmpty(t, strings.TrimPrefix(fields[4], "exe="))
}
