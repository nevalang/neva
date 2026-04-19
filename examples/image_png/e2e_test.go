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
	out, _ := e2e.Run(t, []string{"run", "."})

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file exists.
	repoRoot := e2e.FindRepoRoot(t)
	filename := filepath.Join(repoRoot, "examples", "image_png", "minimal.png")

	_, err := os.ReadFile(filename)
	require.NoError(t, err, out)

	// Remove file output.
	_ = os.Remove(filename)
}
