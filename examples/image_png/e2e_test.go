package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.RunExample(t, "image_png")

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file exists.
	filename := filepath.Join(e2e.ExamplesDir(t), "minimal.png")

	_, err := os.ReadFile(filename)
	require.NoError(t, err, out)

	// Remove file output.
	_ = os.Remove(filename)
}
