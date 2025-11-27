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

	out := e2e.Run(t, "run", "image_png")

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file exists.
	const filename = "minimal.png"

	_, err = os.ReadFile(filename)
	require.NoError(t, err, out)

	// Remove file output.
	os.Remove(filename)
}
