package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "file_read_all")

	want, err := os.ReadFile("file_read_all/main.neva")
	require.NoError(t, err, out)

	require.Equal(
		t,
		string(want),
		strings.TrimSuffix(out, "\n"),
	)
}
